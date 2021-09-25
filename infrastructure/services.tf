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

resource "google_project_iam_binding" "cloudbuild-deploy-binding" {
  project = local.project
  role    = "roles/run.admin"
  members = [
    "serviceAccount:${local.project_number}@cloudbuild.gserviceaccount.com"
  ]
}

resource "google_project_service" "cloudprofiler" {
  project            = local.project
  service            = "cloudprofiler.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "cloudtrace" {
  project            = local.project
  service            = "cloudtrace.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "firestore" {
  project            = local.project
  service            = "firestore.googleapis.com"
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

resource "google_project_service" "appengine" {
  project            = local.project
  service            = "appengine.googleapis.com"
  disable_on_destroy = false
}

resource "google_project_service" "iam" {
  project            = local.project
  service            = "iam.googleapis.com"
  disable_on_destroy = false
}
