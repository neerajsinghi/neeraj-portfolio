output "ecr_repository_url" {
  description = "ECR URL — set as ECR_REPOSITORY in the backend GitHub Actions workflow"
  value       = module.ecr.repository_url
}

output "eks_cluster_name" {
  description = "EKS cluster name — used in aws eks update-kubeconfig"
  value       = module.eks.cluster_name
}

output "eks_cluster_endpoint" {
  value = module.eks.cluster_endpoint
}

output "amplify_app_id" {
  description = "Amplify app ID — set as AMPLIFY_APP_ID in GitHub Actions secrets"
  value       = module.amplify.app_id
}

output "amplify_app_url" {
  description = "Live frontend URL"
  value       = module.amplify.app_url
}

output "github_actions_role_arn" {
  description = "IAM role ARN — set as AWS_ROLE_ARN in GitHub Actions secrets"
  value       = aws_iam_role.github_actions.arn
}
