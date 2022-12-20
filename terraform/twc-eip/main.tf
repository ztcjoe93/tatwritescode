
resource "aws_eip" "lb" {
  vpc = true
  tags = {
    "Name" : "twc-main-eip"
  }
}
