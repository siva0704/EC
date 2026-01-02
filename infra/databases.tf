# Cloud SQL - PostgreSQL (Consolidated: User, Product, Order)
# Standardizing on Postgres reduces operational complexity (single backup/tuning strategy)
resource "google_sql_database_instance" "postgres_instance" {
  name             = "grocery-postgres-main"
  database_version = "POSTGRES_15"
  region           = var.region
  project          = var.project_id

  settings {
    tier = "db-custom-4-15360" # Increased resources for shared load
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.vpc.id
    }
    backup_configuration {
      enabled                        = true
      point_in_time_recovery_enabled = true
    }
    
    # Enable Logical Decoding for CDC (Datastream)
    database_flags {
      name  = "cloudsql.logical_decoding"
      value = "on"
    }
  }
  deletion_protection = false # For demo purposes only
  depends_on          = [google_service_networking_connection.private_vpc_connection]
}

resource "google_sql_database" "user_db" {
  name     = "user_db"
  instance = google_sql_database_instance.postgres_instance.name
  project  = var.project_id
}

resource "google_sql_database" "product_db" {
  name     = "product_db" # CDC Source
  instance = google_sql_database_instance.postgres_instance.name
  project  = var.project_id
}

# Cloud Spanner - Order Service (Horizontal Scale)
resource "google_spanner_instance" "order_spanner_instance" {
  name         = "grocery-orders-spanner"
  config       = "regional-${var.region}"
  display_name = "Grocery Orders Spanner"
  num_nodes    = 1 # Start with 1, scale as needed for 10M MAU
  project      = var.project_id
}

resource "google_spanner_database" "order_spanner_db" {
  instance = google_spanner_instance.order_spanner_instance.name
  name     = "orders"
  project  = var.project_id
  deletion_protection = false

  ddl = [
    "CREATE TABLE Orders (OrderId STRING(36) NOT NULL, UserId STRING(36) NOT NULL, TotalAmount FLOAT64, Status STRING(20), CreatedAt TIMESTAMP) PRIMARY KEY (OrderId)",
    "CREATE INDEX OrdersByUserId ON Orders(UserId)"
  ]
}


# Memorystore for Redis (Cart, Distributed Locking)
resource "google_redis_instance" "redis_main" {
  name           = "grocery-redis-main"
  memory_size_gb = 5
  region         = var.region
  project        = var.project_id
  location_id    = "${var.region}-a"

  # Standard Tier implies HA with a replica
  tier = "STANDARD_HA"

  authorized_network = google_compute_network.vpc.id
  redis_version     = "REDIS_7_0"

  depends_on = [google_service_networking_connection.private_vpc_connection]
}

output "postgres_instance_connection_name" {
  value = google_sql_database_instance.postgres_instance.connection_name
}

output "mysql_instance_connection_name" {
  value = google_sql_database_instance.mysql_instance.connection_name
}

output "redis_host" {
  value = google_redis_instance.redis_main.host
}
