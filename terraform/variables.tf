variable "location" {
  type    = string
  default = "asia-northeast1"
}

variable "project" {
  type = string
}

// credentials json value
variable "google_credentials" {
  type = string
}

variable "name" {
  type    = string
  default = "lawn"
}

variable "gar_repository" {
  type    = string
  default = "ww24"
}

variable "image_name" {
  type    = string
  default = "lawn"
}

variable "image_tag" {
  type    = string
  default = "latest"
}

// cloud run service account
variable "service_account" {
  type = string
}

variable "github_username" {
  type = string
}

variable "cache_control_max_age" {
  type    = string
  default = "6h"
}
