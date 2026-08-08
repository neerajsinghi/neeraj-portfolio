locals {
  enable_full_stack = var.infrastructure_mode == "full"
  frontend_fqdn = var.enable_custom_domains ? (
    var.frontend_subdomain == "" ? var.root_domain : "${var.frontend_subdomain}.${var.root_domain}"
  ) : ""
  admin_fqdn              = "${var.admin_subdomain}.${var.root_domain}"
  resolved_allowed_origin = var.enable_custom_domains ? "https://${local.frontend_fqdn},https://${local.admin_fqdn}" : var.allowed_origin
}

module "cognito" {
  source       = "./modules/cognito"
  project      = var.project
  aws_region   = var.aws_region
  callback_url = var.enable_custom_domains ? "https://${local.admin_fqdn}/auth/callback" : var.admin_callback_url
  logout_url   = var.enable_custom_domains ? "https://${local.admin_fqdn}" : var.admin_logout_url
}

module "lambda_api" {
  count                = var.enable_lambda_backend ? 1 : 0
  source               = "./modules/lambda_api"
  project              = var.project
  aws_region           = var.aws_region
  lambda_zip_path      = var.lambda_zip_path
  lambda_memory_size   = var.lambda_memory_size
  lambda_timeout       = var.lambda_timeout
  allowed_origin       = local.resolved_allowed_origin
  cognito_user_pool_id = module.cognito.user_pool_id
  cognito_client_id    = module.cognito.client_id
  enable_custom_domain = var.enable_custom_domains
  root_domain          = var.root_domain
  api_subdomain        = var.api_subdomain
}

locals {
  resolved_api_base_url = var.enable_lambda_backend ? module.lambda_api[0].api_base_url : var.api_base_url
}

module "vpc" {
  count       = local.enable_full_stack ? 1 : 0
  source      = "./modules/vpc"
  project     = var.project
  environment = var.environment
  aws_region  = var.aws_region
}

module "ecr" {
  count       = local.enable_full_stack ? 1 : 0
  source      = "./modules/ecr"
  project     = var.project
  environment = var.environment
}

module "eks" {
  count         = local.enable_full_stack ? 1 : 0
  source        = "./modules/eks"
  project       = var.project
  environment   = var.environment
  vpc_id        = module.vpc[0].vpc_id
  subnet_ids    = module.vpc[0].private_subnet_ids
  instance_type = var.eks_node_instance_type
  desired_nodes = var.eks_desired_nodes
  min_nodes     = var.eks_min_nodes
  max_nodes     = var.eks_max_nodes
}

module "amplify" {
  source               = "./modules/amplify"
  project              = var.project
  github_repo          = var.github_repo
  github_access_token  = var.github_access_token
  api_base_url         = local.resolved_api_base_url
  branch_name          = var.amplify_branch
  frontend_app_root    = var.amplify_frontend_app_root
  enable_auto_build    = var.amplify_enable_auto_build
  enable_custom_domain = var.enable_custom_domains
  root_domain          = var.root_domain
  frontend_subdomain   = var.frontend_subdomain
  frontend_enable_www  = var.frontend_enable_www
}

module "admin_amplify" {
  source               = "./modules/amplify"
  project              = "${var.project}-admin"
  github_repo          = var.github_repo
  github_access_token  = var.github_access_token
  api_base_url         = local.resolved_api_base_url
  branch_name          = var.amplify_branch
  frontend_app_root    = var.admin_frontend_app_root
  enable_auto_build    = var.amplify_enable_auto_build
  enable_custom_domain = var.enable_custom_domains
  root_domain          = local.admin_fqdn
  frontend_subdomain   = ""
  frontend_enable_www  = false
  environment_variables = {
    NEXT_PUBLIC_COGNITO_DOMAIN       = module.cognito.domain
    NEXT_PUBLIC_COGNITO_CLIENT_ID    = module.cognito.client_id
    NEXT_PUBLIC_COGNITO_REDIRECT_URI = var.enable_custom_domains ? "https://${local.admin_fqdn}/auth/callback" : var.admin_callback_url
  }
}
