resource "aws_amplify_app" "frontend" {
  name         = var.project
  repository   = var.github_repo
  access_token = var.github_access_token
  platform     = "WEB_COMPUTE"

  # Amplify build spec for a monorepo Next.js app.
  # appRoot must be a sibling of "frontend:" (same indent level), not nested inside it.
  build_spec = <<-EOT
    version: 1
    applications:
      - frontend:
          phases:
            preBuild:
              commands:
                - npm ci
            build:
              commands:
                - npm run build
          artifacts:
            baseDirectory: .next
            files:
              - '**/*'
          cache:
            paths:
              - node_modules/**/*
        appRoot: ${var.frontend_app_root}
  EOT

  environment_variables = merge({
    NEXT_PUBLIC_API_BASE      = var.api_base_url
    NEXT_TELEMETRY_DISABLED   = "1"
    AMPLIFY_MONOREPO_APP_ROOT = var.frontend_app_root
  }, var.environment_variables)

  dynamic "custom_rule" {
    for_each = var.enable_custom_domain && var.frontend_enable_www && var.frontend_subdomain == "" ? [1] : []
    content {
      source = "https://www.${var.root_domain}"
      target = "https://${var.root_domain}"
      status = "301"
    }
  }

  # Ignore changes to access_token after initial creation
  lifecycle {
    ignore_changes = [access_token]
  }
}

resource "aws_amplify_branch" "main" {
  app_id      = aws_amplify_app.frontend.id
  branch_name = var.branch_name
  framework   = "Next.js - SSR"
  stage       = "PRODUCTION"

  enable_auto_build = var.enable_auto_build

  environment_variables = merge({
    NEXT_PUBLIC_API_BASE      = var.api_base_url
    NEXT_TELEMETRY_DISABLED   = "1"
    AMPLIFY_MONOREPO_APP_ROOT = var.frontend_app_root
  }, var.environment_variables)
}

resource "aws_amplify_domain_association" "frontend" {
  count       = var.enable_custom_domain ? 1 : 0
  app_id      = aws_amplify_app.frontend.id
  domain_name = var.root_domain

  sub_domain {
    branch_name = aws_amplify_branch.main.branch_name
    prefix      = var.frontend_subdomain
  }

  dynamic "sub_domain" {
    for_each = var.frontend_enable_www && var.frontend_subdomain == "" ? [1] : []
    content {
      branch_name = aws_amplify_branch.main.branch_name
      prefix      = "www"
    }
  }
}
