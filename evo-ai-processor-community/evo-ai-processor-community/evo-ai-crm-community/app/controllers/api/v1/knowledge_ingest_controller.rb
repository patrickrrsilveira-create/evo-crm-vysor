class Api::V1::KnowledgeIngestController < Api::V1::BaseController
  require_permissions({
    file: 'knowledge_bases.update',
    url: 'knowledge_bases.update'
  })
  def file
    processor_url = ENV.fetch('AI_PROCESSOR_URL', 'http://evo_processor:8000')

    if params[:file].blank?
      return render json: { error: 'Arquivo é obrigatório' }, status: :bad_request
    end

    # Forward the multipart file request to the python processor
    begin
      response = RestClient::Request.execute(
        method: :post,
        url: "#{processor_url}/api/v1/knowledge/ingest/file",
        headers: {
          'Authorization' => request.headers['Authorization']
        },
        payload: {
          multipart: true,
          knowledge_base_id: params[:knowledge_base_id],
          title: params[:title],
          file: File.new(params[:file].tempfile.path, 'rb')
        }
      )

      render json: JSON.parse(response.body), status: response.code
    rescue RestClient::ExceptionWithResponse => e
      Rails.logger.error "Error proxying to AI Processor: #{e.response}"
      render json: { error: "Processor proxy error: #{e.response}" }, status: e.http_code
    rescue => e
      Rails.logger.error "Error proxying to AI Processor: #{e.message}"
      render json: { error: "Processor proxy error: #{e.message}" }, status: :internal_server_error
    end
  end

  def url
    processor_url = ENV.fetch('AI_PROCESSOR_URL', 'http://evo_processor:8000')

    begin
      response = HTTParty.post(
        "#{processor_url}/api/v1/knowledge/ingest/url",
        headers: {
          'Authorization' => request.headers['Authorization'],
          'Content-Type' => 'application/json'
        },
        body: {
          knowledge_base_id: params[:knowledge_base_id],
          title: params[:title],
          url: params[:url]
        }.to_json
      )

      render json: response.parsed_response, status: response.code
    rescue => e
      Rails.logger.error "Error proxying to AI Processor: #{e.message}"
      render json: { error: "Processor proxy error: #{e.message}" }, status: :internal_server_error
    end
  end
end
