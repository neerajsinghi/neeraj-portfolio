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

output "custom_domain" {
	value = var.enable_custom_domain ? aws_amplify_domain_association.frontend[0].domain_name : null
}

output "frontend_fqdn" {
	value = var.enable_custom_domain ? (
		var.frontend_subdomain == "" ? var.root_domain : "${var.frontend_subdomain}.${var.root_domain}"
	) : null
}
