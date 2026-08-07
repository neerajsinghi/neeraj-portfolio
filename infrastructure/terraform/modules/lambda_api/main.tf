data "aws_caller_identity" "current" {}

data "aws_ssm_parameters_by_path" "backend_env" {
  path            = "/${var.project}"
  recursive       = false
  with_decryption = true
}

data "aws_route53_zone" "root" {
  count        = var.enable_custom_domain ? 1 : 0
  name         = "${var.root_domain}."
  private_zone = false
}

locals {
  api_custom_fqdn = "${var.api_subdomain}.${var.root_domain}"
  ssm_values_by_key = {
    for idx, name in data.aws_ssm_parameters_by_path.backend_env.names :
    trimprefix(name, "/${var.project}/") => data.aws_ssm_parameters_by_path.backend_env.values[idx]
  }
  cert_validation_options = var.enable_custom_domain ? {
    for dvo in aws_acm_certificate.api[0].domain_validation_options : dvo.domain_name => {
      name   = dvo.resource_record_name
      record = dvo.resource_record_value
      type   = dvo.resource_record_type
    }
  } : {}
}

resource "aws_iam_role" "lambda_exec" {
  name = "${var.project}-lambda-exec"

  assume_role_policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect = "Allow"
      Principal = {
        Service = "lambda.amazonaws.com"
      }
      Action = "sts:AssumeRole"
    }]
  })
}

resource "aws_iam_role_policy" "lambda_inline" {
  name = "${var.project}-lambda-inline"
  role = aws_iam_role.lambda_exec.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Effect = "Allow"
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ]
        Resource = "arn:aws:logs:${var.aws_region}:${data.aws_caller_identity.current.account_id}:*"
      },
      {
        Effect = "Allow"
        Action = [
          "ssm:GetParameter",
          "ssm:GetParameters",
          "ssm:GetParametersByPath"
        ]
        Resource = "arn:aws:ssm:${var.aws_region}:${data.aws_caller_identity.current.account_id}:parameter/${var.project}/*"
      },
      {
        Effect = "Allow"
        Action = [
          "kms:Decrypt",
          "kms:GenerateDataKey"
        ]
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

resource "aws_lambda_function" "backend" {
  function_name = "${var.project}-backend-api"
  role          = aws_iam_role.lambda_exec.arn
  runtime       = "provided.al2023"
  handler       = "bootstrap"
  architectures = ["arm64"]

  filename         = var.lambda_zip_path
  source_code_hash = filebase64sha256(var.lambda_zip_path)

  timeout     = var.lambda_timeout
  memory_size = var.lambda_memory_size

  environment {
    variables = {
      ANTHROPIC_API_KEY = lookup(local.ssm_values_by_key, "anthropic-api-key", "")
      OPENAI_API_KEY    = lookup(local.ssm_values_by_key, "openai-api-key", "")
      XAI_API_KEY       = lookup(local.ssm_values_by_key, "xai-api-key", "")
      GEMINI_API_KEY    = lookup(local.ssm_values_by_key, "gemini-api-key", "")
      ANTHROPIC_MODEL   = lookup(local.ssm_values_by_key, "anthropic-model", "claude-sonnet-4-6")
      OPENAI_MODEL      = lookup(local.ssm_values_by_key, "openai-model", "gpt-4o")
      GROK_MODEL        = lookup(local.ssm_values_by_key, "grok-model", "grok-3")
      GEMINI_MODEL      = lookup(local.ssm_values_by_key, "gemini-model", "gemini-2.0-flash")
      GITHUB_USER       = lookup(local.ssm_values_by_key, "github-user", "neerajsinghi")
      GITHUB_TOKEN      = lookup(local.ssm_values_by_key, "github-token", "")
      PORT              = lookup(local.ssm_values_by_key, "port", "8080")
      ALLOWED_ORIGIN    = lookup(local.ssm_values_by_key, "allowed-origin", var.allowed_origin)
    }
  }
}

resource "aws_apigatewayv2_api" "backend" {
  name          = "${var.project}-backend-api"
  protocol_type = "HTTP"
}

resource "aws_apigatewayv2_integration" "backend" {
  api_id                 = aws_apigatewayv2_api.backend.id
  integration_type       = "AWS_PROXY"
  integration_uri        = aws_lambda_function.backend.invoke_arn
  payload_format_version = "2.0"
}

resource "aws_apigatewayv2_route" "root" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "ANY /"
  target    = "integrations/${aws_apigatewayv2_integration.backend.id}"
}

resource "aws_apigatewayv2_route" "proxy" {
  api_id    = aws_apigatewayv2_api.backend.id
  route_key = "ANY /{proxy+}"
  target    = "integrations/${aws_apigatewayv2_integration.backend.id}"
}

resource "aws_apigatewayv2_stage" "default" {
  api_id      = aws_apigatewayv2_api.backend.id
  name        = "$default"
  auto_deploy = true
}

resource "aws_acm_certificate" "api" {
  count             = var.enable_custom_domain ? 1 : 0
  domain_name       = local.api_custom_fqdn
  validation_method = "DNS"
}

resource "aws_route53_record" "api_cert_validation" {
  for_each = local.cert_validation_options

  zone_id = data.aws_route53_zone.root[0].zone_id
  name    = each.value.name
  type    = each.value.type
  ttl     = 60
  records = [each.value.record]
}

resource "aws_acm_certificate_validation" "api" {
  count                   = var.enable_custom_domain ? 1 : 0
  certificate_arn         = aws_acm_certificate.api[0].arn
  validation_record_fqdns = [for rec in aws_route53_record.api_cert_validation : rec.fqdn]
}

resource "aws_apigatewayv2_domain_name" "api" {
  count       = var.enable_custom_domain ? 1 : 0
  domain_name = local.api_custom_fqdn

  domain_name_configuration {
    certificate_arn = aws_acm_certificate_validation.api[0].certificate_arn
    endpoint_type   = "REGIONAL"
    security_policy = "TLS_1_2"
  }
}

resource "aws_apigatewayv2_api_mapping" "api" {
  count       = var.enable_custom_domain ? 1 : 0
  api_id      = aws_apigatewayv2_api.backend.id
  domain_name = aws_apigatewayv2_domain_name.api[0].id
  stage       = aws_apigatewayv2_stage.default.name
}

resource "aws_route53_record" "api_alias_a" {
  count   = var.enable_custom_domain ? 1 : 0
  zone_id = data.aws_route53_zone.root[0].zone_id
  name    = local.api_custom_fqdn
  type    = "A"

  alias {
    name                   = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_route53_record" "api_alias_aaaa" {
  count   = var.enable_custom_domain ? 1 : 0
  zone_id = data.aws_route53_zone.root[0].zone_id
  name    = local.api_custom_fqdn
  type    = "AAAA"

  alias {
    name                   = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].target_domain_name
    zone_id                = aws_apigatewayv2_domain_name.api[0].domain_name_configuration[0].hosted_zone_id
    evaluate_target_health = false
  }
}

resource "aws_lambda_permission" "apigw_invoke" {
  statement_id  = "AllowExecutionFromAPIGateway"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.backend.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "${aws_apigatewayv2_api.backend.execution_arn}/*/*"
}
