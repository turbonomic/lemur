package topology

type Commodity struct {
	Name  string
	Value float64
}

func newCommodity(name string, value float64) *Commodity {
	return &Commodity{name, value}
}
