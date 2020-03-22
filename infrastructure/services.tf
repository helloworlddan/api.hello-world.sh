resource "google_project_service" "functions" {
  project = local.project
  service = "cloudfunctions.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudrun" {
  project = local.project
  service = "run.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "endpoints" {
  project = local.project
  service = "endpoints.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "registry" {
  project = local.project
  service = "containerregistry.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudbuild" {
  project = local.project
  service = "cloudbuild.googleapis.com"
  disable_on_destroy = false
}