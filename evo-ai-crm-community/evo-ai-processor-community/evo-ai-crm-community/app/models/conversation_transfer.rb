# frozen_string_literal: true

# == Schema Information
#
# Table name: conversation_transfers
#
#  id              :uuid             not null, primary key
#  completed_at    :datetime
#  reason          :text
#  status          :string           default("pending"), not null
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#  conversation_id :uuid             not null
#  from_agent_id   :uuid
#  to_agent_id     :uuid
#
# Indexes
#
#  index_conversation_transfers_on_conversation_id  (conversation_id)
#  index_conversation_transfers_on_from_agent_id    (from_agent_id)
#  index_conversation_transfers_on_to_agent_id      (to_agent_id)
#
# Foreign Keys
#
#  fk_rails_...  (conversation_id => conversations.id)
#  fk_rails_...  (from_agent_id => agent_bots.id)
#  fk_rails_...  (to_agent_id => agent_bots.id)
#
class ConversationTransfer < ApplicationRecord
  belongs_to :conversation
  belongs_to :from_agent, class_name: 'AgentBot', optional: true
  belongs_to :to_agent, class_name: 'AgentBot', optional: true

  validates :status, presence: true
end
