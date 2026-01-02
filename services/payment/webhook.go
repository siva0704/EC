package payment

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"cloud.google.com/go/spanner"
    "google.golang.org/api/iterator"
)

type PaymentService struct {
	spannerClient *spanner.Client
}

type PaymentWebhookPayload struct {
    Signature string `json:"signature"`
    OrderID   string `json:"order_id"`
    Status    string `json:"status"` // SUCCESS, FAILED
    Amount    float64 `json:"amount"`
}

// HandleProviderWebhook processes callbacks from Stripe/Adyen
func (s *PaymentService) HandleProviderWebhook(w http.ResponseWriter, r *http.Request) {
    var payload PaymentWebhookPayload
    if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
        http.Error(w, "Invalid Payload", http.StatusBadRequest)
        return
    }

    // 1. Verify Signature (Security Critical)
    if !verifySignature(payload.Signature) {
        http.Error(w, "Invalid Signature", http.StatusUnauthorized)
        return
    }

    // 2. Idempotent Update to Spanner
    ctx := context.Background()
    _, err := s.spannerClient.ReadWriteTransaction(ctx, func(ctx context.Context, txn *spanner.ReadWriteTransaction) error {
        // A. Check current status
        row, err := txn.ReadRow(ctx, "Orders", spanner.Key{payload.OrderID}, []string{"Status"})
        if err != nil {
            return err
        }
        var currentStatus string
        if err := row.Column(0, &currentStatus); err != nil {
            return err
        }

        // B. Idempotency Check
        if currentStatus == "PAID" || currentStatus == "FAILED" {
            // Already processed, return success to provider to stop retries
            return nil
        }

        // C. Update Status
        newStatus := "FAILED"
        if payload.Status == "SUCCESS" {
            newStatus = "PAID"
        }

        return txn.BufferWrite([]*spanner.Mutation{
            spanner.Update("Orders", []string{"OrderId", "Status"}, []interface{}{payload.OrderID, newStatus}),
        })
    })

    if err != nil {
        if spanner.ErrCode(err) == 6 { // NotFound
             http.Error(w, "Order Not Found", http.StatusNotFound)
             return
        }
        http.Error(w, "Database Error", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

func verifySignature(sig string) bool {
    // Real implementation would check HMAC header against SECRET
    return true
}
