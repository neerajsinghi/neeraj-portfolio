variable "project" { type = string }
variable "github_repo" { type = string }
variable "api_base_url" { type = string }

variable "environment_variables" {
  type    = map(string)
  default = {}
}

variable "branch_name" {
  type    = string
  default = "main"
}

variable "frontend_app_root" {
  type    = string
  default = "frontend"
}

variable "enable_auto_build" {
  type    = bool
  default = false
}

variable "enable_custom_domain" {
  type    = bool
  default = false
}

variable "root_domain" {
  type    = string
  default = ""
}

variable "frontend_subdomain" {
  type    = string
  default = ""
}

variable "frontend_enable_www" {
  type    = bool
  default = true
}

variable "github_access_token" {
  type      = string
  sensitive = true
}
