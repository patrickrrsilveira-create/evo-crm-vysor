# frozen_string_literal: true

class ConversationTransfer < ApplicationRecord
  belongs_to :conversation
  belongs_to :from_agent, class_name: 'AgentBot', optional: true
  belongs_to :to_agent, class_name: 'AgentBot', optional: true

  validates :status, presence: true
end
