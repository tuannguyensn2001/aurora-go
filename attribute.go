package aurora

type attribute struct {
	vals map[string]any
}

func NewAttribute() *attribute {
	return &attribute{
		vals: make(map[string]any),
	}
}

func (a *attribute) Set(key string, val any) {
	a.vals[key] = val
}

func (a *attribute) Get(key string) any {
	return a.vals[key]
}
