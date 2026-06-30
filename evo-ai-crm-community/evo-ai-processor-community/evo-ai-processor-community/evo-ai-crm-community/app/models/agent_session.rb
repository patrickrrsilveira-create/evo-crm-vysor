# frozen_string_literal: true

# == Schema Information
#
# Table name: agent_sessions
#
#  id              :uuid             not null, primary key
#  closed_at       :datetime
#  entities        :jsonb
#  state           :string           default("ACTIVE"), not null
#  summary         :text
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#  agent_bot_id    :uuid             not null
#  conversation_id :uuid             not null
#
# Indexes
#
#  index_agent_sessions_on_agent_bot_id     (agent_bot_id)
#  index_agent_sessions_on_conversation_id  (conversation_id)
#
# Foreign Keys
#
#  fk_rails_...  (agent_bot_id => agent_bots.id)
#  fk_rails_...  (conversation_id => conversations.id)
#
class AgentSession < ApplicationRecord
  belongs_to :conversation
  belongs_to :agent_bot

  validates :state, presence: true
end
