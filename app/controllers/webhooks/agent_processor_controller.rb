# frozen_string_literal: true

class Webhooks::AgentProcessorController < ActionController::API
  before_action :validate_secret

  def handoff
    conversation_id = params[:conversation_id]
    unless conversation_id
      render json: { error: 'conversation_id is required' }, status: :bad_request
      return
    end

    conversation = Conversation.find_by(display_id: conversation_id)
    unless conversation
      render json: { error: 'Conversation not found' }, status: :not_found
      return
    end

    # The DB was already updated by the Python processor, so we just reload
    conversation.reload

    # Dispatch the update event so ActionCable broadcasts to the frontend
    Rails.configuration.dispatcher.dispatch(Conversation::UPDATED, Time.zone.now, conversation: conversation)

    Rails.logger.info "[AgentProcessor::Handoff] Reloaded and broadcasted handoff for conversation #{conversation.display_id}"
    render json: { status: 'success' }, status: :ok
  end

  private

  def validate_secret
    expected_secret = ENV['EVOAI_CRM_API_TOKEN']

    return if expected_secret.blank?

    # Bearer token or custom header
    provided_secret = request.headers['Authorization']&.sub('Bearer ', '') || request.headers['X-Api-Token']
    return if provided_secret == expected_secret

    render json: { error: 'Unauthorized' }, status: :unauthorized
  end
end
