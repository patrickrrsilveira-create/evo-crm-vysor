class RecreateKnowledgeTables < ActiveRecord::Migration[7.0]
  def up
    execute "DROP TABLE IF EXISTS knowledge_base_agent_bots CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_document_chunks CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_documents CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_bases CASCADE;"

    create_table :knowledge_bases, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.string :name, null: false
      t.text :description
      t.integer :documents_count, default: 0

      t.timestamps
    end

    create_table :knowledge_documents, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_base, null: false, foreign_key: true, type: :uuid
      t.string :title, null: false
      t.string :file_url
      t.string :content_type

      t.timestamps
    end

    create_table :knowledge_base_agent_bots, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_base, null: false, foreign_key: true, type: :uuid
      t.references :agent_bot, null: false, foreign_key: true, type: :uuid

      t.timestamps
    end

    add_index :knowledge_base_agent_bots, [:knowledge_base_id, :agent_bot_id], unique: true, name: 'idx_kb_agent_bots_on_kb_id_and_bot_id'

    create_table :knowledge_document_chunks, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_document, null: false, foreign_key: true, type: :uuid
      t.text :content, null: false
      t.jsonb :metadata, default: {}

      t.timestamps
    end

    # Add the vector column via raw SQL since pgvector gem isn't installed
    execute("ALTER TABLE knowledge_document_chunks ADD COLUMN embedding vector(1536)")
    
    # Create an index for vector similarity search
    execute("CREATE INDEX index_knowledge_document_chunks_on_embedding ON knowledge_document_chunks USING hnsw (embedding vector_cosine_ops)")
  end

  def down
    execute "DROP TABLE IF EXISTS knowledge_document_chunks CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_base_agent_bots CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_documents CASCADE;"
    execute "DROP TABLE IF EXISTS knowledge_bases CASCADE;"
  end
end
