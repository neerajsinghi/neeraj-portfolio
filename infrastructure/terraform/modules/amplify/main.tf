resource "aws_amplify_app" "frontend" {
  name         = var.project
  repository   = var.github_repo
  access_token = var.github_access_token
  platform     = "WEB_COMPUTE"

  # Amplify build spec — frontend lives in /frontend sub-directory of the monorepo
  build_spec = <<-EOT
    version: 1
    applications:
      - frontend:
          phases:
            preBuild:
              commands:
                - cd frontend
                - npm ci
            build:
              commands:
                - npm run build
          artifacts:
            baseDirectory: frontend/.next
            files:
              - '**/*'
          cache:
            paths:
              - frontend/node_modules/**/*
        appRoot: frontend
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
  branch_name = "main"
  framework   = "Next.js - SSR"
  stage       = "PRODUCTION"

  enable_auto_build = false # GitHub Actions triggers the build explicitly

  environment_variables = {
    NEXT_PUBLIC_API_BASE    = var.api_base_url
    NEXT_TELEMETRY_DISABLED = "1"
  }
}
