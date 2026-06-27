class Api::V1::Wacalls::SessionsController < Api::BaseController
  before_action :set_channel

  def logout
    success = @channel.provider_service.logout
    
    if success
      render json: { success: true }
    else
      render json: { error: 'Failed to disconnect session' }, status: :unprocessable_entity
    end
  end

  private

  def set_channel
    @channel = Current.account.whatsapp_channels.where(provider: 'wacalls').find_by!("provider_config ->> 'instance_name' = ?", params[:id])
  end
end
