output "load_balancer_url" {
  value = "${layer0_load_balancer.dimsio.url}"
}

output "load_balancer_id" {
  value = "${layer0_load_balancer.dimsio.id}"
}

output "service_id" {
  value = "${layer0_service.dimsio.id}"
}

output "deploy_id" {
  value = "${layer0_deploy.dimsio.id}"
}
