class EvolutionGo::SyncHistoryJob < ApplicationJob
  queue_as :low

  def perform(channel_id, limit = 100)
    channel = Channel::Whatsapp.find_by(id: channel_id)
    return unless channel && channel.provider == 'evolution_go'

    service = Whatsapp::Providers::EvolutionGoService.new(whatsapp_channel: channel)
    
    # 1. Fetch chats
    chats = service.fetch_historical_chats
    return if chats.empty?
    
    Rails.logger.info "[Evolution Go Sync] Found #{chats.size} historical chats for channel #{channel_id}"
    
    chats.each do |chat|
      remote_jid = chat['id'] || chat['remoteJid']
      next unless remote_jid&.include?('@s.whatsapp.net') # Only individual contacts
      
      # 2. Fetch messages for each chat
      messages = service.fetch_historical_messages(remote_jid, limit)
      next if messages.empty?
      
      Rails.logger.info "[Evolution Go Sync] Syncing #{messages.size} messages for #{remote_jid}"
      
      # 3. Use the existing messages.set sync logic by enqueuing a webhook job
      payload = {
        'event' => 'messages.set',
        'data' => messages,
        'instanceId' => channel.provider_config['instance_uuid'],
        'instanceToken' => channel.provider_config['instance_token'],
        'evolution_go' => true
      }
      
      Webhooks::WhatsappEventsJob.perform_later(payload)
    end
    
    Rails.logger.info "[Evolution Go Sync] Completed enqueueing historical messages for channel #{channel_id}"
  end
end
