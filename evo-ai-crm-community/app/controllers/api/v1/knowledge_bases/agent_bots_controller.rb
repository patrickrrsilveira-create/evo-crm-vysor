class Api::V1::KnowledgeBases::AgentBotsController < Api::V1::BaseController
  require_permissions({
    index: 'knowledge_bases.read',
    create: 'knowledge_bases.update',
    destroy: 'knowledge_bases.update'
  })

  before_action :set_knowledge_base

  def index
    @agent_bots = @knowledge_base.agent_bots
    render json: @agent_bots
  end

  def create
    @agent_bot = AgentBot.find(params[:agent_bot_id])
    @knowledge_base.agent_bots << @agent_bot unless @knowledge_base.agent_bots.include?(@agent_bot)
    
    render json: @agent_bot
  end

  def destroy
    @agent_bot = @knowledge_base.agent_bots.find(params[:id])
    @knowledge_base.agent_bots.delete(@agent_bot)
    
    head :ok
  end

  private

  def set_knowledge_base
    base_id = params[:knowledge_base_id] || params[:knowledge_basis_id]
    @knowledge_base = KnowledgeBase.find(base_id)
  end
end
