resource "google_service_account" "proxy" {
  project      = local.project
  account_id   = "${local.prefix}-proxy"
  display_name = "${local.prefix}-proxy"
}

resource "google_service_account_iam_binding" "proxy-sa-user" {
  service_account_id = google_service_account.proxy.name
  role               = "roles/iam.serviceAccountUser"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

resource google_project_iam_member "proxy-cloudlogging" {
  project = local.project
  role    = "roles/logging.logWriter"
  member  = "serviceAccount:${google_service_account.proxy.email}"
}

resource "google_cloud_run_service" "proxy" {
  provider = google-beta
  project  = local.project
  name     = "${local.prefix}-proxy"
  location = local.region

  template {
    spec {
      service_account_name = google_service_account.proxy.email
      containers {
        image = "gcr.io/${local.project}/proxy"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }
}

resource "google_cloud_run_domain_mapping" "proxy" {
  project  = local.project
  location = local.region
  name     = "api.${local.domain}"

  metadata {
    namespace = local.project
  }

  spec {
    route_name = google_cloud_run_service.proxy.name
  }
}

resource "google_cloud_run_service_iam_binding" "public" {
  location = local.region
  project  = local.project
  service  = google_cloud_run_service.proxy.name
  role     = "roles/run.invoker"
  members = [
    "allUsers", # TODO Change to allAuthenticatedUsers
  ]
}

