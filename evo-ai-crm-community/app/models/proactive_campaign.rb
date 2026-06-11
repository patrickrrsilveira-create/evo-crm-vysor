# frozen_string_literal: true

class ProactiveCampaign < ApplicationRecord
  belongs_to :account
  # belongs_to :agent, class_name: 'AgentBot', optional: true

  enum status: { DRAFT: 'DRAFT', ACTIVE: 'ACTIVE', PAUSED: 'PAUSED' }
  enum trigger_type: { LABEL_ADDED: 'LABEL_ADDED', PIPELINE_STAGE_ENTERED: 'PIPELINE_STAGE_ENTERED', SCHEDULED_DATE: 'SCHEDULED_DATE' }

  validates :name, presence: true
  validates :trigger_type, presence: true
  validates :trigger_target, presence: true
  validates :delay_hours, numericality: { greater_than_or_equal_to: 0 }
  validates :message_template, presence: true
end
