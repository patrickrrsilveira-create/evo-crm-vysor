# frozen_string_literal: true

# == Schema Information
#
# Table name: knowledge_bases
#
#  id              :uuid             not null, primary key
#  description     :text
#  documents_count :integer          default(0)
#  name            :string           not null
#  created_at      :datetime         not null
#  updated_at      :datetime         not null
#
class KnowledgeBase < ApplicationRecord
  has_many :knowledge_documents, dependent: :destroy
  has_many :knowledge_base_agent_bots, dependent: :destroy
  has_many :agent_bots, through: :knowledge_base_agent_bots

  validates :name, presence: true
end
