terraform {
  required_version = ">= 1.10.0"

  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.0"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.0"
    }
  }

  # Bootstrap: create the S3 bucket before running init.
  #   aws s3api create-bucket --bucket neeraj-portfolio-tf-state --region ap-south-1
  # Locking uses native S3 conditional writes (use_lockfile) — no DynamoDB table needed.
  backend "s3" {
    bucket       = "neeraj-portfolio-tf-state"
    key          = "portfolio/terraform.tfstate"
    region       = "ap-south-1"
    use_lockfile = true
    encrypt      = true
  }
}

provider "aws" {
  region = var.aws_region

  default_tags {
    tags = {
      Project     = var.project
      Environment = var.environment
      ManagedBy   = "terraform"
    }
  }
}
