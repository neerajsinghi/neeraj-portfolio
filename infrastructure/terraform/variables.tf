variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "ap-south-1"
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

variable "infrastructure_mode" {
  description = "Infrastructure shape: budget (Amplify-only) or full (EKS + ECR + VPC + Amplify)"
  type        = string
  default     = "budget"

  validation {
    condition     = contains(["budget", "full"], var.infrastructure_mode)
    error_message = "infrastructure_mode must be one of: budget, full"
  }
}

variable "enable_lambda_backend" {
  description = "Deploy backend on Lambda + API Gateway"
  type        = bool
  default     = true
}

variable "lambda_zip_path" {
  description = "Path to the zipped Lambda artifact built from backend/cmd/lambda"
  type        = string
  default     = "../../backend/dist/lambda.zip"
}

variable "lambda_memory_size" {
  description = "Lambda memory in MB"
  type        = number
  default     = 512
}

variable "lambda_timeout" {
  description = "Lambda timeout in seconds"
  type        = number
  default     = 60
}

variable "allowed_origin" {
  description = "CORS allowed origin for backend responses"
  type        = string
  default     = "*"
}

variable "enable_custom_domains" {
  description = "Enable custom domains for Amplify frontend and API Gateway backend"
  type        = bool
  default     = false
}

variable "root_domain" {
  description = "Route53 hosted zone domain (e.g. neerajsinghi.com)"
  type        = string
  default     = ""
}

variable "frontend_subdomain" {
  description = "Frontend subdomain prefix; empty string means apex domain"
  type        = string
  default     = ""
}

variable "frontend_enable_www" {
  description = "Also map www.<domain> to the frontend branch when using apex"
  type        = bool
  default     = true
}

variable "api_subdomain" {
  description = "API subdomain prefix under root_domain"
  type        = string
  default     = "api"
}

variable "admin_subdomain" {
  description = "Admin frontend subdomain prefix"
  type        = string
  default     = "admin"
}

variable "admin_callback_url" {
  description = "Cognito callback URL used when custom domains are disabled"
  type        = string
  default     = "http://localhost:3200/auth/callback"
}

variable "admin_logout_url" {
  description = "Cognito logout URL used when custom domains are disabled"
  type        = string
  default     = "http://localhost:3200"
}

# ── EKS ──────────────────────────────────────────────────────────────────────

variable "eks_node_instance_type" {
  description = "EC2 instance type for EKS worker nodes (used only when infrastructure_mode=full)"
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
  description = "Manual backend API URL when Lambda backend is disabled"
  type        = string
  default     = ""
}

variable "amplify_branch" {
  description = "Git branch Amplify should deploy"
  type        = string
  default     = "main"
}

variable "amplify_frontend_app_root" {
  description = "Path to the Next.js frontend app root inside the repository"
  type        = string
  default     = "frontend"
}

variable "admin_frontend_app_root" {
  description = "Path to the admin Next.js app root inside the repository"
  type        = string
  default     = "admin"
}

variable "amplify_enable_auto_build" {
  description = "Whether Amplify should auto-build on git pushes"
  type        = bool
  default     = false
}

# ── GitHub Actions OIDC ───────────────────────────────────────────────────────

variable "github_org_repo" {
  description = "GitHub org/repo slug for the OIDC trust policy (e.g. neerajsinghi/neeraj-portfolio)"
  type        = string
}
