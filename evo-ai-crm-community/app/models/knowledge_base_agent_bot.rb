# frozen_string_literal: true

class KnowledgeBaseAgentBot < ApplicationRecord
  belongs_to :knowledge_base
  belongs_to :agent_bot

  validates :knowledge_base_id, uniqueness: { scope: :agent_bot_id }
end
