data "aws_caller_identity" "current" {}

resource "aws_dynamodb_table" "dimsio" {
  name           = "${var.dynamodb_table_name}"
  read_capacity  = 1
  write_capacity = 1
  hash_key       = "Token"

  attribute {
    name = "Token"
    type = "S"
  }
}

resource "aws_iam_user" "dimsio" {
  name = "${var.iam_user_name}"
}

resource "aws_iam_access_key" "dimsio" {
  user = "${aws_iam_user.dimsio.name}"
}

resource "aws_iam_user_policy" "dimsio" {
  name   = "${aws_iam_user.dimsio.name}"
  user   = "${aws_iam_user.dimsio.name}"
  policy = "${data.template_file.user_policy.rendered}"
}

data "template_file" "user_policy" {
  template = "${file("${path.module}/user_policy.json")}"

  vars {
    dynamodb_table_arn = "${aws_dynamodb_table.dimsio.arn}"
  }
}

resource "layer0_load_balancer" "dimsio" {
  name        = "${var.load_balancer_name}"
  environment = "${var.environment_id}"

  port {
    host_port      = "443"
    container_port = "80"
    protocol       = "https"
    certificate    = "${var.certificate_name}"
  }

  health_check {
    target = "tcp:80"
  }
}

resource "layer0_service" "dimsio" {
  name          = "${var.service_name}"
  environment   = "${var.environment_id}"
  deploy        = "${layer0_deploy.dimsio.id}"
  load_balancer = "${layer0_load_balancer.dimsio.id}"
  scale         = "${var.scale}"
  wait          = true
}

resource "layer0_deploy" "dimsio" {
  name    = "${var.deploy_name}"
  content = "${data.template_file.dimsio.rendered}"
}

data "template_file" "dimsio" {
  template = "${file("${path.module}/Dockerrun.aws.json")}"

  vars {
    docker_image      = "${var.docker_image}"
    aws_access_key    = "${aws_iam_access_key.dimsio.id}"
    aws_secret_key    = "${aws_iam_access_key.dimsio.secret}"
    aws_region        = "${var.aws_region}"
    swagger_host      = "${layer0_load_balancer.dimsio.url}"
    dynamo_table      = "${aws_dynamodb_table.dimsio.name}"
    registry_endpoint = "${data.aws_caller_identity.current.account_id}.dkr.ecr.${var.aws_region}.amazonaws.com"
    auth0_domain      = "${var.auth0_domain}"
    auth0_client_id   = "${var.auth0_client_id}"
    auth0_connection  = "${var.auth0_connection}"
  }
}
