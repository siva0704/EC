# Topic: Order Events (Order -> Payment, Notification)
resource "google_pubsub_topic" "order_events" {
  name    = "order-events"
  project = var.project_id
}

resource "google_pubsub_subscription" "payment_processor_sub" {
  name    = "payment-processor-sub"
  topic   = google_pubsub_topic.order_events.name
  project = var.project_id

  # Dead Letter Policy for robustness
  dead_letter_policy {
    dead_letter_topic     = google_pubsub_topic.order_events_dlq.id
    max_delivery_attempts = 5
  }
}

resource "google_pubsub_topic" "order_events_dlq" {
  name    = "order-events-dlq"
  project = var.project_id
}


# Topic: Product CDC (Product DB -> Search Index)
resource "google_pubsub_topic" "product_cdc" {
  name    = "product-cdc"
  project = var.project_id
}

resource "google_pubsub_subscription" "search_indexer_sub" {
  name    = "search-indexer-sub"
  topic   = google_pubsub_topic.product_cdc.name
  project = var.project_id
}


# Topic: Payment Events (Payment -> Order)
resource "google_pubsub_topic" "payment_events" {
  name    = "payment-events"
  project = var.project_id
}
