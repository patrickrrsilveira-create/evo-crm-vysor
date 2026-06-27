class Whatsapp::IncomingMessageWacallsService < Whatsapp::IncomingMessageBaseService
  def perform
    Rails.logger.info "WaCalls API: Service initialized with inbox: #{@inbox.present? ? @inbox.id : 'NIL'}"

    # WaCalls webhook payload:
    # { session: "sessionName", event: "message", timestamp: 12345678, data: { ... } }
    event_type = processed_params[:event]
    
    Rails.logger.info "WaCalls API: Processing event #{event_type} for session #{processed_params[:session]}"
    
    case event_type
    when 'message'
      process_wacalls_message
    when 'receipt'
      process_read_receipt
    when 'session.status' # If WaCalls sends this
      process_session_status
    else
      Rails.logger.warn "WaCalls API: Unhandled event type: #{event_type}"
    end
  end

  private

  def process_wacalls_message
    @wacalls_data = processed_params[:data]
    return if @wacalls_data.blank?

    Rails.logger.info "WaCalls API: Processing message #{@wacalls_data[:id]}"
    handle_message
  end

  def incoming?
    from_me = @wacalls_data[:fromMe] || @wacalls_data[:from_me]
    !from_me
  end

  def conversation_id
    chat_id = @wacalls_data[:from] if incoming?
    chat_id ||= @wacalls_data[:to]
    chat_id
  end

  def sender_id
    if incoming?
      if group_message?
        @wacalls_data[:author] || @wacalls_data[:participant] || conversation_id
      else
        conversation_id
      end
    else
      whatsapp_channel.phone_number
    end
  end

  def group_message?
    conversation_id&.include?('@g.us')
  end

  def raw_message_id
    @wacalls_data[:id]
  end

  def phone_number_from_jid
    num = sender_id&.split('@')&.first
    "+#{num}" if num.present?
  end

  def message_content
    @wacalls_data[:text] || @wacalls_data[:body] || ''
  end

  def message_type
    if @wacalls_data[:type] == 'image'
      'image'
    elsif @wacalls_data[:type] == 'video'
      'video'
    elsif @wacalls_data[:type] == 'audio' || @wacalls_data[:type] == 'ptt'
      'audio'
    elsif @wacalls_data[:type] == 'document' || @wacalls_data[:type] == 'file'
      'file'
    else
      'text'
    end
  end

  def attach_files
    return unless %w[image video audio file].include?(message_type)
    
    media_url = @wacalls_data[:mediaUrl] || @wacalls_data.dig(:media, :url) || @wacalls_data[:url]
    
    return unless media_url.present?

    begin
      attachment_file = Down.download(
        media_url,
        headers: {
          'Accept' => '*/*'
        }
      )
      
      @message.attachments.new(
        account_id: @message.account_id,
        file_type: message_type,
        file: attachment_file
      )
    rescue StandardError => e
      Rails.logger.error "WaCalls: Could not download attachment: #{e.message}"
    end
  end
  
  def process_read_receipt
    data = processed_params[:data]
    return if data.blank?
    
    message_id = data[:id]
    status = data[:ack] || data[:status]
    
    return unless message_id.present?
    
    message = @inbox.messages.find_by(source_id: message_id)
    return unless message
    
    status_str = case status.to_s.downcase
                 when '1', 'sent' then 'sent'
                 when '2', 'delivered' then 'delivered'
                 when '3', 'read' then 'read'
                 when '4', 'played' then 'read'
                 else 'failed'
                 end
                 
    Messages::StatusUpdateService.new(message, status_str).perform
  end

  def process_session_status
    status = processed_params.dig(:data, :status)
    if status == 'DISCONNECTED' || status == 'STOPPED' || status == 'close'
      whatsapp_channel.update!(provider_connection: { status: 'disconnected' })
    elsif status == 'CONNECTED' || status == 'WORKING' || status == 'open'
      whatsapp_channel.update!(provider_connection: { status: 'connected' })
    end
  end
end
