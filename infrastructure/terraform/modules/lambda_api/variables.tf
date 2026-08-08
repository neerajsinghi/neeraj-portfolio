variable "project" {
  type = string
}

variable "aws_region" {
  type = string
}

variable "lambda_zip_path" {
  type = string
}

variable "lambda_memory_size" {
  type    = number
  default = 512
}

variable "lambda_timeout" {
  type    = number
  default = 60
}

variable "allowed_origin" {
  type    = string
  default = "*"
}

variable "cognito_user_pool_id" { type = string }
variable "cognito_client_id" { type = string }

variable "enable_custom_domain" {
  type    = bool
  default = false
}

variable "root_domain" {
  type    = string
  default = ""
}

variable "api_subdomain" {
  type    = string
  default = "api"
}
