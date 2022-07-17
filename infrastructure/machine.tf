# Allow Proxy SA to invoke this service
resource "google_cloud_run_service_iam_binding" "machine" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.machine.name
  role     = "roles/run.invoker"
  members = [
    "serviceAccount:${google_service_account.proxy.email}"
  ]
}

# SA with perms for this service
resource "google_service_account" "machine" {
  project      = local.project
  account_id   = "${local.prefix}-machine"
  display_name = "${local.prefix}-machine"
}

resource "google_project_iam_member" "machine-compute" {
  project = local.project
  role    = "roles/compute.admin"
  member  = "serviceAccount:${google_service_account.machine.email}"
}

# Service definition
resource "google_cloud_run_service" "machine" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-machine"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.machine.email
      containers {
        image = "gcr.io/${local.project}/machine"
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
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"      = "1"
      }
    }
  }
  traffic {
    percent         = 100
    latest_revision = true
  }
}

# Service build trigger
resource "google_cloudbuild_trigger" "machine" {
  project  = local.project
  provider = google-beta
  github {
    name  = local.repo
    owner = local.repo_owner
    push {
      branch = local.branch
    }
  }
  name        = "${local.prefix}-machine"
  description = "Build pipeline for ${local.prefix}-machine"
  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "machine"
  }
  filename = "services/machine/cloudbuild.yaml"
}

# Allow Cloud Build to bind SA
resource "google_service_account_iam_binding" "machine-sa-user" {
  service_account_id = google_service_account.machine.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}