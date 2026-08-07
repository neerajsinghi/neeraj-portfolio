variable "project"     { type = string }
variable "github_repo" { type = string }
variable "api_base_url" { type = string }

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

variable "github_access_token" {
  type      = string
  sensitive = true
}
