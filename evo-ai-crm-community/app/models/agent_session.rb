# frozen_string_literal: true

class AgentSession < ApplicationRecord
  belongs_to :conversation
  belongs_to :agent_bot

  validates :state, presence: true
end
