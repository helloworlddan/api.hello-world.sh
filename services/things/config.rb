# frozen_string_literal: true

# App
CONFIG = {
  gateway_sa: ENV['GATEWAY_SA'],
  project: ENV['GOOGLE_CLOUD_PROJECT']
}.freeze

# Rack/Sinatra
set :bind, '0.0.0.0'
set :port, ENV['PORT'] || '8080'
set :default_content_type, 'application/json'
set :environment, :production
disable :protection
use Rack::Logger

helpers do
  def log
    request.logger
  end
end
