variable "aws_account_id" {}

variable "aws_region" {
	default = "us-east-1"
}

variable "stage" {
	default = "dev"
}

variable "user_table" {
	default = "gloauth_users"
}

variable "api_name" {
	default = "Gloauth"
}

variable "api_stage" {
	default = "dev"
}