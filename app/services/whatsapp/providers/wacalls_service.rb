require 'base64'

class Whatsapp::Providers::WacallsService < Whatsapp::Providers::BaseService
  def send_message(phone_number, message)
    @message = message
    @phone_number = phone_number

    if message.attachments.present?
      send_attachment_message(phone_number, message)
    elsif message.content_type == 'input_select'
      send_interactive_message(phone_number, message)
    elsif message.content.present?
      send_text_message(phone_number, message)
    else
      @message.update!(is_unsupported: true)
      return
    end
  end

  def send_template(phone_number, template_info)
    Rails.logger.warn "WaCalls API doesn't support template messages natively via Chatwoot structure, sending as text"
    send_text_message(phone_number, build_template_text(template_info))
  end

  def sync_templates
    Rails.logger.debug "WaCalls: Templates are managed internally, no external sync needed"
  end

  def send_text_message(phone_number, message)
    text_content = message.is_a?(String) ? message : html_to_whatsapp(message.content)
    clean_number = phone_number.delete('+')

    body = {
      to: clean_number,
      text: text_content
    }

    response = HTTParty.post(
      "#{api_base_path}/api/sessions/#{instance_name}/messages/text",
      headers: api_headers,
      body: body.to_json,
      timeout: 30
    )

    process_wacalls_response(response)
  rescue StandardError => e
    Rails.logger.error "[WaCalls] Error sending text message: #{e.message}"
    handle_error_fallback(e)
  end

  def send_attachment_message(phone_number, message)
    attachment = message.attachments.first
    return send_text_message(phone_number, message) unless attachment
    
    clean_number = phone_number.delete('+')
    mime_type = attachment.file.content_type
    
    endpoint = if mime_type.start_with?('image/')
                 '/messages/image'
               elsif mime_type.start_with?('audio/')
                 '/messages/audio'
               elsif mime_type.start_with?('video/')
                 '/messages/video'
               else
                 '/messages/document'
               end

    body = {
      to: clean_number,
      url: attachment.file_url,
      mimetype: mime_type
    }

    if %w[/messages/image /messages/video].include?(endpoint)
      caption = html_to_whatsapp(message.content) || ''
      body[:caption] = caption if caption.present?
    elsif endpoint == '/messages/audio'
      body[:ptt] = true
    else
      body[:filename] = attachment.file.filename.to_s
    end

    response = HTTParty.post(
      "#{api_base_path}/api/sessions/#{instance_name}#{endpoint}",
      headers: api_headers,
      body: body.to_json,
      timeout: 30
    )

    process_wacalls_response(response)
  rescue StandardError => e
    Rails.logger.error "[WaCalls] Error sending attachment message: #{e.message}"
    handle_error_fallback(e)
  end

  def send_interactive_message(phone_number, message)
    send_text_message(phone_number, message)
  end

  def get_qr_code
    response = HTTParty.post(
      "#{api_base_path}/api/sessions/#{instance_name}/pair",
      headers: api_headers,
      timeout: 30
    )
    
    if response.success?
      { success: true }
    else
      { success: false, error: response.body }
    end
  rescue StandardError => e
    { success: false, error: e.message }
  end

  def logout
    response = HTTParty.post(
      "#{api_base_path}/api/sessions/#{instance_name}/logout",
      headers: api_headers,
      timeout: 30
    )
    
    response.success?
  rescue StandardError => e
    Rails.logger.error "[WaCalls] Error logging out: #{e.message}"
    false
  end

  private

  def api_base_path
    base_url = whatsapp_channel.provider_config['api_url'].presence || ENV['WACALLS_API_URL'] || 'https://waha.vysortech.app.br'
    base_url.chomp('/')
  end

  def api_headers
    client_id = whatsapp_channel.provider_config['api_key'].presence || ENV['WACALLS_API_KEY'] || 'evo-crm'
    headers = {
      'Content-Type' => 'application/json',
      'Accept' => 'application/json'
    }
    headers['X-Client-Id'] = client_id if client_id.present?
    headers
  end

  def instance_name
    whatsapp_channel.provider_config['instance_name'].presence || whatsapp_channel.name
  end

  def process_wacalls_response(response)
    if response.success?
      parsed_response = response.parsed_response
      if parsed_response.is_a?(Hash) && parsed_response['id']
        parsed_response['id']
      else
        SecureRandom.uuid
      end
    else
      handle_error_fallback(StandardError.new(response.body))
      nil
    end
  end

  def handle_error_fallback(error)
    return if @message.blank?
    Messages::StatusUpdateService.new(@message, 'failed', error.message).perform
  end

  def build_template_text(template_info)
    name = template_info['name']
    components = template_info['components'] || []
    
    text = "*#{name}*\n\n"
    components.each do |component|
      next unless component['type'] == 'BODY'
      text += component['text'] if component['text'].present?
    end
    
    text
  end

  def validate_provider_config
    true
  end

  def error_message
    "Error in WaCalls service"
  end
end
