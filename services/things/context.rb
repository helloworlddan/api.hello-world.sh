module Context
  def self.authenticate!(request)
    h = request.env

    # Protocol verification
    raise AuthZError unless h['HTTP_X_FORWARDED_PROTO']
    raise AuthZError unless h['HTTP_X_FORWARDED_PROTO'] == 'https'

    # Gateway/Proxy verification
    raise AuthZError unless h['HTTP_AUTHORIZATION']
    token = h['HTTP_AUTHORIZATION'].split(' ').last
    result = HTTParty.get "https://oauth2.googleapis.com/tokeninfo?id_token=#{token}"
    raise AuthZError unless h['HTTP_AUTHORIZATION'] unless result.code == 200
    proxy_principal = JSON.parse(result.body, :symbolize_names => true)
    raise AuthZError if proxy_principal.key? :error
    raise AuthZError unless proxy_principal[:iss] == 'https://accounts.google.com'
    raise AuthZError unless proxy_principal[:email] == CONFIG[:gateway_sa]

    # Original user authorization verification
    raise AuthZError unless h['HTTP_X_FORWARDED_AUTHORIZATION']
    token = h['HTTP_X_FORWARDED_AUTHORIZATION'].split(' ').last
    # TODO verify firebase token
    #user_principal = FirebaseIdToken::Signature.verify(token)
    #raise AuthZError unless user_principal
    #raise AuthZError unless user_principal[:iss] == 'https://securetoken.google.com/firebase-id-token'

    # Gateway/Proxy user pre-flight authorization verification
    raise AuthZError unless h['HTTP_X_ENDPOINT_API_USERINFO']
    firebase_principal = JSON.parse(Base64.urlsafe_decode64(h['HTTP_X_ENDPOINT_API_USERINFO']),:symbolize_names => true)
    #raise AuthZError unless user_principal[:user_id] == firebase_principal[:user_id]

    {
      :id =>  firebase_principal[:user_id],
      :email => firebase_principal[:email],
      :name => firebase_principal[:name],
      :picture => firebase_principal[:picture]
    }
  end

  def self.pickup_trace(request)
  end
end