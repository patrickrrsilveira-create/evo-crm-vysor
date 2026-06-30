# frozen_string_literal: true

# == Schema Information
#
# Table name: conversation_contexts
#
#  id              :uuid             not null, primary key
#  customer_data   :jsonb
#  entities        :jsonb
#  intent_data     :jsonb
#  lead_data       :jsonb
#  summary         :text
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#  conversation_id :uuid             not null
#
# Indexes
#
#  index_conversation_contexts_on_conversation_id  (conversation_id)
#
# Foreign Keys
#
#  fk_rails_...  (conversation_id => conversations.id)
#
class ConversationContext < ApplicationRecord
  belongs_to :conversation
end
