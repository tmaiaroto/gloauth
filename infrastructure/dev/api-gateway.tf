# Creates the API "Gloauth API"
resource "aws_api_gateway_rest_api" "GloauthAPI" {
  name = "Gloauth API"
  description = "A simple serverless authentication system"
}

# Creates the /register API resource path
resource "aws_api_gateway_resource" "RegisterResource" {
  rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
  parent_id = "${aws_api_gateway_rest_api.GloauthAPI.root_resource_id}"
  path_part = "register"
}

# Creates the GET method under /register
resource "aws_api_gateway_method" "RegisterMethod" {
  rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
  resource_id = "${aws_api_gateway_resource.RegisterResource.id}"
  http_method = "GET"
  authorization = "NONE"
}

# Creates a method response; 200, 404, 500 etc.
resource "aws_api_gateway_method_response" "200" {
  rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
  resource_id = "${aws_api_gateway_resource.RegisterResource.id}"
  http_method = "${aws_api_gateway_method.RegisterMethod.http_method}"
  status_code = "200"
  response_models = {
    "application/json" = "Empty"
  }
}

# Configures the integration for the Resource Method (what gets triggered for GET /register)
resource "aws_api_gateway_integration" "RegisterIntegration" {
  rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
  resource_id = "${aws_api_gateway_resource.RegisterResource.id}"
  http_method = "${aws_api_gateway_method.RegisterMethod.http_method}"
  type = "AWS"
  integration_http_method = "POST" # Must be POST for invoking Lambda function
  credentials = "${module.iam.gateway_invoke_lambda_role_id}"
  # http://docs.aws.amazon.com/apigateway/api-reference/resource/integration/#uri
  uri = "arn:aws:apigateway:${var.aws_region}:lambda:path/2015-03-31/functions/arn:aws:lambda:${var.aws_region}:${var.aws_account_id}:function:${var.functions.register}/invocations"
}

# Configures the integration (Lambda) response that maps to the method response via the status_code
resource "aws_api_gateway_integration_response" "RegisterResponse" {
  rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
  resource_id = "${aws_api_gateway_resource.RegisterResource.id}"
  http_method = "${aws_api_gateway_method.RegisterMethod.http_method}"
  status_code = "${aws_api_gateway_method_response.200.status_code}"
}

# Creates the API stage
resource "aws_api_gateway_deployment" "dev" {
 depends_on = ["aws_api_gateway_integration.RegisterIntegration"]

 rest_api_id = "${aws_api_gateway_rest_api.GloauthAPI.id}"
 stage_name = "dev"
}