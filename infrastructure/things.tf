# Allow Proxy SA to invoke this service
resource "google_cloud_run_service_iam_binding" "things" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.things.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.proxy.email}"
  ]
}

# SA with perms for this service
resource "google_service_account" "things" {
  project      = local.project
  account_id   = "${local.prefix}-things"
  display_name = "${local.prefix}-things"
}
resource "google_project_iam_member" "things-firestore" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.things.email}"
}
resource "google_project_iam_member" "things-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.things.email}"
}

# Service definition
resource "google_cloud_run_service" "things" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-things"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.things.email
      containers {
        image = "gcr.io/${local.project}/things"
        env {
          name  = "GATEWAY_SA"
          value = google_service_account.proxy.email
        }
        env {
          name  = "ENVIRONMENT"
          value = "prod"
        }
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
resource "google_cloudbuild_trigger" "things" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-things"
  description = "Build pipeline for ${local.prefix}-things"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "things"
  }
  filename = "services/things/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "things-sa-user" {
  service_account_id = google_service_account.things.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}