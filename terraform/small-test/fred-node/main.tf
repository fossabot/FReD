resource "aws_instance" "fred_instance" {
  ami             = data.aws_ami.amazonlinux2.id
  instance_type   = var.instance_type
  key_name        = aws_key_pair.my-test-key.key_name

  security_groups = var.security_groups

  provisioner "file" {
    source      = "./fred-node/config.toml"
    destination = "/tmp/config.toml"
  }


  provisioner "file" {
    source      = "./fred-node/setup_node.sh"
    destination = "/tmp/script.sh"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /tmp/script.sh",
      "/tmp/script.sh ${var.gitlab_repo_username} ${var.gitlab_repo_password}",
    ]
  }

  connection {
    type          = "ssh"
    user          = "ec2-user"
    private_key   = file("terraform.key")
    host          = self.public_ip
  }

  tags = {
    Name = "test-instance"
    type = "fred"
  }
}