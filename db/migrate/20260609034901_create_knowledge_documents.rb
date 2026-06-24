class CreateKnowledgeDocuments < ActiveRecord::Migration[7.0]
  def change
    create_table :knowledge_documents, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.references :knowledge_base, null: false, foreign_key: true, type: :uuid
      t.string :title, null: false
      t.string :file_url
      t.string :content_type

      t.timestamps
    end
  end
end
