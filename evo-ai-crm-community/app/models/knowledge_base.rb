# frozen_string_literal: true

class KnowledgeBase < ApplicationRecord
  belongs_to :account
  has_many :knowledge_documents, dependent: :destroy
  has_many :knowledge_base_agent_bots, dependent: :destroy
  has_many :agent_bots, through: :knowledge_base_agent_bots

  validates :name, presence: true
end
