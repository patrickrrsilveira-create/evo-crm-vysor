module Campaigns
  class SegmentResolverService
    def initialize(account, filters = {})
      @account = account
      @filters = filters
    end

    def resolve
      query = @account.contacts.persons

      if @filters['label_ids'].present?
        label_ids = Array(@filters['label_ids'])
        query = query.joins(:labels).where(labels: { id: label_ids }).distinct
      end

      if @filters['contact_type'].present?
        query = query.where(contact_type: @filters['contact_type'])
      end

      if @filters['timezone'].present?
        query = query.where('additional_attributes ->> ? = ?', 'timezone', @filters['timezone'])
      end

      if @filters['custom_attributes'].present?
        @filters['custom_attributes'].each do |key, value|
          query = query.where('custom_attributes ->> ? = ?', key, value)
        end
      end

      query.where.not(phone_number: [nil, '']).pluck(:id, :phone_number, 'additional_attributes ->> ?', 'timezone').map do |id, phone, tz|
        { contact_id: id, recipient: phone, timezone: tz || 'UTC' }
      end
    end

    def count
      resolve.count
    end
  end
end
