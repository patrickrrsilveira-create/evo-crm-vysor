class Rack::Attack
  # Rate limits
  throttle('req/ip', limit: 300, period: 5.minutes) do |req|
    req.ip # unless req.path.start_with?('/assets')
  end

  # Throttle login attempts by IP
  throttle('logins/ip', limit: 5, period: 60.seconds) do |req|
    if req.path == '/auth/sign_in' && req.post?
      req.ip
    end
  end
end
