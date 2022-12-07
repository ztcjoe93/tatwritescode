terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 4.0"
    }
  }
}

provider "aws" {
  region = "ap-southeast-1"
}

resource "aws_instance" "ec2_instance" {

  ami           = "ami-0af2f764c580cc1f9"
  instance_type = "t2.micro"
  key_name      = aws_key_pair.ssh_key.key_name

  tags = {
    name = "main-server"
  }

  user_data              = templatefile("${path.module}/init.sh", {})
  vpc_security_group_ids = [aws_security_group.main_sg.id]
}

resource "aws_eip" "lb" {
  instance = aws_instance.ec2_instance.id
  vpc      = true
}

resource "aws_security_group" "main_sg" {
  name        = "main_sg"
  description = "Security group for the main server"
}

resource "aws_security_group_rule" "http_access" {
  type              = "ingress"
  from_port         = 80
  to_port           = 80
  cidr_blocks       = ["0.0.0.0/0"]
  protocol          = "tcp"
  security_group_id = aws_security_group.main_sg.id
}

resource "aws_security_group_rule" "https_access" {
  type              = "ingress"
  from_port         = 443
  to_port           = 443
  cidr_blocks       = ["0.0.0.0/0"]
  protocol          = "tcp"
  security_group_id = aws_security_group.main_sg.id
}

resource "aws_security_group_rule" "ssh_ingress" {
  type              = "ingress"
  from_port         = 22
  to_port           = 22
  cidr_blocks       = ["0.0.0.0/0"]
  protocol          = "tcp"
  security_group_id = aws_security_group.main_sg.id
}

resource "aws_security_group_rule" "egress" {
  type              = "egress"
  from_port         = 0
  to_port           = 0
  cidr_blocks       = ["0.0.0.0/0"]
  protocol          = "-1"
  security_group_id = aws_security_group.main_sg.id
}

resource "aws_key_pair" "ssh_key" {
  key_name   = "ssh_key"
  public_key = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCsRA5D2q7Wdei48Maw7ePTC1McvQHr53rZI74dmOB6WvNXmuVjpmD6N81r26UZj3sztZbpiLLorBwHqDbHbfAfN1VqLym1BvHuWHGNzF+JJ4bxaGDxThZ7NF1k5Kqisje7mpNH7mjX/CFNs95IGFNYREkmzXq+wC1eKuBF0vYkVtjys3mPeAnL5A4y3dNLmgCbROj82jlTVp9v6QAhHVJPp4Mu7STrE8Gp86OKb7QyYd/ZHv+7lcFkte8q5GIi/BN2aYGTgnhzdNUaK3uyRT4NMT/h6vg14KOuzSkiW1Yxb3BCuR2VCAVrsZHdV+lpu6C/b+H9DKbZc7rqFKi0xP+r"
}
