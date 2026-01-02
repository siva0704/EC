resource "google_container_cluster" "primary" {
  name     = "grocery-platform-cluster"
  location = var.region
  project  = var.project_id

  # Autopilot for reduced operational overhead
  enable_autopilot = true

  network    = google_compute_network.vpc.id
  subnetwork = google_compute_subnetwork.gke_subnet.id

  # Private Cluster Config
  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false # Allowed public endpoint for ease of access in this demo
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }

  ip_allocation_policy {
    cluster_secondary_range_name  = google_compute_subnetwork.gke_subnet.secondary_ip_range[0].range_name
    services_secondary_range_name = google_compute_subnetwork.gke_subnet.secondary_ip_range[1].range_name
  }

  # Workload Identity is enabled by default on Autopilot
  # But explicit config helps document intent
  workload_identity_config {
    workload_pool = "${var.project_id}.svc.id.goog"
  }

  # Enable Managed Service Mesh (ASM) for mTLS and Observability
  # fleet {
  #   project = var.project_id
  # }
  
  # release_channel {
  #   channel = "REGULAR"
  # }
}

# ASM Feature Registration (Requires Google Beta provider)
# resource "google_gke_hub_membership" "membership" {
#   provider      = google-beta
#   membership_id = "basic"
#   endpoint {
#     gke_cluster {
#       resource_link = "//container.googleapis.com/${google_container_cluster.primary.id}"
#     }
#   }
# }

# Output for credentials
output "kubernetes_cluster_name" {
  value = google_container_cluster.primary.name
}

output "kubernetes_cluster_host" {
  value = google_container_cluster.primary.endpoint
}
