locals {
  availability_zone = "ap-southeast-1a"
  region            = "ap-southeast-1"
  database_name     = "tatwritescode"
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

variable "volume_mount_path" {
  type        = string
  description = "Path for database container's volume"
}

variable "ssl_pem" {
  type        = string
  description = "Public key for SSL cert"
}

variable "ssl_key" {
  type        = string
  description = "Private key for SSL cert"
}
