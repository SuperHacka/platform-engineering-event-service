# Setting up GitHub Actions CI/CD for Google Cloud Run

This guide explains how to configure Google Cloud and GitHub to enable automated deployments.

## 1. Create a Service Account

Create a service account that GitHub Actions will use to deploy your application:

```bash
# Set your project
gcloud config set project axium-assessment

# Create the service account
gcloud iam service-accounts create github-actions-deployer \
    --display-name="GitHub Actions Deployer"
```

## 2. Assign Required Roles

The service account needs permissions to manage Artifact Registry and Cloud Run:

```bash
# Artifact Registry (to push images)
gcloud projects add-iam-policy-binding axium-assessment \
    --member="serviceAccount:github-actions-deployer@axium-assessment.iam.gserviceaccount.com" \
    --role="roles/artifactregistry.admin"

# Cloud Run (to manage services)
gcloud projects add-iam-policy-binding axium-assessment \
    --member="serviceAccount:github-actions-deployer@axium-assessment.iam.gserviceaccount.com" \
    --role="roles/run.admin"

# Service Account User (to act as the identity for Cloud Run)
gcloud iam service-accounts add-iam-policy-binding \
    github-actions-deployer@axium-assessment.iam.gserviceaccount.com \
    --member="serviceAccount:github-actions-deployer@axium-assessment.iam.gserviceaccount.com" \
    --role="roles/iam.serviceAccountUser"

# Storage Admin (if using GCS for Terraform state)
gcloud projects add-iam-policy-binding axium-assessment \
    --member="serviceAccount:github-actions-deployer@axium-assessment.iam.gserviceaccount.com" \
    --role="roles/storage.admin"

# Service Usage Admin (to enable APIs via Terraform)
gcloud projects add-iam-policy-binding axium-assessment \
    --member="serviceAccount:github-actions-deployer@axium-assessment.iam.gserviceaccount.com" \
    --role="roles/serviceusage.serviceUsageAdmin"
```

## 3. Generate Service Account Key

Create and download the JSON key:

```bash
gcloud iam service-accounts keys create gcp-key.json \
    --iam-account=github-actions-deployer@axium-assessment.iam.gserviceaccount.com
```

## 4. Configure GitHub Secrets

1. Go to your GitHub repository: **Settings** > **Secrets and variables** > **Actions**.
2. Click **New repository secret**.
3. Name: `GCP_SA_KEY`.
4. Value: Paste the entire content of the `gcp-key.json` file.
5. Click **Add secret**.

## 5. (Important) Terraform State

The GitHub Action uses **local state** by default in the workflow run. This is not recommended for production because the state will be lost after the run.

**Recommendation**: Follow the "Remote State" instructions in [terraform/README.md](../terraform/README.md) to use a GCS bucket before enabling the GitHub Action.

## 6. Triggering the Workflow

-   **On Pull Request**: The `test` job will run automatically.
-   **On Push to `main`**: All jobs (`test`, `build-and-push`, `deploy`) will run sequentially to update your live service.
