provider "aws" {
	region = "${var.aws_region}"
}

module "iam" {
  source = "../modules/iam"
}

module "api_gateway" {
  source = "../modules/api_gateway"
  api_gateway_aws_account_id = "${var.aws_account_id}"
  api_gateway_aws_region = "${var.aws_region}"
  api_gateway_invoke_lambda_role_arn = "${module.iam.gateway_invoke_lambda_role_arn}"
  api_gateway_stage = "${var.api_stage}"
  api_gateway_api_name = "${var.api_name}"
}

module "dynamodb" {
  source = "../modules/dynamodb"
  dynamodb_user_table = "${var.user_table}"
}

#
output "lambda_function_role_arn" {
  value = "${module.iam.lambda_function_role_arn}"
}

output "gateway_invoke_lambda_role_arn" {
  value = "${module.iam.gateway_invoke_lambda_role_arn}"
}
