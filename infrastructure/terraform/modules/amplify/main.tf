resource "aws_amplify_app" "frontend" {
  name         = var.project
  repository   = var.github_repo
  access_token = var.github_access_token
  platform     = "WEB_COMPUTE"

  # Amplify build spec for a monorepo Next.js app.
  # appRoot scopes build commands, artifact paths, and cache paths.
  build_spec = <<-EOT
    version: 1
    applications:
      - frontend:
          appRoot: ${var.frontend_app_root}
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
  EOT

  environment_variables = {
    NEXT_PUBLIC_API_BASE    = var.api_base_url
    NEXT_TELEMETRY_DISABLED = "1"
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

  environment_variables = {
    NEXT_PUBLIC_API_BASE    = var.api_base_url
    NEXT_TELEMETRY_DISABLED = "1"
  }
}
