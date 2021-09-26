# App
CONFIG = {
  :gateway_sa => ENV['GATEWAY_SA'],
  :project => ENV['GOOGLE_CLOUD_PROJECT']
}

# Rack/Sinatra
set :bind, '0.0.0.0'
set :port, ENV['PORT'] || '8080'
set :default_content_type, 'application/json'
set :show_exceptions, :after_handler
disable :protection