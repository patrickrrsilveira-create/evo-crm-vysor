# frozen_string_literal: true

class Api::V1::Integrations::GoogleCalendar::CredentialsController < Api::ServiceController
  # Service-authenticated endpoint to fetch Google Calendar OAuth credentials
  # Used by evo-ai-processor to get credentials from global config
  # Requires X-Service-Token header for authentication

  def show
    success_response(data: google_calendar_credentials, message: 'Google Calendar credentials retrieved successfully')
  end

  private

  def google_calendar_credentials
    client_id = GlobalConfigService.load('GOOGLE_CALENDAR_CLIENT_ID', nil).presence || GlobalConfigService.load('GOOGLE_OAUTH_CLIENT_ID', nil)
    client_secret = GlobalConfigService.load('GOOGLE_CALENDAR_CLIENT_SECRET', nil).presence || GlobalConfigService.load('GOOGLE_OAUTH_CLIENT_SECRET', nil)

    frontend_url = ENV.fetch('FRONTEND_URL', 'http://localhost:5173').chomp('/')
    default_redirect_uri = "#{frontend_url}/google-calendar/callback"

    {
      google_calendar_client_id: client_id,
      google_calendar_client_secret: client_secret,
      google_calendar_redirect_uri: GlobalConfigService.load('GOOGLE_CALENDAR_REDIRECT_URI', nil).presence || default_redirect_uri
    }
  end
end
