# frozen_string_literal: true

# == Schema Information
#
# Table name: knowledge_base_agent_bots
#
#  id                :uuid             not null, primary key
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  agent_bot_id      :uuid             not null
#  knowledge_base_id :uuid             not null
#
# Indexes
#
#  idx_kb_agent_bots_on_kb_id_and_bot_id                 (knowledge_base_id,agent_bot_id) UNIQUE
#  index_knowledge_base_agent_bots_on_agent_bot_id       (agent_bot_id)
#  index_knowledge_base_agent_bots_on_knowledge_base_id  (knowledge_base_id)
#
# Foreign Keys
#
#  fk_rails_...  (agent_bot_id => agent_bots.id)
#  fk_rails_...  (knowledge_base_id => knowledge_bases.id)
#
class KnowledgeBaseAgentBot < ApplicationRecord
  belongs_to :knowledge_base
  belongs_to :agent_bot

  validates :knowledge_base_id, uniqueness: { scope: :agent_bot_id }
end
