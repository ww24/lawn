variable "location" {
  default = "asia-northeast1"
}

variable "project" {}

// credentials json value
variable "google_credentials" {}

variable "name" {
  default = "lawn"
}

variable "gar_repository" {
  default = "ww24"
}

variable "image_name" {
  default = "lawn"
}

variable "image_tag" {
  default = "latest"
}

// cloud run service account
variable "service_account" {}

variable "github_username" {}

variable "cache_control_max_age" {
  default = "24h"
}
