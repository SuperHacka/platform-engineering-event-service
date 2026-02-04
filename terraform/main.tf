# Enable required APIs
resource "google_project_service" "cloud_run" {
  project = var.project_id
  service = "run.googleapis.com"

  disable_on_destroy = false
}

resource "google_project_service" "artifact_registry" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"

  disable_on_destroy = false
}

# Create Artifact Registry repository for Docker images
resource "google_artifact_registry_repository" "event_service" {
  location      = var.region
  repository_id = var.artifact_registry_repository
  description   = "Docker repository for event-service container images"
  format        = "DOCKER"

  depends_on = [google_project_service.artifact_registry]
}

# Cloud Run service
resource "google_cloud_run_service" "event_service" {
  name     = var.service_name
  location = var.region

  template {
    spec {
      containers {
        # Image must be pushed to Artifact Registry first
        # Format: REGION-docker.pkg.dev/PROJECT_ID/REPOSITORY/IMAGE:TAG
        image = "${var.region}-docker.pkg.dev/${var.project_id}/${var.artifact_registry_repository}/event-service:${var.image_tag}"

        # Environment variables
        dynamic "env" {
          for_each = var.env_vars
          content {
            name  = env.key
            value = env.value
          }
        }

        # Resource limits
        resources {
          limits = {
            cpu    = var.cpu_limit
            memory = var.memory_limit
          }
        }

        # Security: Run as non-root user (matches Dockerfile)
        # Cloud Run runs containers as UID 1000 by default
      }

      # Container concurrency (number of requests per container)
      container_concurrency = 80

      # Timeout for requests
      timeout_seconds = 300
    }

    metadata {
      annotations = {
        # Auto-scaling configuration
        "autoscaling.knative.dev/minScale" = var.min_instances
        "autoscaling.knative.dev/maxScale" = var.max_instances

        # Use second generation execution environment
        "run.googleapis.com/execution-environment" = "gen2"
      }
    }
  }

  traffic {
    percent         = 100
    latest_revision = true
  }

  depends_on = [
    google_project_service.cloud_run,
    google_artifact_registry_repository.event_service
  ]
}

# IAM policy to allow public access (if enabled)
resource "google_cloud_run_service_iam_member" "public_access" {
  count = var.allow_unauthenticated ? 1 : 0

  service  = google_cloud_run_service.event_service.name
  location = google_cloud_run_service.event_service.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
