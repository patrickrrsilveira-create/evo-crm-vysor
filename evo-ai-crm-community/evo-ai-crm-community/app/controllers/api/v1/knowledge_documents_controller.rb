class Api::V1::KnowledgeDocumentsController < Api::V1::BaseController
  require_permissions({
    index: 'knowledge_bases.read',
    destroy: 'knowledge_bases.update'
  })

  before_action :set_knowledge_base
  before_action :set_knowledge_document, only: [:destroy]

  def index
    @documents = @knowledge_base.knowledge_documents.order(created_at: :desc)
    render json: @documents
  end

  def destroy
    @knowledge_document.destroy!
    head :no_content
  end

  private

  def set_knowledge_base
    @knowledge_base = KnowledgeBase.find(params[:knowledge_basis_id] || params[:knowledge_base_id])
  end

  def set_knowledge_document
    @knowledge_document = @knowledge_base.knowledge_documents.find(params[:id])
  end
end
