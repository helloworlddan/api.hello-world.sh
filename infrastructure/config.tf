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

resource "google_service_account" "terraform" {
  project = local.project
  account_id   = "${local.prefix}-terraform-sa"
  display_name = "${local.prefix} Terraform SA"
}

resource "google_project_iam_binding" "terraform-owner-binding" {
  project = local.project
  role    = "roles/owner"
  members = [
    "serviceAccount:${google_service_account.terraform.email}",
  ]
}

resource "google_organization_iam_binding" "terraform-orgadmin-binding" {
  org_id  = local.organization
  role    = "roles/resourcemanager.organizationAdmin"

  members = [
    "serviceAccount:${google_service_account.terraform.email}",
  ]
}

output "terraform_sa" {
  value = google_service_account.terraform.email
}