# frozen_string_literal: true

class Api::V1::Integrations::Openai::CredentialsController < Api::ServiceController
  # Service-authenticated endpoint to fetch OpenAI credentials
  # Used by evo-ai-processor to get credentials from global config
  # Requires X-Service-Token header for authentication

  def show
    success_response(data: openai_credentials, message: 'OpenAI credentials retrieved successfully')
  end

  private

  def openai_credentials
    api_key = GlobalConfigService.load('OPENAI_API_SECRET', nil).presence || GlobalConfigService.load('OPENAI_API_KEY', nil)
    api_url = GlobalConfigService.load('OPENAI_API_URL', nil).presence
    model = GlobalConfigService.load('OPENAI_MODEL', 'gpt-4o-mini').presence

    {
      openai_api_key: api_key,
      openai_api_url: api_url,
      openai_model: model
    }
  end
end
