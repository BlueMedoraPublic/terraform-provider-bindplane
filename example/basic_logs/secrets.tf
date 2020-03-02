data "google_secret_manager_secret_version" "bindplane_svc_act" {
  provider = google-beta
  project = "bpcli-dev"
  secret = "bindplane-service-account"
  version = 1
}
