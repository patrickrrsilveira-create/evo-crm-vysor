# frozen_string_literal: true

class Api::V1::KnowledgeBases::AiAgentsController < Api::V1::BaseController
  require_permissions({
    index: 'knowledge_bases.read'
  })

  before_action :set_knowledge_base

  def index
    links = @knowledge_base.knowledge_base_ai_agents
    
    render json: links.map { |l| { id: l.ai_agent_id, knowledge_base_id: l.knowledge_base_id } }
  end

  private

  def set_knowledge_base
    # Note: Rails inflector maps resources :knowledge_bases to params[:knowledge_basis_id]
    @knowledge_base = KnowledgeBase.find(params[:knowledge_basis_id] || params[:knowledge_base_id])
  end
end
