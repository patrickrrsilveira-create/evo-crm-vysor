module Campaigns
  class MessageListener
    def message_created(message)
      return unless message.account.present?

      service = SuppressionService.new(message.account)
      service.detect_and_suppress(message)
    rescue StandardError => e
      Rails.logger.error("Campaigns::MessageListener error: #{e.message}")
    end
  end
end
