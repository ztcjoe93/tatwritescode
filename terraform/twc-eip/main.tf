
resource "aws_eip" "lb" {
  vpc = true
  tags = {
    "Name" : "Main EIP"
  }
}
