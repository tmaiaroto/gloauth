output "lambda_function_role_id" {
  value = "${aws_iam_role.lambda_function.arn}"
}

output "gateway_invoke_lambda_role_id" {
  value = "${aws_iam_role.gateway_invoke_lambda.arn}"
}