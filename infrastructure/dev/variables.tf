variable "aws_account_id" {
}

variable "aws_region" {
	default = "us-east-1"
}

variable "stage" {
	default = "dev"
}

variable "functions" {
	type = "map"
	default = {
		register = "gloauth_register"
	}
}