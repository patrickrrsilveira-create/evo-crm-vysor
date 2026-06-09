class CreateKnowledgeBaseAgentBots < ActiveRecord::Migration[7.0]
  def change
    create_table :knowledge_base_agent_bots do |t|
      t.references :knowledge_base, null: false, foreign_key: true
      t.references :agent_bot, null: false, foreign_key: true

      t.timestamps
    end

    add_index :knowledge_base_agent_bots, [:knowledge_base_id, :agent_bot_id], unique: true, name: 'idx_kb_agent_bots_on_kb_id_and_bot_id'
  end
end
