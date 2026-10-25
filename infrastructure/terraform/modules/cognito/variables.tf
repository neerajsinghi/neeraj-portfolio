variable "project" { type = string }
variable "aws_region" { type = string }
variable "callback_url" { type = string }
variable "logout_url" { type = string }
variable "passkey_satisfies_mfa" {
  type        = bool
  default     = true
  description = "Allow passkey sign-in, with a user-verified passkey satisfying required MFA. false disables passkey sign-in (Cognito allows passkeys as a first factor under required MFA only when they count as MFA). MFA is set with the AWS CLI, so terraform apply needs a recent AWS CLI v2."
}
