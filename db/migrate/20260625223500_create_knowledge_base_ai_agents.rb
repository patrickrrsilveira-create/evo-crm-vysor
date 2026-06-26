class CreateKnowledgeBaseAiAgents < ActiveRecord::Migration[7.1]
  def change
    create_table :knowledge_base_ai_agents, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.uuid :knowledge_base_id, null: false
      t.uuid :ai_agent_id, null: false
      t.timestamps
    end

    add_index :knowledge_base_ai_agents, :knowledge_base_id
    add_index :knowledge_base_ai_agents, :ai_agent_id
    add_index :knowledge_base_ai_agents, [:knowledge_base_id, :ai_agent_id],
              unique: true,
              name: 'idx_kb_ai_agents_unique'

    add_foreign_key :knowledge_base_ai_agents, :knowledge_bases,
                    column: :knowledge_base_id,
                    on_delete: :cascade
    # NOTE: ai_agent_id references evo_core_agents (managed by evo-ai-processor),
    # so NO FK constraint here — just an application-level reference.
  end
end
