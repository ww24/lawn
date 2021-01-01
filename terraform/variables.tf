variable "region" {
  default = "asia-northeast1"
}

variable "gcr_region" {
  default = "asia"
}

variable "project" {}

variable "name" {
  default = "lawn"
}

variable "image" {
  default = "ghcr.io/ww24/lawn"
}

// Email address of the IAM service account associated with the revision of the
// service. The service account represents the identity of the running revision,
// and determines what permissions the revision has. If not provided, the
// revision will use the project's default service account.
variable "service_account" {}

variable "github_username" {}

variable "cache_control_max_age" {
  default = "24h"
}

variable "image_tag" {
  default = "latest"
}
