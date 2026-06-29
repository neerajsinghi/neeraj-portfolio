variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1"
}

variable "project" {
  description = "Project slug used in resource names and tags"
  type        = string
  default     = "neeraj-portfolio"
}

variable "environment" {
  description = "Deployment environment (prod, staging, dev)"
  type        = string
  default     = "prod"
}

# ── EKS ──────────────────────────────────────────────────────────────────────

variable "eks_node_instance_type" {
  description = "EC2 instance type for EKS worker nodes"
  type        = string
  default     = "t3.small"
}

variable "eks_desired_nodes" {
  type    = number
  default = 2
}

variable "eks_min_nodes" {
  type    = number
  default = 1
}

variable "eks_max_nodes" {
  type    = number
  default = 4
}

# ── Amplify ───────────────────────────────────────────────────────────────────

variable "github_repo" {
  description = "Full GitHub HTTPS URL of the monorepo (e.g. https://github.com/neerajsinghi/neeraj-portfolio)"
  type        = string
}

variable "github_access_token" {
  description = "GitHub personal access token with repo scope — used by Amplify to clone the repo"
  type        = string
  sensitive   = true
}

variable "api_base_url" {
  description = "Public URL of the backend API injected into the Amplify build as NEXT_PUBLIC_API_BASE"
  type        = string
}

# ── GitHub Actions OIDC ───────────────────────────────────────────────────────

variable "github_org_repo" {
  description = "GitHub org/repo slug for the OIDC trust policy (e.g. neerajsinghi/neeraj-portfolio)"
  type        = string
}
