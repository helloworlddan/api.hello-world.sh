resource "google_service_account" "things" {
  project      = local.project
  account_id   = "${local.prefix}-things"
  display_name = "${local.prefix}-things"
}

resource "google_service_account_iam_binding" "things-sa-user" {
  service_account_id = google_service_account.things.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

resource google_project_iam_member "things-firestore" {
  project = local.project
  role    = "roles/datastore.user"
  member  = "serviceAccount:${google_service_account.things.email}"
}

resource google_project_iam_member "things-cloudtrace" {
  project = local.project
  role    = "roles/cloudtrace.agent"
  member  = "serviceAccount:${google_service_account.things.email}"
}

resource google_project_iam_member "things-cloudlogging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.things.email}"
}

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
          name = "HWSH_GATEWAY_SA"
          value = google_service_account.proxy.email
        }
        env {
          name = "HWSH_ENVIRONMENT"
          value = "prod"
        }
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

resource "google_cloud_run_domain_mapping" "things" {
  project  = local.project
  location = local.region
  name     = "things.api.${local.domain}"

  metadata {
    namespace = local.project
  }

  spec {
    route_name = google_cloud_run_service.things.name
  }
}

resource "google_cloud_run_service_iam_binding" "things" {
  project  = local.project
  location = local.region
  service  = google_cloud_run_service.things.name
  role     = "roles/run.invoker"
  members = [
    "allUsers"
  ]
}

resource "google_cloudbuild_trigger" "things" {
  project  = local.project
  provider = google-beta
  github {
    name  = "api.hello-world.sh"
    owner = "helloworlddan"
    push {
      branch = "master"
    }
  }

  name        = "${local.prefix}-things"
  description = "Build pipeline for ${local.prefix}-things"

  substitutions = {
    _REGION  = local.region
    _PREFIX  = local.prefix
    _SERVICE = "things"
    _HWSH_GATEWAY_SA = google_service_account.proxy.email
    _HWSH_ENVIRONMENT = "prod"
  }

  filename = "services/things/cloudbuild.yaml"
}
