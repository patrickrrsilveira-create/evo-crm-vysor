puts "AgentBots:"
AgentBot.all.each do |bot|
  puts "ID: #{bot.id} | Name: #{bot.name} | Provider: #{bot.bot_provider}"
end
puts "---"
puts "AiAgents:"
ActiveRecord::Base.connection.execute("SELECT id, name FROM ai_agents").each do |row|
  puts "ID: #{row['id']} | Name: #{row['name']}"
end
