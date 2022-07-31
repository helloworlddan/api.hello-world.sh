# Disable authentication check to invoke this service
resource "google_cloud_run_service_iam_binding" "static" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.static.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.proxy.email}"
  ]
}

# SA with perms for this service
resource "google_service_account" "static" {
  project      = local.project
  account_id   = "${local.prefix}-static"
  display_name = "${local.prefix}-static"
}

# Service definition
resource "google_cloud_run_service" "static" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-static"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.static.email
      containers {
        image = "gcr.io/${local.project}/static"
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

# Service build trigger
resource "google_cloudbuild_trigger" "static" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-static"
  description = "Build pipeline for ${local.prefix} static"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "static"
  }

  filename = "services/static/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "static-sa-user" {
  service_account_id = google_service_account.static.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}