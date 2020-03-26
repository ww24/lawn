provider "google" {
  # credentials = "${file("../service_account.json")}"
  project = var.project
  region  = var.region
  version = "~> 3.14.0"
}
