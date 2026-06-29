variable "project"     { type = string }
variable "github_repo" { type = string }
variable "api_base_url" { type = string }

variable "github_access_token" {
  type      = string
  sensitive = true
}
