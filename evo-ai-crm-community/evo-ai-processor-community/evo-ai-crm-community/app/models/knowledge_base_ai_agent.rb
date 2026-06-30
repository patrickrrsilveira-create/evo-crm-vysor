# frozen_string_literal: true

# == Schema Information
#
# Table name: knowledge_base_ai_agents
#
#  id                :uuid             not null, primary key
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  knowledge_base_id :uuid             not null (FK → knowledge_bases)
#  ai_agent_id       :uuid             not null (soft ref → evo_core_agents, no FK)
#
class KnowledgeBaseAiAgent < ApplicationRecord
  belongs_to :knowledge_base

  # ai_agent_id is a soft reference to evo_core_agents managed by the processor service
  validates :ai_agent_id, presence: true
  validates :ai_agent_id, uniqueness: { scope: :knowledge_base_id }
end
