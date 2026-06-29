# GitHub Actions authenticates to AWS via OIDC — no long-lived keys stored in secrets.
# Only AWS_ROLE_ARN needs to be stored as a GitHub secret.

data "aws_caller_identity" "current" {}

resource "aws_iam_openid_connect_provider" "github_actions" {
  url             = "https://token.actions.githubusercontent.com"
  client_id_list  = ["sts.amazonaws.com"]
  thumbprint_list = ["6938fd4d98bab03faadb97b34396831e3780aea1"]
}

resource "aws_iam_role" "github_actions" {
  name = "${var.project}-github-actions"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Federated = aws_iam_openid_connect_provider.github_actions.arn
      }
      Action = "sts:AssumeRoleWithWebIdentity"
      Condition = {
        StringEquals = {
          "token.actions.githubusercontent.com:aud" = "sts.amazonaws.com"
        }
        StringLike = {
          "token.actions.githubusercontent.com:sub" = "repo:${var.github_org_repo}:*"
        }
      }
    }]
  })
}

# ── ECR push/pull ─────────────────────────────────────────────────────────────

resource "aws_iam_policy" "github_ecr" {
  name = "${var.project}-github-ecr"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Action = [
        "ecr:GetAuthorizationToken",
        "ecr:BatchCheckLayerAvailability",
        "ecr:GetDownloadUrlForLayer",
        "ecr:BatchGetImage",
        "ecr:InitiateLayerUpload",
        "ecr:UploadLayerPart",
        "ecr:CompleteLayerUpload",
        "ecr:PutImage",
      ]
      Resource = "*"
    }]
  })
}

# ── EKS describe (kubeconfig) ─────────────────────────────────────────────────

resource "aws_iam_policy" "github_eks" {
  name = "${var.project}-github-eks"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["eks:DescribeCluster", "eks:ListClusters"]
      Resource = "*"
    }]
  })
}

# ── Amplify deploy ────────────────────────────────────────────────────────────

resource "aws_iam_policy" "github_amplify" {
  name = "${var.project}-github-amplify"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect   = "Allow"
      Action   = ["amplify:StartJob", "amplify:GetJob"]
      Resource = "*"
    }]
  })
}

# ── SSM Parameter Store (bootstrap write + deploy read) ──────────────────────

resource "aws_iam_policy" "github_ssm" {
  name = "${var.project}-github-ssm"
  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "SSMParams"
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath",
          "ssm:PutParameter",
          "ssm:DeleteParameter",
          "ssm:DescribeParameters",
        ]
        Resource = "arn:aws:ssm:*:${data.aws_caller_identity.current.account_id}:parameter/${var.project}/*"
      },
      {
        # Allow use of the aws/ssm managed CMK for SecureString en/decryption
        Sid    = "SSMKMSDecrypt"
        Effect = "Allow"
        Action = ["kms:Decrypt", "kms:GenerateDataKey"]
        Resource = "*"
        Condition = {
          StringEquals = {
            "kms:ViaService" = "ssm.${var.aws_region}.amazonaws.com"
          }
        }
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "github_ecr" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_ecr.arn
}

resource "aws_iam_role_policy_attachment" "github_eks" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_eks.arn
}

resource "aws_iam_role_policy_attachment" "github_amplify" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_amplify.arn
}

resource "aws_iam_role_policy_attachment" "github_ssm" {
  role       = aws_iam_role.github_actions.name
  policy_arn = aws_iam_policy.github_ssm.arn
}
