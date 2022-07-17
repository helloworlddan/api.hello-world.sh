resource "google_project_service" "cloudrun" {
  project            = local.project
  service            = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudbuild" {
  project            = local.project
  service            = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "identitytoolkit" {
  project            = local.project
  service            = "identitytoolkit.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "securetoken" {
  project            = local.project
  service            = "securetoken.googleapis.com"
  disable_on_destroy = false
}

# Allow Cloud Build to deploy to Cloud Run
resource "google_project_iam_binding" "cloudbuild-deploy-binding" {
  project = local.project
  role    = "roles/run.admin"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

# Allow Cloud Build to read the Service Management API
resource "google_project_iam_member" "cloudbuild-servicecontroller" {
  project = local.project
  role    = "roles/servicemanagement.serviceController"
  member  = "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
}