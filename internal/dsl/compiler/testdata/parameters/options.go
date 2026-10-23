package parameters

type config struct{ offset int }
type Option func(*config)

func add(value int) Option {
	return func(c *config) { c.offset += value }
}

func multiply(value int) Option {
	return func(c *config) { c.offset *= value }
}
