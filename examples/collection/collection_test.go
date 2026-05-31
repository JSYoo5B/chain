package collection

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderWorkflow(t *testing.T) {
	workflow := newOrderWorkflow(map[string]int{
		"coffee": 1200,
		"filter": 300,
	})

	output, err := workflow.Run(context.Background(), []orderLine{
		{id: " A-001 ", item: " Coffee ", quantity: 2},
		{id: "A-002", item: "filter", quantity: 3},
	})

	require.NoError(t, err)
	assert.Equal(t, []orderLine{
		{id: "A-001", item: "coffee", quantity: 2, unitCents: 1200, totalCents: 2400},
		{id: "A-002", item: "filter", quantity: 3, unitCents: 300, totalCents: 900},
	}, output)
}

func TestOrderValidation(t *testing.T) {
	workflow := newOrderWorkflow(map[string]int{"coffee": 1200})

	output, err := workflow.Run(context.Background(), []orderLine{
		{id: "A-001", item: "coffee", quantity: 2},
		{id: "A-002", item: "coffee", quantity: 0},
	})

	require.Error(t, err)
	assert.ErrorContains(t, err, "error occurred at index 1")
	assert.Equal(t, []orderLine{
		{id: "A-001", item: "coffee", quantity: 2},
		{id: "A-002", item: "coffee", quantity: 0},
	}, output)
}

func TestCapacityMap(t *testing.T) {
	action := newCapacityAction()

	output, err := action.Run(context.Background(), map[string]int{
		"KR": 2,
		"US": 3,
	})

	require.NoError(t, err)
	assert.Equal(t, map[string]int{"KR": 7, "US": 8}, output)
}
