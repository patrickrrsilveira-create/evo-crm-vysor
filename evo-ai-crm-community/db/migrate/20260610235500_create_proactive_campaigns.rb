# frozen_string_literal: true

class CreateProactiveCampaigns < ActiveRecord::Migration[7.0]
  def change
    create_table :proactive_campaigns do |t|
      t.string :name, null: false
      t.string :trigger_type, null: false
      t.string :trigger_target, null: false
      t.integer :delay_hours, null: false, default: 0
      t.integer :agent_id
      t.text :message_template, null: false
      t.string :attachment_url
      t.string :status, null: false, default: 'DRAFT'
      t.datetime :last_run_at

      t.timestamps
    end

    add_index :proactive_campaigns, :status
  end
end
