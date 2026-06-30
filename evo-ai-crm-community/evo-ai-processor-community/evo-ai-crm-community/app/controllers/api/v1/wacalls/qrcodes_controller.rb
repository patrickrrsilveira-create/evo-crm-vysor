class Api::V1::Wacalls::QrcodesController < Api::BaseController
  before_action :set_channel

  def show
    result = @channel.provider_service.get_qr_code
    
    if result[:success]
      render json: { success: true }
    else
      render json: { error: result[:error] }, status: :unprocessable_entity
    end
  end

  private

  def set_channel
    @channel = Current.account.whatsapp_channels.where(provider: 'wacalls').find_by!("provider_config ->> 'instance_name' = ?", params[:id])
  end
end

