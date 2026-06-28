# frozen_string_literal: true

# Controller to link/unlink AI agents (from evo_core_agents) to knowledge bases.
# Uses the ai_agents/:ai_agent_id/knowledge_bases route pattern,
# mirroring AiAgents::ProductsController.
#
# The ai_agent_id is a soft UUID reference (no FK) to the evo_core_agents table
# managed by the evo-ai-processor service.
class Api::V1::AiAgents::KnowledgeBasesController < Api::V1::BaseController
  require_permissions({
    index:   'knowledge_bases.read',
    create:  'knowledge_bases.update',
    destroy: 'knowledge_bases.update'
  })

  # GET /api/v1/ai_agents/:ai_agent_id/knowledge_bases
  # Returns all knowledge bases linked to this AI agent
  def index
    linked_kb_ids = KnowledgeBaseAiAgent
                      .where(ai_agent_id: ai_agent_id)
                      .pluck(:knowledge_base_id)

    knowledge_bases = KnowledgeBase.where(id: linked_kb_ids)

    success_response(
      data: knowledge_bases.map { |kb|
        {
          id: kb.id,
          name: kb.name,
          description: kb.description,
          documents_count: kb.documents_count,
          created_at: kb.created_at,
          updated_at: kb.updated_at
        }
      },
      message: 'Knowledge bases retrieved successfully'
    )
  end

  # POST /api/v1/ai_agents/:ai_agent_id/knowledge_bases
  # Body: { knowledge_base_id: uuid }
  def create
    kb_id = params[:knowledge_base_id]

    if kb_id.blank?
      return error_response(
        code: ApiErrorCodes::VALIDATION_ERROR,
        message: 'Provide knowledge_base_id',
        status: :unprocessable_entity
      )
    end

    unless KnowledgeBase.exists?(id: kb_id)
      return error_response(
        code: ApiErrorCodes::RESOURCE_NOT_FOUND,
        message: "Knowledge base not found: #{kb_id}",
        status: :not_found
      )
    end

    record = KnowledgeBaseAiAgent.where(
      ai_agent_id: ai_agent_id,
      knowledge_base_id: kb_id
    ).first_or_create!

    success_response(
      data: { ai_agent_id: ai_agent_id, knowledge_base_id: kb_id },
      message: 'Knowledge base linked to AI agent successfully',
      status: :created
    )
  rescue ActiveRecord::RecordInvalid => e
    error_response(
      code: ApiErrorCodes::VALIDATION_ERROR,
      message: e.message
    )
  end

  # DELETE /api/v1/ai_agents/:ai_agent_id/knowledge_bases/:id
  def destroy
    record = KnowledgeBaseAiAgent.find_by(
      ai_agent_id: ai_agent_id,
      knowledge_base_id: params[:id]
    )

    if record.nil?
      return error_response(
        code: ApiErrorCodes::RESOURCE_NOT_FOUND,
        message: 'Link not found',
        status: :not_found
      )
    end

    record.destroy!

    success_response(
      data: { ai_agent_id: ai_agent_id, knowledge_base_id: params[:id] },
      message: 'Knowledge base unlinked from AI agent successfully'
    )
  end

  private

  def ai_agent_id
    params[:ai_agent_id]
  end
end
