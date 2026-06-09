# frozen_string_literal: true

class KnowledgeDocument < ApplicationRecord
  belongs_to :knowledge_base

  validates :title, presence: true
end
