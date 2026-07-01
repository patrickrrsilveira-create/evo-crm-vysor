class Api::V1::CampaignsController < Api::ServiceController
  before_action :authenticate_service_token!, only: [:resolve_segment]

  def resolve_segment
    account = find_account

    filter_params = params.require(:filters).permit(
      label_ids: [],
      :contact_type,
      :timezone,
      custom_attributes: {}
    )

    service = Campaigns::SegmentResolverService.new(account, filter_params)
    contacts = service.resolve

    render json: {
      account_id: account.id,
      total: contacts.count,
      contacts: contacts
    }, status: :ok
  end

  private

  def find_account
    @account ||= current_user.accounts.find(params[:account_id])
  end

  def authenticate_service_token!
    token = request.headers['X-Campaign-Engine-Token'] || params[:token]
    raise Unauthorized unless valid_service_token?(token)
  end

  def valid_service_token?
    token = ENV['CRM_API_TOKEN']
    token.present? && token == (request.headers['X-Campaign-Engine-Token'] || params[:token])
  end
end
