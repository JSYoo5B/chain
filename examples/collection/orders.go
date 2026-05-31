package collection

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/JSYoo5B/chain"
)

type orderLine struct {
	id         string
	item       string
	quantity   int
	unitCents  int
	totalCents int
}

func newOrderWorkflow(catalog map[string]int) *chain.Workflow[[]orderLine] {
	normalize := chain.AsSequenceSliceAction("NormalizeOrders", normalizeLine(), true)
	price := chain.AsParallelSliceAction("PriceOrders", priceLine(catalog))

	return chain.NewWorkflow("OrderPricing", normalize, price)
}

func newCapacityAction() chain.Action[map[string]int] {
	return chain.AsParallelMapAction[string]("ReserveCapacity", reserveUnits())
}

func normalizeLine() chain.Action[orderLine] {
	normalFunc := func(_ context.Context, input orderLine) (orderLine, error) {
		input.id = strings.TrimSpace(input.id)
		input.item = strings.ToLower(strings.TrimSpace(input.item))

		if input.id == "" {
			return input, errors.New("order id is required")
		}
		if input.quantity <= 0 {
			return input, fmt.Errorf("order %s has invalid quantity", input.id)
		}
		return input, nil
	}
	return chain.NewSimpleAction("NormalizeLine", normalFunc)
}

func priceLine(catalog map[string]int) chain.Action[orderLine] {
	priceFunc := func(_ context.Context, input orderLine) (orderLine, error) {
		unitCents, exists := catalog[input.item]
		if !exists {
			return input, fmt.Errorf("item %s is not in catalog", input.item)
		}

		input.unitCents = unitCents
		input.totalCents = input.quantity * unitCents
		return input, nil
	}
	return chain.NewSimpleAction("PriceLine", priceFunc)
}

func reserveUnits() chain.Action[int] {
	reserveFunc := func(_ context.Context, input int) (int, error) {
		if input < 0 {
			return input, errors.New("unit count cannot be negative")
		}
		return input + 5, nil
	}
	return chain.NewSimpleAction("ReserveUnits", reserveFunc)
}
