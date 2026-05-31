package errorcontrol

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckoutRecovery(t *testing.T) {
	receipts := &receiptFailures{}
	workflow := newCheckoutWorkflow(receipts)

	output, err := workflow.Run(context.Background(), checkout{
		orderID:  "order-1001",
		item:     "coffee-beans",
		quantity: 2,
	})

	require.NoError(t, err)
	assert.Equal(t, 2, output.reservationAttempts)
	assert.Equal(t, 2, output.reservedUnits)
	assert.Equal(t, 5000, output.chargedCents)
	assert.True(t, output.receiptQueued)
	assert.Equal(t, []string{
		"reserved 2 units of coffee-beans",
		"rolled back 2 units of coffee-beans",
		"reserved 2 units of coffee-beans",
		"charged 5000 cents",
		"queued receipt email",
	}, output.events)
	assert.Equal(t, 1, receipts.count)
	assert.EqualError(t, receipts.last, "order order-1001 receipt failed after payment: email provider unavailable")
}
