# frozen_string_literal: true

require 'base64'
require 'json'
require 'sinatra'
require 'httparty'
require 'jwt'

load 'config.rb'
load 'errors.rb'
load 'context.rb'

before do
  @user = Context.authenticate! request
end

get '/things' do
  { user: @user.email }.to_json
end
