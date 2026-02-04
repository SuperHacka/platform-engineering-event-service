terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  # Optional: Configure backend for remote state storage
  # backend "gcs" {
  #   bucket = "axium-assessment-terraform-state"
  #   prefix = "event-service"
  # }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
