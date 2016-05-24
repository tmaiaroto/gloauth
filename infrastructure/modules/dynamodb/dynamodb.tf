resource "aws_dynamodb_table" "GloauthUserTable" {
    name = "${var.dynamodb_user_table}"
    read_capacity = 2
    write_capacity = 2
    hash_key = "id"
    attribute {
      name = "id"
      type = "N"
    }
    attribute {
      name = "email"
      type = "S"
    }
    attribute {
    	name = "last_login"
    	type = "N"
    }
    global_secondary_index {
      name = "EmailIndex"
      hash_key = "email"
      range_key = "last_login"
      write_capacity = 2
      read_capacity = 2
      projection_type = "ALL"
    }
}