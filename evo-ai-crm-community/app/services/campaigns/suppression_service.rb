module Campaigns
  class SuppressionService
    OPT_OUT_KEYWORDS = %w[
      sair parar unsubscribe
      stop remove quit
      tirar lista
    ].freeze

    def initialize(account)
      @account = account
    end

    def detect_and_suppress(message)
      return false unless message.incoming?
      return false unless message.content.present?

      text = message.content.upcase
      return false unless matches_opt_out?(text)

      contact = message.sender
      channel_type = detect_channel_type(message.conversation&.inbox)

      suppress(contact, channel_type)
      true
    end

    def suppress(contact, channel_type)
      SuppressionRecord.find_or_create_by(
        account_id: @account.id,
        contact_id: contact.id,
        channel_type: channel_type
      )
    end

    def is_suppressed?(contact, channel_type = nil)
      q = SuppressionRecord.where(
        account_id: @account.id,
        contact_id: contact.id
      )
      q = q.where(channel_type: channel_type) if channel_type.present?
      q.exists?
    end

    def remove_suppression(contact, channel_type)
      SuppressionRecord.where(
        account_id: @account.id,
        contact_id: contact.id,
        channel_type: channel_type
      ).delete_all
    end

    private

    def matches_opt_out?(text)
      OPT_OUT_KEYWORDS.any? { |kw| text.include?(kw) }
    end

    def detect_channel_type(inbox)
      return 'whatsapp' if inbox&.channel_type == 'Api::Channel::TwitterChannel'
      inbox&.channel_type.to_s.downcase.gsub(/api::channel::|\s+channel/, '').presence || 'unknown'
    end
  end

  class SuppressionRecord < ApplicationRecord
    self.table_name = 'cmp_suppression'
    belongs_to :account
    belongs_to :contact
    validates :account_id, :contact_id, :channel_type, presence: true
    validates :account_id, uniqueness: { scope: [:contact_id, :channel_type] }
  end
end
