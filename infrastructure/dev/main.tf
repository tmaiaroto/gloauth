provider "aws" {
	region = "${var.aws_region}"
}

module "iam" {
  source = "../modules/iam"
}

module "dynamodb" {
  source = "../modules/dynamodb"
  dynamodb_user_table = "${var.user_table}"
}

output "lambda_function_role_id" {
  value = "${module.iam.lambda_function_role_id}"
}
