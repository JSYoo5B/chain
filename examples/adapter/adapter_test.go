package adapter

import (
	"context"
	"testing"

	"github.com/JSYoo5B/chain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCartAdapter(t *testing.T) {
	workflow := chain.NewWorkflow(
		"CartAdapter",
		customerToCart(newNormalizeCustomerAction("NormalizeCustomer")),
		totalCentsToCart(newAddPriceAction("AddCoffeePrice", 1200)),
		totalCentsToCart(newAddPriceAction("AddFilterPrice", 300)),
	)

	output, err := workflow.Run(context.Background(), cart{
		customer:   "  Ada   Lovelace ",
		totalCents: 0,
	})

	require.NoError(t, err)
	assert.Equal(t, "Ada Lovelace", output.customer)
	assert.Equal(t, 1500, output.totalCents)
}

func TestWorkflowAdapter(t *testing.T) {
	addBundle := chain.NewWorkflow(
		"AddBundle",
		newAddPriceAction("AddCoffeePrice", 1200),
		newAddPriceAction("AddFilterPrice", 300),
		newAddPriceAction("AddMugPrice", 1500),
	)

	workflow := chain.NewWorkflow(
		"BundleCart",
		customerToCart(newNormalizeCustomerAction("NormalizeCustomer")),
		totalCentsToCart(addBundle),
	)

	output, err := workflow.Run(context.Background(), cart{
		customer:   "  Grace   Hopper ",
		totalCents: 500,
	})

	require.NoError(t, err)
	assert.Equal(t, "Grace Hopper", output.customer)
	assert.Equal(t, 3500, output.totalCents)
}
