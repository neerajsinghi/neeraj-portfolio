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
