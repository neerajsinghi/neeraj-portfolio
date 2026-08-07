output "ecr_repository_url" {
  description = "ECR URL (only when infrastructure_mode=full)"
  value       = var.infrastructure_mode == "full" ? module.ecr[0].repository_url : null
}

output "eks_cluster_name" {
  description = "EKS cluster name (only when infrastructure_mode=full)"
  value       = var.infrastructure_mode == "full" ? module.eks[0].cluster_name : null
}

output "eks_cluster_endpoint" {
  value = var.infrastructure_mode == "full" ? module.eks[0].cluster_endpoint : null
}

output "lambda_function_name" {
  description = "Lambda backend function name"
  value       = var.enable_lambda_backend ? module.lambda_api[0].function_name : null
}

output "backend_api_base_url" {
  description = "Backend API base URL used by frontend"
  value       = var.enable_lambda_backend ? module.lambda_api[0].api_base_url : var.api_base_url
}

output "backend_api_custom_domain" {
  description = "Backend API custom domain (if enabled)"
  value       = var.enable_lambda_backend ? module.lambda_api[0].api_custom_domain : null
}

output "amplify_app_id" {
  description = "Amplify app ID — set as AMPLIFY_APP_ID in GitHub Actions secrets"
  value       = module.amplify.app_id
}

output "amplify_app_url" {
  description = "Live frontend URL"
  value       = module.amplify.branch_web_url
}

output "frontend_custom_domain" {
  description = "Frontend custom domain host (if enabled)"
  value       = module.amplify.frontend_fqdn
}

output "resolved_allowed_origin" {
  description = "Allowed origin value wired into backend"
  value       = local.resolved_allowed_origin
}

output "amplify_default_domain" {
  description = "Amplify default domain"
  value       = module.amplify.app_default_domain
}

output "amplify_branch" {
  description = "Amplify branch currently managed by Terraform"
  value       = module.amplify.branch_name
}

output "github_actions_role_arn" {
  description = "IAM role ARN — set as AWS_ROLE_ARN in GitHub Actions secrets"
  value       = aws_iam_role.github_actions.arn
}
