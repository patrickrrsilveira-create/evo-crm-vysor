class Api::V1::UserToursController < Api::V1::BaseController
  before_action :set_user

  def index
    render json: @user.user_tours.all
  end

  def create
    tour = @user.user_tours.find_or_initialize_by(tour_key: tour_params[:tour_key])
    tour.status = tour_params.fetch(:status, 'completed')
    tour.completed_at = Time.current
    tour.save!
    render json: tour, status: :ok
  end

  def destroy
    tour = @user.user_tours.find_by!(tour_key: params[:tour_key])
    tour.destroy!
    render json: { deleted: true }, status: :ok
  end

  private

  def set_user
    @user = current_user
  end

  def tour_params
    params.require(:tour).permit(:tour_key, :status)
  end
end
