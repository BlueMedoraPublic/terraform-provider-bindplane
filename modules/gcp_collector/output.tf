output "project" {
  value = google_compute_instance.collector.project
}

output "collector_name_prefix" {
  value = var.collector_name_prefix
}

output "compute_instance_name" {
  value = google_compute_instance.collector.name
}

output "network_zone" {
  value = google_compute_instance.collector.zone
}

output "compute_instance_size" {
  value = google_compute_instance.collector.machine_type
}

output "compute_image" {
  value = var.compute_image
}

output "network" {
  value = var.network
}

output "disk_type" {
  value = var.disk_type
}

output "disk_size" {
  value = var.disk_size
}

output "id" {
  value = bindplane_collector.collector.id
}
