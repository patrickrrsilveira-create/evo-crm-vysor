class CreateKnowledgeDocumentChunks < ActiveRecord::Migration[7.0]
  def up
    # Ensure pgvector is enabled
    enable_extension 'vector' unless extension_enabled?('vector')

    create_table :knowledge_document_chunks, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_document, null: false, foreign_key: true, type: :uuid
      t.text :content, null: false
      t.vector :embedding, limit: 1536 # OpenAI text-embedding-ada-002 size

      t.timestamps
    end
    
    # Create an index for vector similarity search
    add_index :knowledge_document_chunks, :embedding, using: :hnsw, opclass: :vector_cosine_ops
  end

  def down
    drop_table :knowledge_document_chunks
  end
end
