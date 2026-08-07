output "app_id" {
	value = aws_amplify_app.frontend.id
}

output "app_default_domain" {
	value = aws_amplify_app.frontend.default_domain
}

output "branch_name" {
	value = aws_amplify_branch.main.branch_name
}

output "branch_web_url" {
	value = aws_amplify_branch.main.web_url
}

output "app_url" {
	value = aws_amplify_branch.main.web_url
}
