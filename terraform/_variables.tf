locals {
  region = "ap-southeast-1"
}

variable "instance_type" {
  type        = string
  description = "EC2 instance type"
  default     = "t2.micro"
}
