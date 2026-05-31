package errorcontrol

import (
	"context"
	"errors"
	"fmt"

	"github.com/JSYoo5B/chain"
)

type checkout struct {
	orderID             string
	item                string
	quantity            int
	reservedUnits       int
	reservationAttempts int
	chargedCents        int
	receiptQueued       bool
	events              []string
}

type receiptFailures struct {
	count int
	last  error
}

func newCheckoutWorkflow(receipts *receiptFailures) *chain.Workflow[checkout] {
	reserve := chain.AsRetryableAction(
		"ReserveInventoryWithRetry",
		reserveInventory(),
		rollbackInventoryReservation(),
		2,
	)
	charge := chargePayment(2500)
	receipt := chain.AsBestEffortAction(
		queueReceipt(),
		func(_ context.Context, input checkout, err error) {
			receipts.count++
			receipts.last = fmt.Errorf("order %s receipt failed after payment: %w", input.orderID, err)
		},
	)

	return chain.NewWorkflow("Checkout", reserve, charge, receipt)
}

func reserveInventory() chain.Action[checkout] {
	return chain.NewSimpleAction("ReserveInventory", func(_ context.Context, input checkout) (checkout, error) {
		input.reservationAttempts++
		input.reservedUnits += input.quantity
		input.events = append(input.events, fmt.Sprintf("reserved %d units of %s", input.quantity, input.item))

		if input.reservationAttempts == 1 {
			return input, errors.New("inventory service timeout")
		}
		return input, nil
	})
}

func rollbackInventoryReservation() chain.Action[checkout] {
	return chain.NewSimpleAction("RollbackInventoryReservation", func(_ context.Context, input checkout) (checkout, error) {
		input.reservedUnits -= input.quantity
		input.events = append(input.events, fmt.Sprintf("rolled back %d units of %s", input.quantity, input.item))
		return input, nil
	})
}

func chargePayment(unitPriceCents int) chain.Action[checkout] {
	return chain.NewSimpleAction("ChargePayment", func(_ context.Context, input checkout) (checkout, error) {
		if input.reservedUnits < input.quantity {
			return input, errors.New("cannot charge before inventory is reserved")
		}

		input.chargedCents = input.quantity * unitPriceCents
		input.events = append(input.events, fmt.Sprintf("charged %d cents", input.chargedCents))
		return input, nil
	})
}

func queueReceipt() chain.Action[checkout] {
	return chain.NewSimpleAction("QueueReceipt", func(_ context.Context, input checkout) (checkout, error) {
		input.receiptQueued = true
		input.events = append(input.events, "queued receipt email")
		return input, errors.New("email provider unavailable")
	})
}
