package evaluator

import (
	"math"
	"reflect"
	"strings"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

type Operator string

const (
	Equal              Operator = "equal"
	NotEqual           Operator = "notEqual"
	GreaterThan        Operator = "greaterThan"
	LessThan           Operator = "lessThan"
	GreaterThanOrEqual Operator = "greaterThanOrEqual"
	LessThanOrEqual    Operator = "lessThanOrEqual"
	Contains           Operator = "contains"
	In                 Operator = "in"
	NotIn              Operator = "notIn"
)

const epsilon = 1e-9

func EqualOp(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	ta := va.Type()
	tb := vb.Type()

	if ta == tb {
		return reflect.DeepEqual(a, b)
	}

	if isNumeric(ta) && isNumeric(tb) {
		return compareNumeric(va, vb)
	}

	return false
}

func NotEqualOp(a, b any) bool {
	return !EqualOp(a, b)
}

func GreaterThanOp(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	if isNumeric(ta) && isNumeric(tb) {
		fa, ok1 := toFloat64(va)
		fb, ok2 := toFloat64(vb)
		if !ok1 || !ok2 {
			return false
		}
		return fa > fb
	}

	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return va.String() > vb.String()
	}

	return false
}

func LessThanOp(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	if isNumeric(ta) && isNumeric(tb) {
		fa, ok1 := toFloat64(va)
		fb, ok2 := toFloat64(vb)
		if !ok1 || !ok2 {
			return false
		}
		return fa < fb
	}

	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return va.String() < vb.String()
	}

	return false
}

func GreaterThanOrEqualOp(a, b any) bool {
	return GreaterThanOp(a, b) || EqualOp(a, b)
}

func LessThanOrEqualOp(a, b any) bool {
	return LessThanOp(a, b) || EqualOp(a, b)
}

func ContainsOp(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return strings.Contains(va.String(), vb.String())
	}

	if (ta.Kind() == reflect.Slice || ta.Kind() == reflect.Array) && ta.Elem() == tb {
		for i := 0; i < va.Len(); i++ {
			if reflect.DeepEqual(va.Index(i).Interface(), b) {
				return true
			}
		}
	}

	return false
}

func InOp(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	vb := reflect.ValueOf(b)
	tb := vb.Type()

	if tb.Kind() != reflect.Slice && tb.Kind() != reflect.Array {
		return false
	}

	for i := 0; i < vb.Len(); i++ {
		elem := vb.Index(i)

		if EqualOp(a, elem.Interface()) {
			return true
		}
	}

	return false
}

func NotInOp(a, b any) bool {
	return !InOp(a, b)
}

func isNumeric(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64:
		return true
	default:
		return false
	}
}

func toFloat64(v reflect.Value) (float64, bool) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return float64(v.Int()), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return float64(v.Uint()), true
	case reflect.Float32, reflect.Float64:
		return v.Float(), true
	default:
		return 0, false
	}
}

func compareNumeric(va, vb reflect.Value) bool {
	fa, ok1 := toFloat64(va)
	fb, ok2 := toFloat64(vb)
	if !ok1 || !ok2 {
		return false
	}
	return math.Abs(fa-fb) < epsilon
}

var DefaultOperators = map[Operator]func(a, b any) bool{
	Equal:              EqualOp,
	NotEqual:           NotEqualOp,
	GreaterThan:        GreaterThanOp,
	LessThan:           LessThanOp,
	GreaterThanOrEqual: GreaterThanOrEqualOp,
	LessThanOrEqual:    LessThanOrEqualOp,
	Contains:           ContainsOp,
	In:                 InOp,
	NotIn:              NotInOp,
}

func EvaluateConstraint(constraint auroratype.Constraint, attr map[string]any, operators map[Operator]func(a, b any) bool) bool {
	if operators == nil {
		operators = DefaultOperators
	}

	op := operators[Operator(constraint.Operator)]
	if op == nil {
		return false
	}

	attrValue := attr[constraint.Field]
	return op(attrValue, constraint.Value)
}
