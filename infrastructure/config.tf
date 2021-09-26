terraform {
  backend "gcs" {
    bucket = "hwsh-api-terraform-state"
    prefix = "terraform/state"
  }
}

provider "google-beta" {
  region = local.region
}

provider "google" {
  region = local.region
}