module AgentBots
  class WebhookJob < ::WebhookJob
    queue_as :high
  end
end
