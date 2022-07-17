# Disable authentication check to invoke this service
resource "google_cloud_run_service_iam_binding" "top" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.top.name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}

# SA with perms for this service
resource "google_service_account" "top" {
  project      = local.project
  account_id   = "${local.prefix}-top"
  display_name = "${local.prefix}-top"
}

# Service definition
resource "google_cloud_run_service" "top" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-top"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.top.email
      containers {
        image = "gcr.io/${local.project}/top"
        env {
          name  = "GOOGLE_CLOUD_PROJECT"
          value = local.project
        }
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Custom domain mapping for this service
resource "google_cloud_run_domain_mapping" "top" {
  project  = local.project
  location = local.region
  name     = "top.${local.domain}"
  metadata {
    namespace = local.project
  }
  spec {
    route_name = google_cloud_run_service.top.name
  }
}

# Service build trigger
resource "google_cloudbuild_trigger" "top" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-top"
  description = "Build pipeline for ${local.prefix} top"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "top"
  }

  filename = "services/top/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "top-sa-user" {
  service_account_id = google_service_account.top.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}