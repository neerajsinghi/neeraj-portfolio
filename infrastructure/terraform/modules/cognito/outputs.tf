output "user_pool_id" { value = aws_cognito_user_pool.admin.id }
output "client_id" { value = aws_cognito_user_pool_client.admin.id }
output "domain" { value = "https://${aws_cognito_user_pool_domain.admin.domain}.auth.${var.aws_region}.amazoncognito.com" }