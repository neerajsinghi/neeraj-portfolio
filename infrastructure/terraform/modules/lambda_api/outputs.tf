output "function_name" {
  value = aws_lambda_function.backend.function_name
}

output "function_arn" {
  value = aws_lambda_function.backend.arn
}

output "api_base_url" {
  value = aws_apigatewayv2_stage.default.invoke_url
}
