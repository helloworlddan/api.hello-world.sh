resource "google_app_engine_application" "database" {
  project       = local.project
  provider      = google-beta
  location_id   = replace(local.region, "1", "")
  database_type = "CLOUD_FIRESTORE"
}