# Terraform Configuration for Event Service on Google Cloud Run

This directory contains Terraform configuration to deploy the event-service to Google Cloud Run with all necessary infrastructure.

## Prerequisites

1. **Google Cloud SDK**: Install and authenticate
   ```bash
   gcloud auth application-default login
   gcloud config set project axium-assessment
   ```

2. **Terraform**: Install Terraform >= 1.0
   ```bash
   # Check version
   terraform version
   ```

3. **Docker Image**: Build and push your Docker image to Artifact Registry first (see deployment steps below)

4. **GCP APIs**: The configuration will automatically enable required APIs:
   - Cloud Run API
   - Artifact Registry API

## Project Configuration

- **Project ID**: `axium-assessment`
- **Region**: `us-central1` (configurable)
- **Artifact Registry Repository**: `axium-asssessment-image-bucket`
- **Service Name**: `event-service`

## Deployment Steps

### 1. Build and Push Docker Image

First, authenticate Docker with Artifact Registry:
```bash
gcloud auth configure-docker us-central1-docker.pkg.dev
```

Build and tag your image:
```bash
# From project root directory
docker build -t event-service .

# Tag for Artifact Registry
docker tag event-service us-central1-docker.pkg.dev/axium-assessment/axium-asssessment-image-bucket/event-service:latest
```

**Note**: You'll need to manually create the Artifact Registry repository first, or run `terraform apply` and let it fail on the Cloud Run deployment, then push your image and re-run `terraform apply`.

Alternatively, create the repository manually:
```bash
gcloud artifacts repositories create axium-asssessment-image-bucket \
  --repository-format=docker \
  --location=us-central1 \
  --description="Docker repository for event-service"
```

Then push the image:
```bash
docker push us-central1-docker.pkg.dev/axium-assessment/axium-asssessment-image-bucket/event-service:latest
```

### 2. Initialize Terraform

```bash
cd terraform
terraform init
```

### 3. Review the Plan

```bash
terraform plan
```

This will show you what resources will be created:
- Artifact Registry repository
- Cloud Run service
- IAM bindings for public access

### 4. Apply the Configuration

```bash
terraform apply
```

Review the plan and type `yes` to confirm.

### 5. Access Your Service

After deployment completes, Terraform will output the service URL:
```bash
terraform output service_url
```

Test the service:
```bash
curl $(terraform output -raw service_url)/health
```

## Configuration Options

### Customize Variables

You can override default variables by creating a `terraform.tfvars` file:

```bash
cp terraform.tfvars.example terraform.tfvars
# Edit terraform.tfvars with your values
```

Or pass variables via command line:
```bash
terraform apply -var="image_tag=v1.0.0" -var="max_instances=20"
```

### Available Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `project_id` | `axium-assessment` | GCP project ID |
| `region` | `us-central1` | GCP region |
| `service_name` | `event-service` | Cloud Run service name |
| `artifact_registry_repository` | `axium-asssessment-image-bucket` | Artifact Registry repository name |
| `image_tag` | `latest` | Container image tag |
| `allow_unauthenticated` | `true` | Allow public access |
| `min_instances` | `0` | Minimum instances (scale to zero) |
| `max_instances` | `10` | Maximum instances |
| `cpu_limit` | `1` | CPU cores per instance |
| `memory_limit` | `512Mi` | Memory per instance |

### Environment Variables

Configure application environment variables in `terraform.tfvars`:
```hcl
env_vars = {
  ENV                 = "production"
  PROCESSING_DELAY_MS = "500"
}
```

## Updating the Service

To deploy a new version:

1. Build and push new image with a new tag:
   ```bash
   docker build -t event-service .
   docker tag event-service us-central1-docker.pkg.dev/axium-assessment/axium-asssessment-image-bucket/event-service:v2
   docker push us-central1-docker.pkg.dev/axium-assessment/axium-asssessment-image-bucket/event-service:v2
   ```

2. Update the image tag:
   ```bash
   terraform apply -var="image_tag=v2"
   ```

## Tearing Down

To destroy all resources:
```bash
terraform destroy
```

## Troubleshooting

### Image Not Found Error

If Cloud Run can't find your image:
1. Verify the image exists in Artifact Registry:
   ```bash
   gcloud artifacts docker images list us-central1-docker.pkg.dev/axium-assessment/axium-asssessment-image-bucket
   ```

2. Ensure you've pushed the image with the correct tag matching `var.image_tag`

### Permission Denied

Ensure your account has the necessary IAM roles:
- `roles/run.admin` - Cloud Run Admin
- `roles/artifactregistry.admin` - Artifact Registry Admin
- `roles/iam.serviceAccountUser` - Service Account User

### API Not Enabled

If you get API errors, manually enable required APIs:
```bash
gcloud services enable run.googleapis.com
gcloud services enable artifactregistry.googleapis.com
```

## Remote State (Optional)

For team collaboration, store Terraform state in GCS:

1. Create a bucket for state:
   ```bash
   gcloud storage buckets create gs://axium-assessment-terraform-state --location=us-central1
   ```

2. Uncomment the backend configuration in `provider.tf`:
   ```hcl
   backend "gcs" {
     bucket = "axium-assessment-terraform-state"
     prefix = "event-service"
   }
   ```

3. Re-initialize:
   ```bash
   terraform init -migrate-state
   ```
