output "app_id" {
	value = aws_amplify_app.frontend.id
}

output "app_default_domain" {
	value = aws_amplify_app.frontend.default_domain
}

output "branch_name" {
	value = aws_amplify_branch.main.branch_name
}

# aws_amplify_branch has no "web_url" attribute; Amplify branch URLs follow this fixed pattern.
output "branch_web_url" {
	value = "https://${aws_amplify_branch.main.branch_name}.${aws_amplify_app.frontend.default_domain}"
}

output "app_url" {
	value = "https://${aws_amplify_branch.main.branch_name}.${aws_amplify_app.frontend.default_domain}"
}

output "custom_domain" {
	value = var.enable_custom_domain ? aws_amplify_domain_association.frontend[0].domain_name : null
}

output "frontend_fqdn" {
	value = var.enable_custom_domain ? (
		var.frontend_subdomain == "" ? var.root_domain : "${var.frontend_subdomain}.${var.root_domain}"
	) : null
}
