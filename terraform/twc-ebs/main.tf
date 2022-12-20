resource "aws_ebs_volume" "ebs_volume" {
  availability_zone = local.availability_zone
  type              = "gp3"
  size              = 30

  tags = {
    "Name" : "twc-main-ebs-1a"
  }
}
