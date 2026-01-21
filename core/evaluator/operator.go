package evaluator

import (
	"fmt"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

type OperatorFunc func(a, b any) bool

func Equal(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return fmt.Sprintf("%v", a) == fmt.Sprintf("%v", b)
}

func NotEqual(a, b any) bool {
	return !Equal(a, b)
}

func GreaterThan(a, b any) bool {
	af, ok1 := ToFloat64(a)
	bf, ok2 := ToFloat64(b)
	if !ok1 || !ok2 {
		return false
	}
	return af > bf
}

func LessThan(a, b any) bool {
	af, ok1 := ToFloat64(a)
	bf, ok2 := ToFloat64(b)
	if !ok1 || !ok2 {
		return false
	}
	return af < bf
}

func GreaterThanOrEqual(a, b any) bool {
	return GreaterThan(a, b) || Equal(a, b)
}

func LessThanOrEqual(a, b any) bool {
	return LessThan(a, b) || Equal(a, b)
}

func Contains(a, b any) bool {
	as, ok1 := a.(string)
	bs, ok2 := b.(string)
	if !ok1 || !ok2 {
		return false
	}
	return len(as) >= len(bs) && (as == bs || containsSubstring(as, bs))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func In(a, b any) bool {
	arr, ok := b.([]any)
	if !ok {
		return false
	}
	for _, elem := range arr {
		if Equal(a, elem) {
			return true
		}
	}
	return false
}

func NotIn(a, b any) bool {
	return !In(a, b)
}

func ToFloat64(v any) (float64, bool) {
	switch val := v.(type) {
	case int:
		return float64(val), true
	case int8:
		return float64(val), true
	case int16:
		return float64(val), true
	case int32:
		return float64(val), true
	case int64:
		return float64(val), true
	case uint:
		return float64(val), true
	case uint8:
		return float64(val), true
	case uint16:
		return float64(val), true
	case uint32:
		return float64(val), true
	case uint64:
		return float64(val), true
	case float32:
		return float64(val), true
	case float64:
		return val, true
	default:
		return 0, false
	}
}

var DefaultOperators = map[string]OperatorFunc{
	"equal":              Equal,
	"notEqual":           NotEqual,
	"greaterThan":        GreaterThan,
	"lessThan":           LessThan,
	"greaterThanOrEqual": GreaterThanOrEqual,
	"lessThanOrEqual":    LessThanOrEqual,
	"contains":           Contains,
	"in":                 In,
	"notIn":              NotIn,
}

func EvaluateConstraint(constraint auroratype.Constraint, attr map[string]any, operators map[string]OperatorFunc) bool {
	if operators == nil {
		operators = DefaultOperators
	}

	op := operators[constraint.Operator]
	if op == nil {
		return false
	}

	attrValue := attr[constraint.Field]
	return op(attrValue, constraint.Value)
}
