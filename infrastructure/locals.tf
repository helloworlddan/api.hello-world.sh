locals {
  prefix = "hwsh-api"
  region = "europe-west1"
  project = lookup(data.external.project.result, "project", "null")
  domain = "api.hello-world.sh"
}