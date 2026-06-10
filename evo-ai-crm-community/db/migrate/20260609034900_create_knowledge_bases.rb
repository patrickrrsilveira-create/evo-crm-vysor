class CreateKnowledgeBases < ActiveRecord::Migration[7.0]
  def change
    create_table :knowledge_bases, id: :uuid, default: -> { "gen_random_uuid()" } do |t|
      t.string :name, null: false
      t.text :description

      t.timestamps
    end
  end
end
