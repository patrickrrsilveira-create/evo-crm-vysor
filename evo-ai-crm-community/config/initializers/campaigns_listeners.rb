Rails.application.config.after_initialize do
  listener = Campaigns::MessageListener.new
  Wisper.subscribe(listener, prefix: 'message')

  Rails.logger.info 'Campaigns::MessageListener registered successfully'
end
