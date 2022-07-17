# frozen_string_literal: true

class User
  def initialize(id, email, name)
    @id = id
    @email = email
    @name = name
  end

  attr_reader :id, :email, :name
end

module Context
  def self.authenticate!(request)
    h = request.env
    log = request.logger

    # Protocol verification
    raise AuthZError unless h['HTTP_X_FORWARDED_PROTO']
    raise AuthZError unless h['HTTP_X_FORWARDED_PROTO'] == 'https'

    # Proxy verification
    raise AuthZError unless h['HTTP_AUTHORIZATION']
    proxy = verify_proxy h['HTTP_AUTHORIZATION'].split(' ').last

    # Original user authorization verification
    raise AuthZError unless h['HTTP_X_FORWARDED_AUTHORIZATION']
    client = verify_client h['HTTP_X_FORWARDED_AUTHORIZATION'].split(' ').last

    # Gateway/Proxy user pre-flight authorization verification
    raise AuthZError unless h['HTTP_X_ENDPOINT_API_USERINFO']
    user = extract_user h['HTTP_X_ENDPOINT_API_USERINFO']

    raise AuthZError unless client['user_id'] == user['user_id']

    User.new(user['user_id'], user['email'], user['name'])
  end

  def self.verify_proxy(token)
    raise AuthZError unless token
    raise AuthZError if token.empty?

    result = HTTParty.get "https://oauth2.googleapis.com/tokeninfo?id_token=#{token}"
    raise AuthZError unless result.code == 200

    principal = JSON.parse result.body
    raise AuthZError if principal.key? :error
    raise AuthZError unless principal['iat']
    raise AuthZError unless principal['iss'] == 'https://accounts.google.com'
    raise AuthZError unless principal['email'] == CONFIG[:gateway_sa]

    principal
  end

  def self.verify_client(token)
    raise AuthZError unless token
    raise AuthZError if token.empty?

    decoded = JWT.decode token, nil, false
    payload = decoded[0]
    raise AuthZError unless payload['iat']
    raise AuthZError unless payload['aud']
    raise AuthZError unless payload['aud'] == CONFIG[:project]
    raise AuthZError unless payload['iss']
    raise AuthZError unless payload['iss'] == "https://securetoken.google.com/#{CONFIG[:project]}"
    raise AuthZError unless payload['user_id']
    raise AuthZError unless payload['email']

    result = HTTParty.get 'https://www.googleapis.com/robot/v1/metadata/x509/securetoken@system.gserviceaccount.com'
    raise AuthZError unless result.code == 200

    certs = JSON.parse result.body
    cert = OpenSSL::X509::Certificate.new(certs[decoded[1]['kid']])

    decoded = JWT.decode token, cert.public_key, true, { algorithm: 'RS256' }
    principal = decoded[0]
    raise AuthZError unless principal['iat']
    raise AuthZError unless principal['aud']
    raise AuthZError unless principal['aud'] == CONFIG[:project]
    raise AuthZError unless principal['iss']
    raise AuthZError unless principal['iss'] == "https://securetoken.google.com/#{CONFIG[:project]}"
    raise AuthZError unless principal['user_id']
    raise AuthZError unless principal['email']

    principal
  end

  def self.extract_user(token)
    raise AuthZError unless token
    raise AuthZError if token.empty?

    decoded = Base64.urlsafe_decode64 token

    principal = JSON.parse decoded
    raise AuthZError unless principal['iat']
    raise AuthZError unless principal['aud']
    raise AuthZError unless principal['aud'] == CONFIG[:project]
    raise AuthZError unless principal['iss']
    raise AuthZError unless principal['iss'] == "https://securetoken.google.com/#{CONFIG[:project]}"
    raise AuthZError unless principal['user_id']
    raise AuthZError unless principal['email']
    raise AuthZError unless principal['name']

    principal
  end

  def self.pickup_trace(request); end
end
