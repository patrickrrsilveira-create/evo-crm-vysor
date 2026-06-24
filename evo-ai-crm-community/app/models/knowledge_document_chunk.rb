# frozen_string_literal: true

# == Schema Information
#
# Table name: knowledge_document_chunks
#
#  id                    :uuid             not null, primary key
#  content               :text             not null
#  embedding             :vector(1536)
#  metadata              :jsonb
#  created_at            :datetime         not null
#  updated_at            :datetime         not null
#  knowledge_document_id :uuid             not null
#
# Indexes
#
#  index_knowledge_document_chunks_on_embedding              (embedding) USING hnsw
#  index_knowledge_document_chunks_on_knowledge_document_id  (knowledge_document_id)
#
# Foreign Keys
#
#  fk_rails_...  (knowledge_document_id => knowledge_documents.id)
#
class KnowledgeDocumentChunk < ApplicationRecord
  belongs_to :knowledge_document

  validates :content, presence: true
end
