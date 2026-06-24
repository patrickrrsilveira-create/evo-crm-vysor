# frozen_string_literal: true

require 'net/http'
require 'uri'
require 'json'

module MicrosoftTeams
  class MeetingService
    def initialize(webhook_url)
      @webhook_url = webhook_url
    end

    def configured?
      @webhook_url.present?
    end

    def generate_meeting_link(subject = 'Reunião via Evo CRM', payload = {})
      raise 'A URL do Webhook do Microsoft Teams não está configurada para este agente.' unless configured?

      uri = URI(@webhook_url)
      request = Net::HTTP::Post.new(uri)
      request['Content-Type'] = 'application/json'

      request.body = JSON.dump(payload.merge({
        "subject" => subject,
        "requested_at" => Time.now.utc.iso8601
      }))

      response = Net::HTTP.start(uri.hostname, uri.port, use_ssl: uri.scheme == 'https') do |http|
        http.request(request)
      end

      unless response.is_a?(Net::HTTPSuccess)
        Rails.logger.error("Falha ao disparar Webhook do Teams (n8n): #{response.body}")
        raise "Erro ao acionar integração do Teams via Webhook: #{response.body}"
      end

      data = JSON.parse(response.body) rescue {}
      
      link = data['link'] || data['joinWebUrl'] || data['url']
      
      unless link.present?
        Rails.logger.error("Resposta do Webhook do Teams não continha o link da reunião: #{response.body}")
        raise "O webhook do n8n respondeu, mas não enviou o link da reunião (esperado: 'link', 'joinWebUrl' ou 'url'). Resposta: #{response.body[0..100]}"
      end

      link
    end
  end
end
