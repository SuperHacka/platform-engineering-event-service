terraform {
  required_version = ">= 1.0"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = "~> 5.0"
    }
  }

  backend "gcs" {
    bucket = "axium-assessment-terraform-state"
    prefix = "event-service"
  }
}

provider "google" {
  project = var.project_id
  region  = var.region
}
