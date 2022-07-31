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
        env {
          name  = "TOP_SESSION"
          value = local.session
        }
        env {
          name  = "TOP_OWNER"
          value = local.owner
        }
        env {
          name = "TOP_MACHINE"
          value = local.machine
        }
        env {
          name = "TOP_ZONE"
          value = local.zone
        }
      }
    }
    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale"  = "1"
        "client.knative.dev/user-image"     = "gcr.io/hwsh-api/machine"
        "run.googleapis.com/client-name"    = "gcloud"
        "run.googleapis.com/client-version" = "392.0.0"
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

# SA with perms for the top machine
resource "google_service_account" "machine-top" {
  project      = local.project
  account_id   = "${local.prefix}-machine-top"
  display_name = "${local.prefix}-machine-top"
}

resource "google_compute_instance" "machine-top" {
  project      = local.project
  name         = local.machine
  machine_type = "e2-medium"
  zone         = "${local.region}-a"

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
      type  = "pd-ssd"
      size  = "20"
    }
  }

  network_interface {
    network = "default"

    access_config {
      // Ephemeral public IP
    }
  }

  scheduling {
    preemptible = true
    automatic_restart = false
  }

  service_account {
    email  = google_service_account.machine-top.email
    scopes = ["cloud-platform"]
  }
}
