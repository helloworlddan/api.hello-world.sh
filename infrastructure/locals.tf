locals {
  prefix         = "hwsh-api"
  region         = "europe-west4"
  project        = "hwsh-api"
  project_number = "546978254761"
  domain         = "hello-world.sh"
  repo           = "api.hello-world.sh"
  repo_owner     = "helloworlddan"
  branch         = "master"
  organization   = "892444794895"
}

output "project" {
  value = local.project
}

output "region" {
  value = local.region
}

output "prefix" {
  value = local.prefix
}