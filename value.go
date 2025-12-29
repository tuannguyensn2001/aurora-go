package aurora

type resolvedValue struct {
	value   any
	matched bool
}

func NewResolvedValue(value any, matched bool) *resolvedValue {
	return &resolvedValue{
		value:   value,
		matched: matched,
	}
}

func (r *resolvedValue) Boolean(defaultValue bool) bool {
	if !r.matched {
		return defaultValue
	}
	if r.value == nil {
		return defaultValue
	}

	val, ok := r.value.(bool)
	if !ok {
		return defaultValue
	}

	return val
}

func (r *resolvedValue) String(defaultValue string) string {
	if !r.matched {
		return defaultValue
	}

	if r.value == nil {
		return defaultValue
	}

	val, ok := r.value.(string)
	if !ok {
		return defaultValue
	}

	return val
}

func (r *resolvedValue) Int(defaultValue int) int {
	if !r.matched {
		return defaultValue
	}

	if r.value == nil {
		return defaultValue
	}

	val, ok := r.value.(int)
	if !ok {
		return defaultValue
	}

	return val
}

func (r *resolvedValue) Float(defaultValue float64) float64 {
	if !r.matched {
		return defaultValue
	}

	if r.value == nil {
		return defaultValue
	}

	val, ok := r.value.(float64)
	if !ok {
		return defaultValue
	}

	return val
}
