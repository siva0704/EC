package checkout

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/spanner"
	"github.com/go-redis/redis/v8"
	"github.com/bsm/redislock"
)

type CheckoutService struct {
	redisClient *redis.Client
	locker      *redislock.Client
	spannerClient *spanner.Client
}

// handlePurchase implements Reserve-then-Commit pattern
func (s *CheckoutService) HandlePurchase(ctx context.Context, userID string, items []AppItem) error {
	// 1. Acquire Distributed Lock (Redlock)
	// Keyed by UserID to prevent double submission, or ItemID for hot items
	lockKey := fmt.Sprintf("lock:checkout:%s", userID)
	lock, err := s.locker.Obtain(ctx, lockKey, 5*time.Second, nil)
	if err != nil {
		return fmt.Errorf("could not acquire lock: %v", err)
	}
	defer lock.Release(ctx)

	// 2. Reserve Inventory (Virtual Stock in Redis)
	// Lua script for atomicity
	for _, item := range items {
		stockKey := fmt.Sprintf("stock:%s", item.ID)
		newStock, err := s.redisClient.Decr(ctx, stockKey).Result()
		if err != nil {
			return err
		}
		if newStock < 0 {
			// Rollback (Compensating Transaction)
			s.redisClient.Incr(ctx, stockKey)
			return fmt.Errorf("out of stock for item %s", item.ID)
		}
	}

	// 3. Process Payment (Simulated)
	if err := processPayment(userID, items); err != nil {
		// Rollback Stock on Payment Failure
		for _, item := range items {
			s.redisClient.Incr(ctx, fmt.Sprintf("stock:%s", item.ID))
		}
		return fmt.Errorf("payment failed: %v", err)
	}

	// 4. Commit to Source of Truth (Cloud Spanner)
	// Using a Spanner ReadWriteTransaction
	_, err = s.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
		stmt := spanner.Statement{
			SQL: `INSERT INTO Orders (OrderId, UserId, Status, CreatedAt) VALUES (@orderId, @userId, 'CONFIRMED', PENDING_COMMIT_TIMESTAMP())`,
			Params: map[string]interface{}{
				"orderId": generateUUID(),
				"userId":  userID,
			},
		}
		return txn.Update(ctx, stmt)
	})

	if err != nil {
		// CRITICAL: Payment succeeded but DB failed. 
		// This enters the "Reaper" territory or requires manual reconciliation.
		// For now, log error for the background sweeper.
		logCriticalError("Order commit failed after payment", userID, err)
		return err
	}

	return nil
}
