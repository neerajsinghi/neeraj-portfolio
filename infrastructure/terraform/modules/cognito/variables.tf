variable "project" { type = string }
variable "aws_region" { type = string }
variable "callback_url" { type = string }
variable "logout_url" { type = string }
variable "passkey_satisfies_mfa" {
  type        = bool
  default     = true
  description = "Let a user-verified passkey satisfy required MFA. Runs the AWS CLI locally, so it needs a recent AWS CLI v2 where terraform apply runs."
}
