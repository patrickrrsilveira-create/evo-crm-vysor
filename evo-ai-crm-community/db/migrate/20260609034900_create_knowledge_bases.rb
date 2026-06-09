class CreateKnowledgeBases < ActiveRecord::Migration[7.0]
  def change
    create_table :knowledge_bases do |t|
      t.string :name, null: false
      t.text :description
      t.references :account, null: false, foreign_key: true

      t.timestamps
    end
  end
end
