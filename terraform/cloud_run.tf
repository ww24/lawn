data "google_cloud_run_service" "lawn" {
  name     = var.name
  location = var.location
}

locals {
  current_image = data.google_cloud_run_service.lawn.template != null ? data.google_cloud_run_service.lawn.template[0].spec[0].containers[0].image : null
  new_image     = "${var.location}-docker.pkg.dev/${var.project}/${var.gar_repository}/${var.image_name}:${var.image_tag}"
  image         = (local.current_image != null && var.image_tag == "latest") ? local.current_image : local.new_image
}

resource "google_cloud_run_service" "lawn" {
  name     = var.name
  location = var.location
  project  = var.project

  template {
    spec {
      service_account_name = var.service_account

      containers {
        image = local.image

        resources {
          limits = {
            cpu    = "1000m"
            memory = "128Mi"
          }
        }

        env {
          name  = "GITHUB_USERNAME"
          value = var.github_username
        }

        env {
          name  = "MAX_AGE"
          value = var.cache_control_max_age
        }
      }
    }

    metadata {
      annotations = {
        "autoscaling.knative.dev/maxScale" = "1"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  lifecycle {
    ignore_changes = [

    ]
  }
}

resource "google_cloud_run_service_iam_policy" "noauth" {
  location = google_cloud_run_service.lawn.location
  project  = google_cloud_run_service.lawn.project
  service  = google_cloud_run_service.lawn.name

  policy_data = data.google_iam_policy.noauth.policy_data
}

data "google_iam_policy" "noauth" {
  binding {
    role = "roles/run.invoker"
    members = [
      "allUsers",
    ]
  }
}
