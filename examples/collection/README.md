# Collection Actions in Chain: A Case Study with Order Pricing

Full implementation: [orders.go](orders.go)

A `chain.Action[T]` processes one value of type `T`. Real workflows often need
to apply the same action to a collection, such as a slice of order lines or a
map of regional capacities.

You could write a custom `Action[[]orderLine]` every time, but the logic is
usually the same:

- copy the collection
- run one item action per element
- preserve indexes or keys
- decide how to aggregate errors

Collection action builders provide that wrapper while keeping the item action
small and reusable.

## Order Line

This example prices a small list of order lines:

```go
type orderLine struct {
    id         string
    item       string
    quantity   int
    unitCents  int
    totalCents int
}
```

Two single-item actions do the actual work:

```go
func normalizeLine() chain.Action[orderLine]
func priceLine(catalog map[string]int) chain.Action[orderLine]
```

The workflow needs actions over `[]orderLine`, so each item action is wrapped
as a collection action.

## Sequential Slice Action

Normalization should happen in input order and stop before pricing if a line is
invalid. `AsSequenceSliceAction` turns `Action[orderLine]` into
`Action[[]orderLine]`.

```go
normalize := chain.AsSequenceSliceAction("NormalizeOrders", normalizeLine(), true)
```

The last argument is `stopOnError`. In this example it is `true`, so the first
invalid line stops the slice action.

```mermaid
graph LR;
    classDef inputNode fill:#f4f4f5,stroke:#52525b,stroke-width:1.5px;
    classDef wrapperNode fill:#ede9fe,stroke:#6d28d9,stroke-width:1.5px;
    classDef itemNode fill:#fff7ed,stroke:#c2410c,stroke-width:1.5px;
    classDef outputNode fill:#e0f2fe,stroke:#0369a1,stroke-width:1.5px;
    classDef errorNode fill:#fee2e2,stroke:#b91c1c,stroke-width:1.5px;

    input(["[]orderLine input"]):::inputNode
    wrapper["NormalizeOrders<br/>Action[[]orderLine]"]:::wrapperNode
    item0["index 0<br/>NormalizeLine"]:::itemNode
    item1["index 1<br/>NormalizeLine"]:::itemNode
    output(["[]orderLine output"]):::outputNode
    stop{{"stop on error"}}:::errorNode

    input --> wrapper
    wrapper --> item0
    item0 --> item1
    item1 --> output
    item1 -. invalid quantity .-> stop
```

`TestOrderValidation` verifies that an invalid quantity stops the workflow after
normalization, before the pricing action runs.

## Parallel Slice Action

After normalization, each line can be priced independently. `AsParallelSliceAction`
keeps the output indexes aligned with the input indexes.

```go
price := chain.AsParallelSliceAction("PriceOrders", priceLine(catalog))
```

```mermaid
graph LR;
    classDef inputNode fill:#f4f4f5,stroke:#52525b,stroke-width:1.5px;
    classDef wrapperNode fill:#ede9fe,stroke:#6d28d9,stroke-width:1.5px;
    classDef itemNode fill:#fff7ed,stroke:#c2410c,stroke-width:1.5px;
    classDef outputNode fill:#e0f2fe,stroke:#0369a1,stroke-width:1.5px;

    input(["normalized []orderLine"]):::inputNode
    wrapper["PriceOrders<br/>Action[[]orderLine]"]:::wrapperNode
    item0["index 0<br/>PriceLine"]:::itemNode
    item1["index 1<br/>PriceLine"]:::itemNode
    output(["priced []orderLine<br/>same indexes"]):::outputNode

    input --> wrapper
    wrapper --> item0
    wrapper --> item1
    item0 --> output
    item1 --> output
```

## Order Pricing Workflow

The final slice workflow is just two collection actions:

```go
func newOrderWorkflow(catalog map[string]int) *chain.Workflow[[]orderLine] {
    normalize := chain.AsSequenceSliceAction("NormalizeOrders", normalizeLine(), true)
    price := chain.AsParallelSliceAction("PriceOrders", priceLine(catalog))

    return chain.NewWorkflow("OrderPricing", normalize, price)
}
```

```mermaid
graph LR;
    classDef startNode fill:#cce5ff,stroke:#003366,stroke-width:1.5px;
    classDef sequenceNode fill:#e6ffe6,stroke:#339966,stroke-width:1.5px;
    classDef parallelNode fill:#fff7ed,stroke:#c2410c,stroke-width:1.5px;
    classDef endNode fill:#ffe6e6,stroke:#cc3333,stroke-width:1.5px;

    start((Start)):::startNode
    normalize[NormalizeOrders<br/>sequential slice]:::sequenceNode
    price[PriceOrders<br/>parallel slice]:::parallelNode
    terminate((End)):::endNode

    start --> normalize
    normalize --> price
    price --> terminate
```

`TestOrderWorkflow` verifies this path:

```text
raw order lines -> NormalizeOrders -> PriceOrders -> priced order lines
```

## Parallel Map Action

Map processing is separate from the order workflow so the map behavior stays
clear. `AsParallelMapAction` turns `Action[int]` into `Action[map[string]int]`.

```go
func newCapacityAction() chain.Action[map[string]int] {
    return chain.AsParallelMapAction[string]("ReserveCapacity", reserveUnits())
}
```

Each region value is processed independently, and the original keys are
preserved.

```mermaid
graph LR;
    classDef inputNode fill:#f4f4f5,stroke:#52525b,stroke-width:1.5px;
    classDef wrapperNode fill:#ede9fe,stroke:#6d28d9,stroke-width:1.5px;
    classDef itemNode fill:#fff7ed,stroke:#c2410c,stroke-width:1.5px;
    classDef outputNode fill:#e0f2fe,stroke:#0369a1,stroke-width:1.5px;

    input(["map[string]int<br/>KR: 2, US: 3"]):::inputNode
    wrapper["ReserveCapacity<br/>Action[map[string]int]"]:::wrapperNode
    kr["key KR<br/>ReserveUnits<br/>2 -> 7"]:::itemNode
    us["key US<br/>ReserveUnits<br/>3 -> 8"]:::itemNode
    output(["map[string]int<br/>KR: 7, US: 8"]):::outputNode

    input --> wrapper
    wrapper --> kr
    wrapper --> us
    kr --> output
    us --> output
```

`TestCapacityMap` verifies that keys are preserved while values are processed in
parallel.
