class CreateKnowledgeDocumentChunks < ActiveRecord::Migration[7.0]
  def up
    # Ensure pgvector is enabled
    enable_extension 'vector' unless extension_enabled?('vector')

    create_table :knowledge_document_chunks, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_document, null: false, foreign_key: true, type: :uuid
      t.text :content, null: false

      t.timestamps
    end

    # Add the vector column via raw SQL since pgvector gem isn't installed
    execute("ALTER TABLE knowledge_document_chunks ADD COLUMN embedding vector(1536)")
    
    # Create an index for vector similarity search
    execute("CREATE INDEX index_knowledge_document_chunks_on_embedding ON knowledge_document_chunks USING hnsw (embedding vector_cosine_ops)")
  end

  def down
    drop_table :knowledge_document_chunks
  end
end
