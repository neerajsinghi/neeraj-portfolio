locals {
  enable_full_stack = var.infrastructure_mode == "full"
}

module "lambda_api" {
  count              = var.enable_lambda_backend ? 1 : 0
  source             = "./modules/lambda_api"
  project            = var.project
  aws_region         = var.aws_region
  lambda_zip_path    = var.lambda_zip_path
  lambda_memory_size = var.lambda_memory_size
  lambda_timeout     = var.lambda_timeout
  allowed_origin     = var.allowed_origin
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
  count          = local.enable_full_stack ? 1 : 0
  source         = "./modules/eks"
  project        = var.project
  environment    = var.environment
  vpc_id         = module.vpc[0].vpc_id
  subnet_ids     = module.vpc[0].private_subnet_ids
  instance_type  = var.eks_node_instance_type
  desired_nodes  = var.eks_desired_nodes
  min_nodes      = var.eks_min_nodes
  max_nodes      = var.eks_max_nodes
}

module "amplify" {
  source              = "./modules/amplify"
  project             = var.project
  github_repo         = var.github_repo
  github_access_token = var.github_access_token
  api_base_url        = local.resolved_api_base_url
  branch_name         = var.amplify_branch
  frontend_app_root   = var.amplify_frontend_app_root
  enable_auto_build   = var.amplify_enable_auto_build
}
