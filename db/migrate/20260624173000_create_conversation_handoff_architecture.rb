class CreateConversationHandoffArchitecture < ActiveRecord::Migration[7.0]
  def change
    # 1. Modificar a tabela conversations
    add_column :conversations, :active_agent_id, :uuid
    add_column :conversations, :state, :string, default: 'ACTIVE', null: false
    add_column :conversations, :state_version, :integer, default: 0, null: false
    add_column :conversations, :transfer_lock, :boolean, default: false, null: false

    # Add foreign key for active_agent_id to agent_bots table
    add_foreign_key :conversations, :agent_bots, column: :active_agent_id

    # 2. Tabela agent_sessions
    create_table :agent_sessions, id: :uuid do |t|
      t.references :conversation, type: :uuid, null: false, foreign_key: { on_delete: :cascade }
      t.references :agent_bot, type: :uuid, null: false, foreign_key: { on_delete: :cascade }
      t.text :summary
      t.jsonb :entities, default: {}
      t.string :state, default: 'ACTIVE', null: false
      t.datetime :closed_at

      t.timestamps
    end

    # 3. Tabela conversation_transfers
    create_table :conversation_transfers, id: :uuid do |t|
      t.references :conversation, type: :uuid, null: false, foreign_key: { on_delete: :cascade }
      t.references :from_agent, type: :uuid, foreign_key: { to_table: :agent_bots, on_delete: :cascade }
      t.references :to_agent, type: :uuid, foreign_key: { to_table: :agent_bots, on_delete: :cascade }
      t.text :reason
      t.string :status, default: 'pending', null: false
      t.datetime :completed_at

      t.timestamps
    end

    # 4. Tabela conversation_contexts
    create_table :conversation_contexts, id: :uuid do |t|
      t.references :conversation, type: :uuid, null: false, foreign_key: { on_delete: :cascade }
      t.jsonb :customer_data, default: {}
      t.jsonb :lead_data, default: {}
      t.jsonb :intent_data, default: {}
      t.jsonb :entities, default: {}
      t.text :summary

      t.timestamps
    end
  end
end
