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
  machine        = "top"
  zone           = "${local.region}-a"
  owner          = "dan@hello-world.sh"
  session        = "ca683f00-d51c-4f1a-af5e-5f9a25b3f4a8"
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
