output "source_name" {
  value = bindplane_log_source.eks.name
}

output "source_type_id" {
  value = bindplane_log_source.eks.source_type_id
}

output "source_type_version" {
  value = bindplane_log_source.eks.source_version
}

output "container_logs" {
  value = var.container_logs
}

output "container_file_exp" {
  value = var.container_file_exp
}

output "kube_proxy_logs" {
  value = var.kube_proxy_logs
}

output "kubelet_logs" {
  value = var.kubelet_logs
}

output "controller_manager_logs" {
  value = var.controller_manager_logs
}

output "scheduler_logs" {
  value = var.scheduler_logs
}

output "apiserver_logs" {
  value = var.apiserver_logs
}

output "template_name" {
  value = bindplane_log_template.eks.name
}

output "template_source_config_ids" {
  value = bindplane_log_template.eks.source_config_ids
}

output "template_destination_config_id" {
  value = bindplane_log_template.eks.destination_config_id
}

output "template_id" {
  value = bindplane_log_template.eks.id
}
