# frozen_string_literal: true

require 'net/http'
require 'uri'
require 'json'

module MicrosoftTeams
  class MeetingService
    def initialize
      @tenant_id = GlobalConfigService.load('MICROSOFT_TEAMS_TENANT_ID', nil).presence || ENV['MICROSOFT_TEAMS_TENANT_ID']
      @client_id = GlobalConfigService.load('MICROSOFT_TEAMS_CLIENT_ID', nil).presence || ENV['MICROSOFT_TEAMS_CLIENT_ID']
      @client_secret = GlobalConfigService.load('MICROSOFT_TEAMS_CLIENT_SECRET', nil).presence || ENV['MICROSOFT_TEAMS_CLIENT_SECRET']
      @owner_id = GlobalConfigService.load('MICROSOFT_TEAMS_MEETING_OWNER_ID', nil).presence || ENV['MICROSOFT_TEAMS_MEETING_OWNER_ID']
    end

    def configured?
      @tenant_id.present? && @client_id.present? && @client_secret.present? && @owner_id.present?
    end

    def generate_meeting_link(subject = 'Reunião via Evo CRM')
      raise 'Microsoft Teams não está configurado. Verifique as configurações globais.' unless configured?

      token = fetch_access_token

      uri = URI("https://graph.microsoft.com/v1.0/users/#{@owner_id}/onlineMeetings")
      request = Net::HTTP::Post.new(uri)
      request['Authorization'] = "Bearer #{token}"
      request['Content-Type'] = 'application/json'

      request.body = JSON.dump({
        "startDateTime" => Time.now.utc.iso8601,
        "endDateTime" => (Time.now.utc + 1.hour).iso8601,
        "subject" => subject
      })

      response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
        http.request(request)
      end

      unless response.is_a?(Net::HTTPSuccess)
        Rails.logger.error("Falha ao criar reunião no Teams: #{response.body}")
        raise "Erro ao gerar link do Teams: #{response.body}"
      end

      data = JSON.parse(response.body)
      data['joinWebUrl']
    end

    private

    def fetch_access_token
      uri = URI("https://login.microsoftonline.com/#{@tenant_id}/oauth2/v2.0/token")
      request = Net::HTTP::Post.new(uri)
      request.set_form_data(
        'client_id' => @client_id,
        'scope' => 'https://graph.microsoft.com/.default',
        'client_secret' => @client_secret,
        'grant_type' => 'client_credentials'
      )

      response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: true) do |http|
        http.request(request)
      end

      unless response.is_a?(Net::HTTPSuccess)
        Rails.logger.error("Falha na autenticação com a Microsoft: #{response.body}")
        raise "Erro de autenticação com Azure AD. Verifique suas credenciais."
      end

      data = JSON.parse(response.body)
      data['access_token']
    end
  end
end
