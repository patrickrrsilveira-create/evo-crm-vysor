class AgentBots::ResponseProcessor
  def initialize(agent_bot, payload)
    @agent_bot = agent_bot
    @payload = payload
  end

  def process(response)
    return unless response

    status_code = response.code.to_i
    Rails.logger.info "[AgentBot HTTP] Response Status: #{response.code} #{response.message}"

    if success_response?(status_code)
      handle_success_response(response)
    else
      handle_error_response(response)
    end
  end

  private

  def success_response?(status_code)
    status_code >= 200 && status_code < 300
  end

  def handle_success_response(response)
    Rails.logger.info "[AgentBot HTTP] Success: #{response.code}"

    begin
      parsed_response = JSON.parse(response.body)
      Rails.logger.info "[AgentBot HTTP] Parsed Response: #{parsed_response}"
      process_bot_response(parsed_response)
    rescue JSON::ParserError => e
      Rails.logger.error "[AgentBot HTTP] JSON parsing failed: #{e.message}"
    end
  end

  def handle_error_response(response)
    Rails.logger.error "[AgentBot HTTP] Error Response: #{response.code} - #{response.body}"
  end

  def process_bot_response(parsed_response)
    artifacts = extract_artifacts(parsed_response)
    return unless artifacts

    extracted = extract_content_from_artifacts(artifacts)
    text_content = extracted[:text] || ''
    file_part = extracted[:file]
    
    return if text_content.blank? && file_part.blank?

    conversation = AgentBots::ConversationFinder.new(@agent_bot, @payload).find_conversation
    return unless conversation

    select_part = extracted[:select]
    select_items = select_part&.dig('items')
    
    attachments = []
    if file_part.present?
      # Handle A2A Spec format (nested under 'file') or custom format (flat in part)
      file_info = file_part['file'] || file_part
      
      if file_info.present?
        base64_data = nil
        mime_type = file_info['mimeType'] || 'application/octet-stream'
        filename = file_info['name'] || "attachment_#{SecureRandom.hex(4)}"
        downloaded_io = nil

        if file_info['bytes'].present?
          base64_data = file_info['bytes']
        elsif file_info['url'].to_s.start_with?('data:')
          # Extract base64 from data URL: data:audio/ogg;base64,T2dn...
          match = file_info['url'].match(/data:(.*?);base64,(.*)/)
          if match
            mime_type = match[1] unless match[1].blank?
            base64_data = match[2]
            
            # Set default extension based on mime_type if name is generic
            if filename.start_with?('attachment_')
              ext = mime_type.split('/').last || 'bin'
              filename = "#{filename}.#{ext}"
            end
          end
        elsif file_info['url'].to_s.match?(/^https?:\/\//)
          begin
            require 'open-uri'
            downloaded_io = URI.open(file_info['url'])
          rescue StandardError => e
            Rails.logger.error "[AgentBot HTTP] Error downloading file from URL: #{e.message}"
          end
        end

        if downloaded_io.present?
          attachments << {
            io: downloaded_io,
            filename: filename,
            content_type: mime_type
          }
          Rails.logger.info "[AgentBot HTTP] Successfully extracted attachment from URL: #{filename} (#{mime_type})"
        elsif base64_data.present?
          begin
            decoded_bytes = Base64.decode64(base64_data)
            io = StringIO.new(decoded_bytes)
            
            attachments << {
              io: io,
              filename: filename,
              content_type: mime_type
            }
            Rails.logger.info "[AgentBot HTTP] Successfully extracted attachment: #{filename} (#{mime_type})"
          rescue StandardError => e
            Rails.logger.error "[AgentBot HTTP] Error decoding file bytes: #{e.message}"
          end
        else
          Rails.logger.warn "[AgentBot HTTP] Could not extract base64 data or download from file_info"
        end
      end
    end

    # Check if text segmentation is enabled for this agent bot
    if attachments.blank? && select_items.blank? && @agent_bot.text_segmentation_enabled && ['evo_ai_provider', 'n8n_provider'].include?(@agent_bot.bot_provider) && text_content.present?
      process_segmented_response(text_content, conversation)
    else
      # Process as a single message with signature
      final_content = build_message_with_signature(text_content)
      Rails.logger.info "[AgentBot HTTP] Bot Response Message: #{final_content}"
      
      # Try to create message normally first
      message_creator = AgentBots::MessageCreator.new(@agent_bot)
      content_type = select_items.present? ? 'input_select' : 'text'
      content_attributes = select_items.present? ? { items: select_items } : nil
      message = message_creator.create_bot_reply(final_content, conversation, content_type: content_type, content_attributes: content_attributes, attachments: attachments)
      
      # If message creation failed (conversation not eligible, e.g., after transfer),
      # try to force create it anyway (for final responses after transfer)
      unless message
        Rails.logger.info "[AgentBot HTTP] Message creation failed (conversation not eligible), attempting force create..."
        message = message_creator.create_bot_reply(final_content, conversation, force: true, content_type: content_type, content_attributes: content_attributes, attachments: attachments)
      end
      
      message
    end
  end

  def extract_artifacts(parsed_response)
    artifacts = parsed_response.dig('result', 'artifacts')
    return unless artifacts&.any?

    artifacts
  end

  def extract_content_from_artifacts(artifacts)
    text = nil
    select = nil
    file = nil

    artifacts.each do |artifact|
      next unless artifact.is_a?(Hash) && artifact['parts'].is_a?(Array)

      artifact['parts'].each do |part|
        next unless part.is_a?(Hash)

        if text.nil? && part['type'] == 'text' && part['text'].present?
          text = part['text']
        end

        if select.nil? && part['type'] == 'select'
          select = part
        end

        if file.nil? && part['type'] == 'file'
          file = part
        end
      end
    end

    { text: text, select: select, file: file }
  end


  def process_segmented_response(text_content, conversation)
    # Create segmentation service with bot's configuration
    segmentation_service = AgentBots::TextSegmentationService.new(
      @agent_bot.text_segmentation_limit || 300,
      @agent_bot.text_segmentation_min_size || 50
    )

    # Segment the text
    segments = segmentation_service.segment_text(text_content)

    Rails.logger.info "[AgentBot HTTP] Text segmented into #{segments.length} parts"
    segments.each_with_index do |segment, index|
      Rails.logger.info "[AgentBot HTTP] Segment #{index + 1}: #{segment[0..100]}#{'...' if segment.length > 100}"
    end

    # Create messages using the segmented message creator
    message_creator = AgentBots::SegmentedMessageCreator.new(@agent_bot)
    message_creator.create_messages(segments, conversation)
  end

  def build_message_with_signature(content)
    # Sempre retornar apenas o conteúdo, ignorando a assinatura do bot.
    # O painel Chatwoot tenta colocar a assinatura, mas a IA já gerencia as respostas de forma inteligente.
    # Isso resolve o bug do texto ficar "PaulaSão 22h49" mesmo com segmentação desativada.
    content
  end
end
