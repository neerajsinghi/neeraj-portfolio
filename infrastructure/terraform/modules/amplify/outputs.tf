output "app_id"  { value = aws_amplify_app.frontend.id }
output "app_url" { value = "https://main.${aws_amplify_app.frontend.id}.amplifyapp.com" }
