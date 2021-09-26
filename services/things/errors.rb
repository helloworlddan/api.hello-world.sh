class APIError < StandardError
  def initialize(msg="unknown server error", code=500)
    super(msg)
    @code = code
  end
  attr_reader :code
end

class AuthZError < APIError
  def initialize(msg="user could not be authorized", code=401)
    super(msg, code)
  end
end

error APIError do
  e = env['sinatra.error']
  status e.code
  {:message => e.message}.to_json
end

error do
  status 500
  {:message => 'internal server error'}.to_json
end