output "function_name" {
  value = aws_lambda_function.backend.function_name
}

output "function_arn" {
  value = aws_lambda_function.backend.arn
}

output "api_base_url" {
  value = var.enable_custom_domain ? "https://${local.api_custom_fqdn}" : aws_apigatewayv2_stage.default.invoke_url
}

output "api_custom_domain" {
  value = var.enable_custom_domain ? local.api_custom_fqdn : null
}
