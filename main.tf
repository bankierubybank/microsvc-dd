terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.84"
    }
  }

  required_version = ">= 1.2.0"
}

provider "aws" {
  region = "ap-southeast-7"
}

resource "aws_instance" "app_server" {
  ami           = "ami-01b7b0559e63daaf3"
  instance_type = "t3.micro"

  tags = {
    Name = "tarot"
  }
}
