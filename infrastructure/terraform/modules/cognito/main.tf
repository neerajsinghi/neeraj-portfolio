data "aws_caller_identity" "current" {}

locals {
  domain_prefix = "${var.project}-${data.aws_caller_identity.current.account_id}"
  # Passkeys are bound to the domain that hosts managed login.
  relying_party_id = "${local.domain_prefix}.auth.${var.aws_region}.amazoncognito.com"
}

resource "aws_cognito_user_pool" "admin" {
  name                     = "${var.project}-admin"
  username_attributes      = ["email"]
  auto_verified_attributes = ["email"]
  mfa_configuration        = "ON"
  # Passkeys are unavailable on the Lite tier.
  user_pool_tier = "ESSENTIALS"

  software_token_mfa_configuration {
    enabled = true
  }

  sign_in_policy {
    allowed_first_auth_factors = ["PASSWORD", "WEB_AUTHN"]
  }

  web_authn_configuration {
    relying_party_id  = local.relying_party_id
    user_verification = "required"
  }

  password_policy {
    minimum_length                   = 14
    require_lowercase                = true
    require_numbers                  = true
    require_symbols                  = true
    require_uppercase                = true
    temporary_password_validity_days = 2
  }

  account_recovery_setting {
    recovery_mechanism {
      name     = "verified_email"
      priority = 1
    }
  }

  admin_create_user_config {
    allow_admin_create_user_only = true
  }
}

resource "aws_cognito_user_pool_client" "admin" {
  name         = "${var.project}-admin-web"
  user_pool_id = aws_cognito_user_pool.admin.id

  generate_secret                      = false
  allowed_oauth_flows_user_pool_client = true
  allowed_oauth_flows                  = ["code"]
  allowed_oauth_scopes                 = ["openid", "email", "profile"]
  explicit_auth_flows                  = ["ALLOW_USER_AUTH", "ALLOW_USER_SRP_AUTH", "ALLOW_REFRESH_TOKEN_AUTH"]
  supported_identity_providers         = ["COGNITO"]
  callback_urls                        = [var.callback_url]
  logout_urls                          = [var.logout_url]
  prevent_user_existence_errors        = "ENABLED"
  enable_token_revocation              = true
  access_token_validity                = 1
  id_token_validity                    = 1
  refresh_token_validity               = 7

  token_validity_units {
    access_token  = "hours"
    id_token      = "hours"
    refresh_token = "days"
  }
}

resource "aws_cognito_user_pool_domain" "admin" {
  domain       = local.domain_prefix
  user_pool_id = aws_cognito_user_pool.admin.id
  # Passkey sign-in and enrollment exist only in managed login, not the classic hosted UI.
  managed_login_version = 2
}

resource "aws_cognito_managed_login_branding" "admin" {
  user_pool_id                = aws_cognito_user_pool.admin.id
  client_id                   = aws_cognito_user_pool_client.admin.id
  use_cognito_provided_values = true
}

# The AWS provider does not expose WebAuthnConfiguration.FactorConfiguration yet.
# Without it a passkey is only a first factor and TOTP is still prompted; with it,
# a passkey that verifies the user (biometric/PIN) satisfies the required MFA.
# The provider rewrites the MFA config whenever these inputs change, so rerun then.
resource "terraform_data" "passkey_mfa" {
  count = var.passkey_satisfies_mfa ? 1 : 0

  triggers_replace = [
    aws_cognito_user_pool.admin.id,
    aws_cognito_user_pool.admin.mfa_configuration,
    aws_cognito_user_pool.admin.software_token_mfa_configuration,
    aws_cognito_user_pool.admin.web_authn_configuration,
    aws_cognito_user_pool.admin.user_pool_tier,
  ]

  provisioner "local-exec" {
    command = <<-EOT
      aws cognito-idp set-user-pool-mfa-config \
        --region ${var.aws_region} \
        --user-pool-id ${aws_cognito_user_pool.admin.id} \
        --mfa-configuration ON \
        --software-token-mfa-configuration Enabled=true \
        --web-authn-configuration RelyingPartyId=${local.relying_party_id},UserVerification=required,FactorConfiguration=MULTI_FACTOR_WITH_USER_VERIFICATION
    EOT
  }
}

resource "aws_cognito_user_group" "admin" {
  name         = "admin"
  user_pool_id = aws_cognito_user_pool.admin.id
  description  = "Can publish, archive, and delete blog posts"
  precedence   = 1
}

resource "aws_cognito_user_group" "editor" {
  name         = "editor"
  user_pool_id = aws_cognito_user_pool.admin.id
  description  = "Can create and edit draft blog posts"
  precedence   = 2
}