resource "google_service_account" "app" {
  project      = local.project
  account_id   = "${local.prefix}-app"
  display_name = "${local.prefix}-app"
}

resource "google_service_account_iam_binding" "app-sa-user" {
  service_account_id = google_service_account.app.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

resource google_project_iam_member "app-cloudlogging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.app.email}"
}

resource "google_cloud_run_service" "app" {
  project  = local.project
  provider = google-beta
  name     = "${local.prefix}-app"
  location = local.region
  template {
    spec {
      service_account_name = google_service_account.app.email
      containers {
        image = "gcr.io/${local.project}/app"
        env {
          name = "GOOGLE_CLOUD_PROJECT"
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

resource "google_cloud_run_domain_mapping" "app" {
  project  = local.project
  location = local.region
  name     = "app.${local.domain}"

  metadata {
    namespace = local.project
  }

  spec {
    route_name = google_cloud_run_service.app.name
  }
}

resource "google_cloud_run_service_iam_binding" "app" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.app.name
  role     = "roles/run.invoker"
  members = [
    "allUsers",
  ]
}

resource "google_cloudbuild_trigger" "app" {
  project  = local.project
  provider = google-beta
  github {
    name  = "api.hello-world.sh"
    owner = "helloworlddan"
    push {
      branch = "master"
    }
  }

  name        = "${local.prefix}-app"
  description = "Build pipeline for ${local.prefix} app"

  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "app"
  }

  filename = "services/app/cloudbuild.yaml"
}