variable "project_id" {
  description = "GCP project ID"
  type        = string
  default     = "axium-assessment"
}

variable "region" {
  description = "GCP region for resources"
  type        = string
  default     = "us-central1"
}

variable "service_name" {
  description = "Cloud Run service name"
  type        = string
  default     = "event-service"
}

variable "artifact_registry_repository" {
  description = "Artifact Registry repository name"
  type        = string
  default     = "axium-asssessment-image-bucket"
}

variable "image_tag" {
  description = "Container image tag"
  type        = string
  default     = "latest"
}

variable "allow_unauthenticated" {
  description = "Allow unauthenticated access to the Cloud Run service"
  type        = bool
  default     = true
}

variable "env_vars" {
  description = "Environment variables for the Cloud Run service"
  type        = map(string)
  default = {
    ENV                 = "production"
    PROCESSING_DELAY_MS = "1000"
  }
}

variable "min_instances" {
  description = "Minimum number of instances"
  type        = number
  default     = 0
}

variable "max_instances" {
  description = "Maximum number of instances"
  type        = number
  default     = 10
}

variable "cpu_limit" {
  description = "CPU limit for each instance"
  type        = string
  default     = "1"
}

variable "memory_limit" {
  description = "Memory limit for each instance"
  type        = string
  default     = "512Mi"
}
