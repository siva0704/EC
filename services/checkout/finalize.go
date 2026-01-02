package checkout

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/google/uuid"
	"google.golang.org/api/iterator"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
)

// OrderItem represents a line item in the order
type OrderItem struct {
	ItemID   string
	Quantity int
    Price    float64
}

// FinalizeOrder executes the strict consistency phase of the purchase
func (s *CheckoutService) FinalizeOrder(ctx context.Context, userID string, items []OrderItem) error {
    // 1. Transaction: Write Order, Items, and Outbox Event atomically
	_, err := s.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
        
        // A. Validate "Real" Inventory (Double Check) if needed, 
        // but primarily we trust the Redis reservation from the previous step.
        // However, if we wanted to support "Partial Stock" explicitly here:
        // We would read current stock from Spanner (if we were syncing it back) or trust Redis.
        // For this implementation, we assume Redis reservation was successful.
        
        orderID := uuid.New().String()
        totalAmount := calculateTotal(items)

        // B. Insert Order
        orderMutation := spanner.Insert("Orders",
            []string{"OrderId", "UserId", "TotalAmount", "Status", "CreatedAt"},
            []interface{}{orderID, userID, totalAmount, "CREATED", spanner.CommitTimestamp},
        )

        // C. Insert Order Items
        var mutations []*spanner.Mutation
        mutations = append(mutations, orderMutation)

        for _, item := range items {
            itemMutation := spanner.Insert("OrderItems",
                []string{"OrderId", "ItemId", "Quantity", "Price"},
                []interface{}{orderID, item.ItemID, item.Quantity, item.Price},
            )
            mutations = append(mutations, itemMutation)
        }

        // D. Outbox Pattern: Insert Event for Asynchronous Processing (Payment, Email)
        // This ensures the event is ONLY produced if the Order is committed.
        outboxID := uuid.New().String()
        outboxMutation := spanner.Insert("Outbox",
            []string{"EventId", "AggregateId", "Type", "Payload", "CreatedAt"},
            []interface{}{outboxID, orderID, "ORDER_CREATED", fmt.Sprintf(`{"user_id": "%s", "total": %f}`, userID, totalAmount), spanner.CommitTimestamp},
        )
        mutations = append(mutations, outboxMutation)

		return txn.BufferWrite(mutations)
	})

	if err != nil {
		return fmt.Errorf("spanner transaction failed: %v", err)
	}

    // Note: A separate "Outbox Relay" process will tail the Outbox table 
    // and publish to Pub/Sub. We do NOT publish here to avoid dual-write issues.
	return nil
}

// CheckStockAndGetModifiedCart handles the "Partial Stock" scenario check before generating the order
func (s *CheckoutService) CheckStockAndGetModifiedCart(ctx context.Context, items []OrderItem) ([]OrderItem, bool, error) {
    var modifiedItems []OrderItem
    hasChanges := false

    for _, item := range items {
        // Check Redis for current availability
        stockKey := fmt.Sprintf("stock:%s", item.ItemID)
        val, err := s.redisClient.Get(ctx, stockKey).Int()
        if err != nil {
            return nil, false, fmt.Errorf("failed to check stock for %s: %v", item.ItemID, err)
        }

        if val < item.Quantity {
            // Partial Stock Scenario
            hasChanges = true
            msg := fmt.Sprintf("Requested %d, only %d available", item.Quantity, val)
            // Logic to ask user: return modified quantity (e.g. max available)
            if val > 0 {
                modifiedItems = append(modifiedItems, OrderItem{
                    ItemID: item.ItemID,
                    Quantity: val, // Take what is left
                    Price: item.Price,
                })
            }
            // If 0, item is removed from modified cart
            fmt.Println(msg) // Logging constraint
        } else {
            modifiedItems = append(modifiedItems, item)
        }
    }

    if hasChanges {
        // In a real API, we would return a 409 Conflict with this payload
        return modifiedItems, true, nil 
    }

    return items, false, nil
}

func calculateTotal(items []OrderItem) float64 {
    total := 0.0
    for _, item := range items {
        total += item.Price * float64(item.Quantity)
    }
    return total
}
