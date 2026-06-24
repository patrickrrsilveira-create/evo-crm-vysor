class Messages::AudioTranscriptionService
  include Events::Types
  pattr_initialize [:attachment!]

  def perform
    Rails.logger.info "AudioTranscriptionService: Starting for attachment #{attachment.id}"

    unless attachment.audio?
      Rails.logger.warn "AudioTranscriptionService: Attachment #{attachment.id} is not audio"
      return { error: 'Attachment is not audio' }
    end

    if attachment.meta&.[]('transcribed_text').present?
      Rails.logger.info "AudioTranscriptionService: Transcription already exists for attachment #{attachment.id}"
      return { error: 'Transcription already exists' }
    end

    unless transcription_enabled?
      Rails.logger.warn "AudioTranscriptionService: Transcription not enabled"
      return { error: 'Transcription not enabled' }
    end

    Rails.logger.info "AudioTranscriptionService: Transcription enabled, starting transcription..."
    transcribed_text = transcribe_audio

    unless transcribed_text.present?
      Rails.logger.warn "AudioTranscriptionService: Transcription returned empty result for attachment #{attachment.id}"
      return { error: 'Transcription failed' }
    end

    Rails.logger.info "AudioTranscriptionService: Transcription successful, saving to attachment #{attachment.id}"

    # Save transcription to attachment meta
    attachment.meta ||= {}
    attachment.meta['transcribed_text'] = transcribed_text
    attachment.save!

    # Reload attachment and message to ensure fresh data for broadcast
    message = attachment.message
    attachment.reload
    message.reload

    # Clear attachments association cache to ensure fresh data when push_event_data is called
    # This forces Rails to reload attachments from database when push_event_data accesses them
    message.association(:attachments).reset

    # Broadcast message update to frontend so transcription appears
    Rails.configuration.dispatcher.dispatch(
      MESSAGE_UPDATED,
      Time.zone.now,
      message: message,
      previous_changes: { 'attachments' => [attachment.id] }
    )

    Rails.logger.info "AudioTranscriptionService: Transcription saved successfully for attachment #{attachment.id}"
    { success: true, transcribed_text: transcribed_text }
  rescue StandardError => e
    Rails.logger.error "AudioTranscriptionService: Error for attachment #{attachment.id}: #{e.message}"
    Rails.logger.error e.backtrace.join("\n")
    { error: e.message }
  end

  private

  def transcription_enabled?
    # Priority 1: Check new AUDIO_TRANSCRIPTION_ENABLED global configuration
    global_enabled = GlobalConfigService.load('AUDIO_TRANSCRIPTION_ENABLED', nil)
    
    # Fallback to OPENAI_ENABLE_AUDIO_TRANSCRIPTION for backward compatibility
    if global_enabled.nil?
      global_enabled = GlobalConfigService.load('OPENAI_ENABLE_AUDIO_TRANSCRIPTION', nil)
    end
    
    unless global_enabled.nil?
      # Convert to boolean - handle both boolean and string values
      enabled = if global_enabled.is_a?(TrueClass) || global_enabled.is_a?(FalseClass)
                  global_enabled
                else
                  case global_enabled.to_s.downcase
                  when 'true', '1', 'yes', 'on'
                    true
                  when 'false', '0', 'no', 'off', ''
                    false
                  else
                    false
                  end
                end

      Rails.logger.info "AudioTranscriptionService: Global config value: #{global_enabled.inspect} (#{global_enabled.class}), converted to: #{enabled.inspect}"

      if enabled
        Rails.logger.info "AudioTranscriptionService: Transcription enabled via global config"
        api_key = get_audio_api_key
        unless api_key.present?
          Rails.logger.warn "AudioTranscriptionService: Global config enabled but API key not configured"
        end
        return api_key.present?
      else
        Rails.logger.info "AudioTranscriptionService: Transcription disabled via global config"
        return false
      end
    end

    # Priority 2: Check OpenAI integration hook settings (legacy)
    openai_hook = Hook.find_by(app_id: 'openai')
    return false unless openai_hook&.enabled?
    return false unless openai_hook.settings&.[]('enable_audio_transcription') == true

    openai_hook.settings&.[]('api_key').present?
  end

  def transcribe_audio
    return nil unless attachment.file.attached?

    # Get API key from integration hook or global config
    api_key = get_audio_api_key
    return nil unless api_key.present?

    # Download audio file
    audio_file = download_audio_file
    return nil unless audio_file

    # Call Audio Transcription API
    response = call_audio_transcription_api(api_key, audio_file)

    # Clean up temp file
    File.delete(audio_file.path) if File.exist?(audio_file.path)

    response&.dig('text')
  rescue StandardError => e
    Rails.logger.error "Audio Transcription API error: #{e.message}"
    nil
  end

  def get_audio_api_key
    # Priority 1: Dedicated Audio Transcription configuration
    global_api_key = GlobalConfigService.load('AUDIO_TRANSCRIPTION_API_SECRET', nil)
    if global_api_key.present?
      Rails.logger.info "AudioTranscriptionService: Using dedicated Audio Transcription API key"
      return global_api_key
    end
    
    # Priority 2: Fallback to global OpenAI configuration
    fallback_api_key = GlobalConfigService.load('OPENAI_API_SECRET', nil)
    if fallback_api_key.present?
      Rails.logger.info "AudioTranscriptionService: Using fallback global OpenAI API key"
      return fallback_api_key
    end

    # Priority 3: Fallback to hook settings for backward compatibility
    openai_hook = Hook.find_by(app_id: 'openai')
    unless openai_hook&.enabled?
      Rails.logger.warn "AudioTranscriptionService: OpenAI hook not found or not enabled"
      return nil
    end

    api_key = openai_hook.settings&.[]('api_key')
    if api_key.present?
      Rails.logger.info "AudioTranscriptionService: Using hook OpenAI API key as fallback"
    else
      Rails.logger.warn "AudioTranscriptionService: OpenAI hook exists but API key is not configured"
    end
    api_key
  end

  def download_audio_file
    return nil unless attachment.file.attached?

    # Retry with exponential backoff to handle race condition where file
    # might not be fully uploaded to S3 yet
    max_retries = 3
    retry_delay = 1 # seconds

    max_retries.times do |attempt|
      begin
        # Create temp file with valid extension
        file_extension = attachment.extension.presence || 'ogg'
        temp_file = Tempfile.new(['audio', ".#{file_extension}"])
        temp_file.binmode

        # Download from ActiveStorage
        attachment.file.download do |chunk|
          temp_file.write(chunk)
        end

        temp_file.rewind
        Rails.logger.info "AudioTranscriptionService: Successfully downloaded audio file (attempt #{attempt + 1})"
        return temp_file
      rescue ActiveStorage::FileNotFoundError => e
        if attempt < max_retries - 1
          wait_time = retry_delay * (2 ** attempt)
          Rails.logger.warn "AudioTranscriptionService: File not found, retrying in #{wait_time}s (attempt #{attempt + 1}/#{max_retries})"
          sleep(wait_time)
        else
          Rails.logger.error "AudioTranscriptionService: Error downloading audio file after #{max_retries} attempts: #{e.message}"
          return nil
        end
      rescue StandardError => e
        Rails.logger.error "AudioTranscriptionService: Error downloading audio file: #{e.message}"
        return nil
      end
    end

    nil
  end

  def call_audio_transcription_api(api_key, audio_file)
    require 'net/http'
    require 'uri'

    # Get base URL, fallback to OpenAI if not set
    base_url = GlobalConfigService.load('AUDIO_TRANSCRIPTION_API_URL', nil)
    base_url = GlobalConfigService.load('OPENAI_API_URL', 'https://api.openai.com/v1') if base_url.blank?
    transcription_url = "#{base_url}/audio/transcriptions"

    uri = URI(transcription_url)
    http = Net::HTTP.new(uri.host, uri.port)
    http.use_ssl = (uri.scheme == 'https')

    request = Net::HTTP::Post.new(uri.path)
    request['Authorization'] = "Bearer #{api_key}"

    # Ensure we have a valid extension (default to 'ogg' for audio files)
    file_extension = attachment.extension.presence || 'ogg'
    filename = "audio.#{file_extension}"

    # Get model, fallback to OpenAI model, then whisper-1
    model_name = GlobalConfigService.load('AUDIO_TRANSCRIPTION_MODEL', nil)
    model_name = GlobalConfigService.load('OPENAI_MODEL', 'whisper-1') if model_name.blank?

    if base_url.include?('openrouter.ai')
      request['Content-Type'] = 'application/json'
      
      require 'base64'
      audio_data = Base64.strict_encode64(audio_file.read)
      audio_file.rewind
      
      payload = {
        model: model_name,
        input_audio: {
          data: audio_data,
          format: file_extension
        }
      }
      
      detected_language = detect_language
      payload[:language] = detected_language if detected_language.present?
      
      request.body = payload.to_json
    else
      form_data = [
        ['file', audio_file, { filename: filename }],
        ['model', model_name]
      ]

      # Only add language if detect_language returns a non-nil value
      detected_language = detect_language
      form_data << ['language', detected_language] if detected_language.present?

      request.set_form(form_data, 'multipart/form-data')
    end

    Rails.logger.info "Audio Transcription API request to #{transcription_url} for attachment #{attachment.id} using model #{model_name}"
    response = http.request(request)
    Rails.logger.info "Audio Transcription API response: #{response.code} - #{response.body[0..200]}"

    if response.code == '200'
      JSON.parse(response.body)
    else
      Rails.logger.error "Audio Transcription API error: #{response.code} - #{response.body}"
      nil
    end
  rescue StandardError => e
    Rails.logger.error "Audio Transcription API request error: #{e.message}"
    Rails.logger.error e.backtrace.join("\n")
    nil
  end

  def detect_language
    # Try to detect language from global locale config
    locale = GlobalConfigService.load('DEFAULT_LOCALE', nil)
    return 'pt' if locale&.start_with?('pt')
    return 'es' if locale&.start_with?('es')
    return 'fr' if locale&.start_with?('fr')
    return 'de' if locale&.start_with?('de')
    return 'it' if locale&.start_with?('it')

    # Default to auto-detect
    nil
  end
end

