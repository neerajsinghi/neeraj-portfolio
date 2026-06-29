terraform {
  required_version = ">= 1.6.0"

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

  # Bootstrap: create the S3 bucket and DynamoDB table before running init.
  #   aws s3api create-bucket --bucket neeraj-portfolio-tf-state --region us-east-1
  #   aws dynamodb create-table --table-name neeraj-portfolio-tf-lock \
  #     --attribute-definitions AttributeName=LockID,AttributeType=S \
  #     --key-schema AttributeName=LockID,KeyType=HASH \
  #     --billing-mode PAY_PER_REQUEST
  backend "s3" {
    bucket         = "neeraj-portfolio-tf-state"
    key            = "portfolio/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "neeraj-portfolio-tf-lock"
    encrypt        = true
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
