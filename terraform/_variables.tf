locals {
  region        = "ap-southeast-1"
  database_name = "tatwritescode"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t2.micro"
}

variable "database_user" {
  type        = string
  description = "Name of database"
}

variable "database_password" {
  type        = string
  description = "Password for database_user"
}

variable "database_root_password" {
  type        = string
  description = "Password for root database user"
}
