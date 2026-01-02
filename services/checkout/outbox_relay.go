package checkout

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/spanner"
	"google.golang.org/api/iterator"
)

type OutboxRelay struct {
	spannerClient *spanner.Client
	pubsubClient  *pubsub.Client
    topicID       string
}

type OutboxEvent struct {
    EventID     string
    AggregateID string
    Type        string
    Payload     string
}

// StartPoller begins the loop to process outbox events
func (r *OutboxRelay) StartPoller(ctx context.Context) {
    ticker := time.NewTicker(200 * time.Millisecond) // Poll frequency
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            r.processBatch(ctx)
        }
    }
}

func (r *OutboxRelay) processBatch(ctx context.Context) {
    // 1. Transaction: Lock methods or just Delete-Returning pattern
    // Spanner doesn't support DELETE RETURNING in the same way Postgres does for queuing.
    // We will use a ReadWriteTransaction to Read, Publish, Delete.
    
    _, err := r.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
        // A. Read pending events (Limit 10 to avoid congestion)
        // Ensure index on CreatedAt ASC
        stmt := spanner.Statement{SQL: "SELECT EventId, AggregateId, Type, Payload FROM Outbox ORDER BY CreatedAt ASC LIMIT 10"}
        iter := txn.Query(ctx, stmt)
        
        var events []OutboxEvent
        err := iter.Do(func(row *spanner.Row) error {
            var e OutboxEvent
            if err := row.ToStruct(&e); err != nil {
                return err
            }
            events = append(events, e)
            return nil
        })
        if err != nil {
            return err
        }

        if len(events) == 0 {
            return nil
        }

        // B. Publish to Pub/Sub (Inside transaction? No, usually outside or robust retry inside)
        // If we fail after publish, we might re-publish (At-Least-Once). 
        // Better to publish here. If publish fails, we abort transaction and retry later.
        
        var deleteKeys []spanner.KeySet
        topic := r.pubsubClient.Topic(r.topicID)
        
        for _, e := range events {
            msgData, _ := json.Marshal(e)
            res := topic.Publish(ctx, &pubsub.Message{
                Data: msgData,
                Attributes: map[string]string{
                    "type": e.Type,
                    "aggregate_id": e.AggregateID,
                },
            })
            
            // Block for result to ensure Pub/Sub accepted it before deleting from DB
            // This hurts latency but ensures consistency.
            if _, err := res.Get(ctx); err != nil {
                return fmt.Errorf("failed to publish event %s: %v", e.EventID, err)
            }
            
            deleteKeys = append(deleteKeys, spanner.Key{e.EventID})
        }

        // C. Delete processed events from Outbox
        return txn.BufferWrite([]*spanner.Mutation{
            spanner.Delete("Outbox", spanner.KeySets(deleteKeys...)),
        })
    })

    if err != nil {
        log.Printf("Outbox Relay Error: %v", err)
    }
}
