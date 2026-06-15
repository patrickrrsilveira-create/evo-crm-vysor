require_relative 'config/environment'

channel = Channel::Whatsapp.where(provider: 'evolution_go').first
puts "Using channel: #{channel.id}"

api_url = channel.provider_config['api_url']
instance_token = channel.provider_config['instance_token']
phone_number = "5511999999999" # any valid number format

headers = {
  'apikey' => instance_token,
  'Content-Type' => 'application/json'
}

body = {
  number: phone_number,
  state: 'composing',
  isAudio: false
}

puts "Sending to: #{api_url}/message/presence"
puts "Headers: #{headers.keys}"
puts "Body: #{body.to_json}"

response = HTTParty.post(
  "#{api_url}/message/presence",
  headers: headers,
  body: body.to_json
)

puts "Response code: #{response.code}"
puts "Response body: #{response.body}"

