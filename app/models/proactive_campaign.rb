# frozen_string_literal: true

# == Schema Information
#
# Table name: proactive_campaigns
#
#  id               :bigint           not null, primary key
#  attachment_url   :string
#  delay_hours      :integer          default(0), not null
#  last_run_at      :datetime
#  message_template :text             not null
#  name             :string           not null
#  status           :string           default("DRAFT"), not null
#  trigger_target   :string           not null
#  trigger_type     :string           not null
#  created_at       :datetime         not null
#  updated_at       :datetime         not null
#  agent_id         :integer
#
# Indexes
#
#  index_proactive_campaigns_on_status  (status)
#
class ProactiveCampaign < ApplicationRecord
  # belongs_to :agent, class_name: 'AgentBot', optional: true

  enum status: { DRAFT: 'DRAFT', ACTIVE: 'ACTIVE', PAUSED: 'PAUSED' }
  enum trigger_type: { 
    LABEL_ADDED: 'LABEL_ADDED', 
    PIPELINE_STAGE_ENTERED: 'PIPELINE_STAGE_ENTERED', 
    SCHEDULED_DATE: 'SCHEDULED_DATE',
    CONTACT_CREATED: 'CONTACT_CREATED',
    CONVERSATION_OPENED: 'CONVERSATION_OPENED',
    CONVERSATION_RESOLVED: 'CONVERSATION_RESOLVED'
  }

  validates :name, presence: true
  validates :trigger_type, presence: true
  validates :trigger_target, presence: true
  validates :delay_hours, numericality: { greater_than_or_equal_to: 0 }
  validates :message_template, presence: true
end
