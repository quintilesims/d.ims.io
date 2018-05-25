variable "dynamodb_table_name" {
  default = "d.ims.io"
}

variable "dynamodb_account_table_name" {
  default = "d.ims.io.account"
}

variable "dynamodb_read_capacity" {
  default = 1
}

variable "dynamodb_write_capacity" {
  default = 1
}

variable "debug" {
  default = false
}

variable "s3_bucket_name" {
  default = "d.ims.io-backups"
}

variable "iam_user_name" {
  default = "d.ims.io"
}

variable "aws_region" {
  default = "us-west-2"
}

variable "environment_id" {
  description = "ID of the Layer0 environment to build the service"
}

variable "load_balancer_name" {
  description = "Name of the Layer0 load balancer to create"
  default     = "dimsio"
}

variable "service_name" {
  description = "Name of the Layer0 service to create"
  default     = "dimsio"
}

variable "deploy_name" {
  description = "Name of the Layer0 deploy to create"
  default     = "dimsio"
}

variable "scale" {
  description = "The scale of the service"
  default     = 1
}

variable "certificate_arn" {
  description = "The AWS ARN of the ACM certificate to use on the load balancer"
}

variable "docker_image" {
  description = "Docker image to use"
  default     = "quintilesims/d.ims.io:latest"
}

variable "auth0_domain" {
  default = "https://imshealth.auth0.com"
}

variable "auth0_connection" {
  description = "The name of the Auth0 connnection to use for AD validation"
}

variable "auth0_client_id" {
  description = "The client id associated with the Auth0 connection"
}

variable "lambda_function_name" {
  default = "dimsio_backups"
}
