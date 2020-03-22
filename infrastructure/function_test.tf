data "archive_file" "test" {
 type = "zip"
 source_dir = "function_test"
 output_path = "../test.zip"
}

resource "google_storage_bucket_object" "test" {
  name = "artifacts/functions/test.zip"
  bucket = "${local.prefix}-admin"
  source = "../test.zip"
}

resource "google_cloudfunctions_function" "test" {
  project = local.project
  name = "${local.prefix}-test"
  available_memory_mb = 128
  source_archive_bucket ="${local.prefix}-admin"
  source_archive_object = google_storage_bucket_object.test.name
  timeout = 10
  entry_point = "Test"
  trigger_http = true
  runtime = "go113"

  labels = {
    path = "test"
  }

  provisioner "local-exec" {
    command = "gcloud alpha functions remove-iam-policy-binding ${google_cloudfunctions_function.test.name} --member=allUsers --role=roles/cloudfunctions.invoker --region=${local.region}"
  }
}

resource "google_cloudfunctions_function_iam_member" "invoker" {
  project        = local.project
  region         = local.region
  cloud_function = google_cloudfunctions_function.test.name

  role   = "roles/cloudfunctions.invoker"
  member = "domain:hello-world.sh"
}