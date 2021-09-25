resource "google_identity_platform_default_supported_idp_config" "google" {
  project       = local.project
  enabled       = true
  idp_id        = "google.com"
  client_id     = lookup(lookup(local.idps, "google", "undefined"), "client_id", "undefined")
  client_secret = lookup(lookup(local.idps, "google", "undefined"), "client_secret", "undefined")
}

resource "google_identity_platform_default_supported_idp_config" "github" {
  project       = local.project
  enabled       = true
  idp_id        = "github.com"
  client_id     = lookup(lookup(local.idps, "github", "undefined"), "client_id", "undefined")
  client_secret = lookup(lookup(local.idps, "github", "undefined"), "client_secret", "undefined")
}
