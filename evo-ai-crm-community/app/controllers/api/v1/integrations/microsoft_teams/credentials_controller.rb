# frozen_string_literal: true

class Api::V1::Integrations::MicrosoftTeams::CredentialsController < Api::ServiceController
  # Service-authenticated endpoint to fetch Microsoft Teams OAuth credentials
  # Used by evo-ai-processor to get credentials from global config
  # Requires X-Service-Token header for authentication

  def show
    success_response(data: microsoft_teams_credentials, message: 'Microsoft Teams credentials retrieved successfully')
  end

  private

  def microsoft_teams_credentials
    tenant_id = GlobalConfigService.load('MICROSOFT_TEAMS_TENANT_ID', nil).presence || ENV['MICROSOFT_TEAMS_TENANT_ID']
    client_id = GlobalConfigService.load('MICROSOFT_TEAMS_CLIENT_ID', nil).presence || ENV['MICROSOFT_TEAMS_CLIENT_ID']
    client_secret = GlobalConfigService.load('MICROSOFT_TEAMS_CLIENT_SECRET', nil).presence || ENV['MICROSOFT_TEAMS_CLIENT_SECRET']
    owner_id = GlobalConfigService.load('MICROSOFT_TEAMS_MEETING_OWNER_ID', nil).presence || ENV['MICROSOFT_TEAMS_MEETING_OWNER_ID']

    {
      microsoft_teams_tenant_id: tenant_id,
      microsoft_teams_client_id: client_id,
      microsoft_teams_client_secret: client_secret,
      microsoft_teams_meeting_owner_id: owner_id
    }
  end
end
