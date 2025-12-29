package aurora

import "reflect"

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

func equalOperator(a, b any) bool {
	// Handle nil cases
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Get reflect values
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)

	// Get types
	ta := va.Type()
	tb := vb.Type()

	// If types are exactly the same, use reflect.DeepEqual
	if ta == tb {
		return reflect.DeepEqual(a, b)
	}

	// Handle numeric type conversions
	if isNumeric(ta) && isNumeric(tb) {
		return compareNumeric(va, vb)
	}

	// Types don't match and not both numeric, return false
	return false
}

// isNumeric checks if a type is numeric
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

// compareNumeric compares two numeric values by converting to float64
func compareNumeric(va, vb reflect.Value) bool {
	fa, ok1 := getFloat64(va)
	fb, ok2 := getFloat64(vb)
	if !ok1 || !ok2 {
		return false
	}
	return fa == fb
}

// getFloat64 converts a reflect.Value to float64
func getFloat64(v reflect.Value) (float64, bool) {
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

// notEqualOperator checks if a != b
func notEqualOperator(a, b any) bool {
	return !equalOperator(a, b)
}

// greaterThanOperator checks if a > b
func greaterThanOperator(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	// Handle numeric comparisons
	if isNumeric(ta) && isNumeric(tb) {
		fa, ok1 := getFloat64(va)
		fb, ok2 := getFloat64(vb)
		if !ok1 || !ok2 {
			return false
		}
		return fa > fb
	}

	// Handle string comparisons
	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return va.String() > vb.String()
	}

	return false
}

// lessThanOperator checks if a < b
func lessThanOperator(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	// Handle numeric comparisons
	if isNumeric(ta) && isNumeric(tb) {
		fa, ok1 := getFloat64(va)
		fb, ok2 := getFloat64(vb)
		if !ok1 || !ok2 {
			return false
		}
		return fa < fb
	}

	// Handle string comparisons
	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return va.String() < vb.String()
	}

	return false
}

// greaterThanOrEqualOperator checks if a >= b
func greaterThanOrEqualOperator(a, b any) bool {
	return greaterThanOperator(a, b) || equalOperator(a, b)
}

// lessThanOrEqualOperator checks if a <= b
func lessThanOrEqualOperator(a, b any) bool {
	return lessThanOperator(a, b) || equalOperator(a, b)
}

// containsOperator checks if string a contains string b, or if slice/array a contains b
func containsOperator(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	ta := va.Type()
	tb := vb.Type()

	// Handle string contains
	if ta.Kind() == reflect.String && tb.Kind() == reflect.String {
		return contains(va.String(), vb.String())
	}

	// Handle slice/array contains
	if (ta.Kind() == reflect.Slice || ta.Kind() == reflect.Array) && ta.Elem() == tb {
		for i := 0; i < va.Len(); i++ {
			if reflect.DeepEqual(va.Index(i).Interface(), b) {
				return true
			}
		}
	}

	return false
}

// contains checks if string s contains substring substr
func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(s) < len(substr) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// inOperator checks if value a is in array/slice b
func inOperator(a, b any) bool {
	if a == nil || b == nil {
		return false
	}

	vb := reflect.ValueOf(b)
	tb := vb.Type()

	// b must be a slice or array
	if tb.Kind() != reflect.Slice && tb.Kind() != reflect.Array {
		return false
	}

	// Iterate through the array/slice
	for i := 0; i < vb.Len(); i++ {
		elem := vb.Index(i)

		// Use equalOperator logic for comparison (handles type conversions)
		if equalOperator(a, elem.Interface()) {
			return true
		}
	}

	return false
}

// notInOperator checks if value a is NOT in array/slice b
func notInOperator(a, b any) bool {
	return !inOperator(a, b)
}
