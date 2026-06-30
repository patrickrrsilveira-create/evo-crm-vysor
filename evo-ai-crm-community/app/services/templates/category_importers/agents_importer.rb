# frozen_string_literal: true

module Templates
  module CategoryImporters
    # Creates AgentBots with secrets forcibly nil. User must reconfigure
    # api_key and outgoing_url in the agent settings before the bot will work.
    class AgentsImporter < Base
      CATEGORY = 'agents'
      MODEL = ::AgentBot
      UNIQUE_FIELD = :name

      private

      def attributes_for(item)
        attrs = item.except('slug')
        
        # Map unknown attributes from Evo-Nexus exports into bot_config
        known_attrs = %w[name description outgoing_url bot_type bot_config api_key bot_provider message_signature text_segmentation_enabled text_segmentation_limit text_segmentation_min_size delay_per_character debounce_time id created_at updated_at]
        
        bot_config = attrs['bot_config'] || {}
        bot_config = bot_config.dup if bot_config.is_a?(Hash)
        
        attrs.keys.each do |key|
          next if known_attrs.include?(key)
          bot_config[key] = attrs.delete(key)
        end
        attrs['bot_config'] = bot_config

        # Defense in depth: zero secrets even if Sanitizer.zero_blocked_fields! missed.
        attrs['api_key'] = "configure-#{SecureRandom.hex(4)}"
        attrs['outgoing_url'] = nil
        attrs
      end
    end
  end
end
