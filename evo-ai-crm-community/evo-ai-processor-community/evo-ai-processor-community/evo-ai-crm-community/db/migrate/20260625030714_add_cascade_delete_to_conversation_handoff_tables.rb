class AddCascadeDeleteToConversationHandoffTables < ActiveRecord::Migration[7.0]
  def up
    # Remove existing non-cascading foreign keys
    remove_foreign_key :agent_sessions, :conversations
    remove_foreign_key :agent_sessions, :agent_bots
    remove_foreign_key :conversation_transfers, :conversations
    remove_foreign_key :conversation_transfers, column: :from_agent_id
    remove_foreign_key :conversation_transfers, column: :to_agent_id
    remove_foreign_key :conversation_contexts, :conversations

    # Add cascading foreign keys
    add_foreign_key :agent_sessions, :conversations, on_delete: :cascade
    add_foreign_key :agent_sessions, :agent_bots, on_delete: :cascade
    add_foreign_key :conversation_transfers, :conversations, on_delete: :cascade
    add_foreign_key :conversation_transfers, :agent_bots, column: :from_agent_id, on_delete: :cascade
    add_foreign_key :conversation_transfers, :agent_bots, column: :to_agent_id, on_delete: :cascade
    add_foreign_key :conversation_contexts, :conversations, on_delete: :cascade
  end

  def down
    remove_foreign_key :agent_sessions, :conversations
    remove_foreign_key :agent_sessions, :agent_bots
    remove_foreign_key :conversation_transfers, :conversations
    remove_foreign_key :conversation_transfers, column: :from_agent_id
    remove_foreign_key :conversation_transfers, column: :to_agent_id
    remove_foreign_key :conversation_contexts, :conversations

    add_foreign_key :agent_sessions, :conversations
    add_foreign_key :agent_sessions, :agent_bots
    add_foreign_key :conversation_transfers, :conversations
    add_foreign_key :conversation_transfers, :agent_bots, column: :from_agent_id
    add_foreign_key :conversation_transfers, :agent_bots, column: :to_agent_id
    add_foreign_key :conversation_contexts, :conversations
  end
end
