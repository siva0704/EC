resource "google_compute_network" "vpc" {
  name                    = "grocery-platform-vpc"
  auto_create_subnetworks = false
  project                 = var.project_id
}

resource "google_compute_subnetwork" "gke_subnet" {
  name          = "gke-subnet"
  ip_cidr_range = "10.0.0.0/20"
  region        = var.region
  network       = google_compute_network.vpc.id
  project       = var.project_id

  secondary_ip_range {
    range_name    = "gke-pods"
    ip_cidr_range = "10.48.0.0/14"
  }
  secondary_ip_range {
    range_name    = "gke-services"
    ip_cidr_range = "10.52.0.0/20"
  }

  private_ip_google_access = true
}

# Private Service Access for Cloud SQL / Memorystore
resource "google_compute_global_address" "private_ip_range" {
  name          = "private-ip-range"
  purpose       = "VPC_PEERING"
  address_type  = "INTERNAL"
  prefix_length = 16
  network       = google_compute_network.vpc.id
  project       = var.project_id
}

resource "google_service_networking_connection" "private_vpc_connection" {
  network                 = google_compute_network.vpc.id
  service                 = "servicenetworking.googleapis.com"
  reserved_peering_ranges = [google_compute_global_address.private_ip_range.name]
}

# Firewall Rules - Best Practices
# Deny all ingress is implicit in GCP, but we add allow rules strictly.

# Allow IAP for SSH (if needed for debug nodes)
resource "google_compute_firewall" "allow_iap_ssh" {
  name    = "allow-iap-ssh"
  network = google_compute_network.vpc.name
  project = var.project_id

  allow {
    protocol = "tcp"
    ports    = ["22"]
  }

  source_ranges = ["35.235.240.0/20"] # IAP IP Range
  target_tags   = ["bastion"]
}

# Allow Health Checks
resource "google_compute_firewall" "allow_health_checks" {
  name    = "allow-health-checks"
  network = google_compute_network.vpc.name
  project = var.project_id

  allow {
    protocol = "tcp"
  }

  source_ranges = ["130.211.0.0/22", "35.191.0.0/16"]
  target_tags   = ["gke-node"]
}
