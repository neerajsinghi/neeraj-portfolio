module "vpc" {
  source      = "./modules/vpc"
  project     = var.project
  environment = var.environment
  aws_region  = var.aws_region
}

module "ecr" {
  source      = "./modules/ecr"
  project     = var.project
  environment = var.environment
}

module "eks" {
  source         = "./modules/eks"
  project        = var.project
  environment    = var.environment
  vpc_id         = module.vpc.vpc_id
  subnet_ids     = module.vpc.private_subnet_ids
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
  api_base_url        = var.api_base_url
}
