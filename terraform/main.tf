data "aws_caller_identity" "current" {}

resource "aws_ecr_repository" "mod" {
  name = "${var.ecr_repo_name}"
}

resource "aws_dynamodb_table" "mod" {
  name           = "${var.dynamodb_table_name}"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Token"

  attribute {
    name = "Token"
    type = "S"
  }
}

resource "aws_iam_user" "mod" {
  name = "${var.iam_user_name}"
}

resource "aws_iam_access_key" "mod" {
  user = "${aws_iam_user.mod.name}"
}

resource "aws_iam_user_policy" "mod" {
  name   = "${aws_iam_user.mod.name}"
  user   = "${aws_iam_user.mod.name}"
  policy = "${data.template_file.user_policy.rendered}"
}

data "template_file" "user_policy" {
  template = "${file("${path.module}/user_policy.json")}"

  vars {
    region     = "${var.region}"
    account_id = "${data.aws_caller_identity.current.account_id}"
    repository = "${aws_ecr_repository.mod.name}"
  }
}
