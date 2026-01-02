package main

import (
	"context"
	"encoding/json"
	"log"
	"os"

	"cloud.google.com/go/pubsub"
	"database/sql"
	"github.com/elastic/go-elasticsearch/v8"
	_ "github.com/lib/pq"
)

// Reconciler checks for drift between Postgres (Source) and Elastic (Sink)
func main() {
	ctx := context.Background()
	log.Println("Starting Reconciliation Job...")

	// 1. Connect to Dependencies
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}
	defer db.Close()

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating ES client: %v", err)
	}

	pubsubClient, err := pubsub.NewClient(ctx, os.Getenv("PROJECT_ID"))
	if err != nil {
		log.Fatalf("Error creating PubSub client: %v", err)
	}
	topic := pubsubClient.Topic(os.Getenv("CDC_TOPIC"))

	// 2. Iterate Source (Postgres)
	// In prod, use cursor/pagination. Simplified here.
	rows, err := db.Query("SELECT id, updated_at FROM products")
	if err != nil {
		log.Fatalf("DB Query failed: %v", err)
	}
	defer rows.Close()

	driftCount := 0

	for rows.Next() {
		var id string
		var updatedAt string
		if err := rows.Scan(&id, &updatedAt); err != nil {
			log.Println("Row scan error:", err)
			continue
		}

		// 3. Check Sink (Elasticsearch)
		// Optimization: Use Multi-Get (MGET) in batches for performance
		exists := checkElasticExistence(es, id) // Stubbed

		if !exists {
			log.Printf("DRIFT DETECTED: Product %s missing in Search Index", id)
			
			// 4. Repair: Trigger CDC Event
			msg := map[string]string{
				"type": "REPAIR",
				"table": "products",
				"id": id,
			}
			data, _ := json.Marshal(msg)
			topic.Publish(ctx, &pubsub.Message{Data: data})
			driftCount++
		}
	}

	log.Printf("Reconciliation Complete. Drifts repaired: %d", driftCount)
}

func checkElasticExistence(es *elasticsearch.Client, id string) bool {
    // Implement HEAD /products/_doc/ID
    return true // Stub: Assume synced for now
}
