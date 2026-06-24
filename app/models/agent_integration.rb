# frozen_string_literal: true

# == Schema Information
#
# Table name: evo_core_agent_integrations
#
#  id         :uuid             not null, primary key
#  config     :jsonb
#  provider   :string(100)      not null
#  created_at :datetime         not null
#  updated_at :datetime         not null
#  agent_id   :uuid             not null
#
# Indexes
#
#  idx_evo_core_agent_integrations_agent     (agent_id)
#  idx_evo_core_agent_integrations_provider  (provider)
#  unique_agent_integration                  (agent_id,provider) UNIQUE
#
# Foreign Keys
#
#  evo_core_agent_integrations_agent_id_fkey  (agent_id => evo_core_agents.id) ON DELETE => cascade
#
class AgentIntegration < ApplicationRecord
  self.table_name = 'evo_core_agent_integrations'

  # The column is `agent_id` which references `agent_bots.id`.
  belongs_to :agent_bot, foreign_key: 'agent_id'
end
