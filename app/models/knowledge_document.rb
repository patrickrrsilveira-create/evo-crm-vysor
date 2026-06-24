# frozen_string_literal: true

# == Schema Information
#
# Table name: knowledge_documents
#
#  id                :uuid             not null, primary key
#  content_type      :string
#  file_url          :string
#  title             :string           not null
#  created_at        :datetime         not null
#  updated_at        :datetime         not null
#  knowledge_base_id :uuid             not null
#
# Indexes
#
#  index_knowledge_documents_on_knowledge_base_id  (knowledge_base_id)
#
# Foreign Keys
#
#  fk_rails_...  (knowledge_base_id => knowledge_bases.id)
#
class KnowledgeDocument < ApplicationRecord
  belongs_to :knowledge_base
  has_many :knowledge_document_chunks, dependent: :destroy

  validates :title, presence: true
end
