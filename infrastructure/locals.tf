locals {
  prefix         = "hwsh-api"
  region         = "europe-west1"
  project        = "hwsh-api"
  project_number = "546978254761"
  domain         = "hello-world.sh"
  support_email  = "api@hello-world.sh"
  organization   = "892444794895"
  idps = {
    google = {
      client_id     = ""
      client_secret = ""
    }
    github = {
      client_id     = ""
      client_secret = ""
    }
  }
}