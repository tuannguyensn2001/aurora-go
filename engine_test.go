package core

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/core/evaluator"
)

func TestNewEngine(t *testing.T) {
	e := newEngine()
	assert.NotNil(t, e)
	assert.NotNil(t, e.operators)
	assert.Len(t, e.operators, 0) // No operators registered yet
}

func TestEngineBootstrap(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	assert.Contains(t, e.operators, evaluator.Equal)
	assert.Contains(t, e.operators, evaluator.NotEqual)
	assert.Contains(t, e.operators, evaluator.GreaterThan)
	assert.Contains(t, e.operators, evaluator.LessThan)
	assert.Contains(t, e.operators, evaluator.GreaterThanOrEqual)
	assert.Contains(t, e.operators, evaluator.LessThanOrEqual)
	assert.Contains(t, e.operators, evaluator.Contains)
	assert.Contains(t, e.operators, evaluator.In)
	assert.Contains(t, e.operators, evaluator.NotIn)
}

func TestEngineRegisterOperator(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	customOperator := func(a, b any) bool {
		s1, ok1 := a.(string)
		s2, ok2 := b.(string)
		if !ok1 || !ok2 {
			return false
		}
		return len(s1) > len(s2)
	}

	e.registerOperator(evaluator.Operator("longerThan"), customOperator)
	assert.Contains(t, e.operators, evaluator.Operator("longerThan"))
	assert.True(t, e.operators[evaluator.Operator("longerThan")]("hello world", "hi"))
	assert.False(t, e.operators[evaluator.Operator("longerThan")]("hi", "hello world"))
}

func TestEvaluateParameterNoRules(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules:        []auroratype.Rule{},
	}

	ctx := context.Background()
	attr := NewAttribute()

	result := e.evaluateParameter(ctx, "test_param", param, attr)

	assert.False(t, result.matched)
	assert.Equal(t, "default", result.value)
}

func TestEvaluateParameterSingleRuleMatching(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "value1",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "production"},
				},
			},
		},
	}

	ctx := context.Background()
	attr := NewAttribute()
	attr.Set("env", "production")

	result := e.evaluateParameter(ctx, "test_param", param, attr)

	assert.True(t, result.matched)
	assert.Equal(t, "value1", result.value)
}

func TestEvaluateParameterSingleRuleNotMatching(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "value1",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "production"},
				},
			},
		},
	}

	ctx := context.Background()
	attr := NewAttribute()
	attr.Set("env", "development")

	result := e.evaluateParameter(ctx, "test_param", param, attr)

	assert.False(t, result.matched)
	assert.Equal(t, "default", result.value)
}

func TestEvaluateParameterMultipleRulesFirstMatch(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "value1",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "production"},
				},
			},
			{
				RolloutValue: "value2",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "development"},
				},
			},
		},
	}

	ctx := context.Background()
	attr := NewAttribute()
	attr.Set("env", "production")

	result := e.evaluateParameter(ctx, "test_param", param, attr)

	assert.True(t, result.matched)
	assert.Equal(t, "value1", result.value)
}

func TestEvaluateParameterMultipleConstraintsAllMustMatch(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "value1",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "equal", Value: "production"},
					{Field: "version", Operator: "greaterThan", Value: 2},
				},
			},
		},
	}

	tests := []struct {
		name        string
		env         string
		version     int
		expectedVal string
		matched     bool
	}{
		{"both match", "production", 3, "value1", true},
		{"env match, version not", "production", 1, "default", false},
		{"env not match", "development", 3, "default", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			attr := NewAttribute()
			attr.Set("env", tt.env)
			attr.Set("version", tt.version)

			result := e.evaluateParameter(context.Background(), "test_param", param, attr)

			assert.Equal(t, tt.matched, result.matched)
			assert.Equal(t, tt.expectedVal, result.value)
		})
	}
}

func TestEvaluateRuleTimeBased(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	pastTime := time.Now().Add(-time.Hour).Unix()
	futureTime := time.Now().Add(time.Hour).Unix()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "past_rule",
				EffectiveAt:  &pastTime,
			},
			{
				RolloutValue: "future_rule",
				EffectiveAt:  &futureTime,
			},
		},
	}

	t.Run("rule with past effectiveAt matches", func(t *testing.T) {
		attr := NewAttribute()
		result := e.evaluateParameter(context.Background(), "test_param", param, attr)
		assert.True(t, result.matched)
		assert.Equal(t, "past_rule", result.value)
	})

	t.Run("rule with future effectiveAt does not match", func(t *testing.T) {
		futureParam := auroratype.Parameter{
			DefaultValue: "default",
			Rules: []auroratype.Rule{
				{
					RolloutValue: "future_rule",
					EffectiveAt:  &futureTime,
				},
			},
		}
		attr := NewAttribute()
		result := e.evaluateParameter(context.Background(), "test_param", futureParam, attr)
		assert.False(t, result.matched)
		assert.Equal(t, "default", result.value)
	})
}

func TestEvaluateRulePercentage(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	percentage := 100
	hashAttr := "user_id"

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue:  "rollout_value",
				Percentage:    &percentage,
				HashAttribute: &hashAttr,
				Constraints:   []auroratype.Constraint{},
			},
		},
	}

	t.Run("100% rollout should always match", func(t *testing.T) {
		attr := NewAttribute()
		attr.Set("user_id", "any_user")

		result := e.evaluateParameter(context.Background(), "test_param", param, attr)

		assert.True(t, result.matched)
		assert.Equal(t, "rollout_value", result.value)
	})

	t.Run("0% rollout should never match", func(t *testing.T) {
		zeroPercent := 0
		zeroParam := auroratype.Parameter{
			DefaultValue: "default",
			Rules: []auroratype.Rule{
				{
					RolloutValue:  "rollout_value",
					Percentage:    &zeroPercent,
					HashAttribute: &hashAttr,
					Constraints:   []auroratype.Constraint{},
				},
			},
		}
		attr := NewAttribute()
		attr.Set("user_id", "any_user")

		result := e.evaluateParameter(context.Background(), "test_param", zeroParam, attr)

		assert.False(t, result.matched)
		assert.Equal(t, "default", result.value)
	})

	t.Run("missing hash attribute should not match", func(t *testing.T) {
		attr := NewAttribute()

		result := e.evaluateParameter(context.Background(), "test_param", param, attr)

		assert.False(t, result.matched)
		assert.Equal(t, "default", result.value)
	})
}

func TestEvaluateRuleWithUnknownOperator(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := auroratype.Parameter{
		DefaultValue: "default",
		Rules: []auroratype.Rule{
			{
				RolloutValue: "value1",
				Constraints: []auroratype.Constraint{
					{Field: "env", Operator: "unknownOperator", Value: "production"},
				},
			},
		},
	}

	ctx := context.Background()
	attr := NewAttribute()
	attr.Set("env", "production")

	result := e.evaluateParameter(ctx, "test_param", param, attr)

	assert.False(t, result.matched)
	assert.Equal(t, "default", result.value)
}

func TestCalculateHashPercentage(t *testing.T) {
	tests := []struct {
		name       string
		percentage int
		expected   bool
	}{
		{"0% returns false", 0, false},
		{"100% returns true", 100, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash := evaluator.CalculateHash("test_value", "test_param")
			result := evaluator.IsInPercentageRange(hash, tt.percentage)
			assert.Equal(t, tt.expected, result)
		})
	}

	t.Run("same value always returns same result", func(t *testing.T) {
		hash1 := evaluator.CalculateHash("test_value", "test_param")
		hash2 := evaluator.CalculateHash("test_value", "test_param")
		result1 := evaluator.IsInPercentageRange(hash1, 50)
		result2 := evaluator.IsInPercentageRange(hash2, 50)
		assert.Equal(t, result1, result2)
	})

	t.Run("different values can have different results", func(t *testing.T) {
		results := make(map[bool]int)
		for i := 0; i < 100; i++ {
			hash := evaluator.CalculateHash("user_"+string(rune('a'+i)), "test_param")
			result := evaluator.IsInPercentageRange(hash, 50)
			results[result]++
		}
		assert.True(t, results[true] > 0)
		assert.True(t, results[false] > 0)
	})

	t.Run("different parameters have independent rollouts", func(t *testing.T) {
		hash1 := evaluator.CalculateHash("user1", "param1")
		hash2 := evaluator.CalculateHash("user1", "param2")
		result1 := evaluator.IsInPercentageRange(hash1, 100)
		result2 := evaluator.IsInPercentageRange(hash2, 0)
		assert.True(t, result1)
		assert.False(t, result2)
	})
}

func TestEngineEvaluateRuleComplex(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	t.Run("constraints must pass before percentage rollout", func(t *testing.T) {
		percentage := 100
		hashAttr := "user_id"

		param := auroratype.Parameter{
			DefaultValue: "default",
			Rules: []auroratype.Rule{
				{
					RolloutValue:  "rollout_value",
					Percentage:    &percentage,
					HashAttribute: &hashAttr,
					Constraints: []auroratype.Constraint{
						{Field: "env", Operator: "equal", Value: "production"},
					},
				},
			},
		}

		t.Run("constraint fails, returns default", func(t *testing.T) {
			attr := NewAttribute()
			attr.Set("user_id", "user1")
			attr.Set("env", "development")

			result := e.evaluateParameter(context.Background(), "test_param", param, attr)

			assert.False(t, result.matched)
			assert.Equal(t, "default", result.value)
		})

		t.Run("constraint passes, returns rollout value", func(t *testing.T) {
			attr := NewAttribute()
			attr.Set("user_id", "user1")
			attr.Set("env", "production")

			result := e.evaluateParameter(context.Background(), "test_param", param, attr)

			assert.True(t, result.matched)
			assert.Equal(t, "rollout_value", result.value)
		})
	})
}

func TestEngineWithDifferentConstraintOperators(t *testing.T) {
	e := newEngine()
	e.bootstrap()

	param := func(constraints []auroratype.Constraint) auroratype.Parameter {
		return auroratype.Parameter{
			DefaultValue: "default",
			Rules: []auroratype.Rule{
				{
					RolloutValue: "matched",
					Constraints:  constraints,
				},
			},
		}
	}

	tests := []struct {
		name       string
		constraint auroratype.Constraint
		attrVal    any
		expected   bool
	}{
		{"equal - matches", auroratype.Constraint{Field: "key", Operator: "equal", Value: "test"}, "test", true},
		{"equal - no match", auroratype.Constraint{Field: "key", Operator: "equal", Value: "test"}, "other", false},
		{"notEqual - matches", auroratype.Constraint{Field: "key", Operator: "notEqual", Value: "test"}, "other", true},
		{"notEqual - no match", auroratype.Constraint{Field: "key", Operator: "notEqual", Value: "test"}, "test", false},
		{"greaterThan - matches", auroratype.Constraint{Field: "key", Operator: "greaterThan", Value: 10}, 20, true},
		{"greaterThan - no match", auroratype.Constraint{Field: "key", Operator: "greaterThan", Value: 10}, 5, false},
		{"lessThan - matches", auroratype.Constraint{Field: "key", Operator: "lessThan", Value: 10}, 5, true},
		{"lessThan - no match", auroratype.Constraint{Field: "key", Operator: "lessThan", Value: 10}, 20, false},
		{"contains - matches", auroratype.Constraint{Field: "key", Operator: "contains", Value: "world"}, "hello world", true},
		{"contains - no match", auroratype.Constraint{Field: "key", Operator: "contains", Value: "world"}, "hello", false},
		{"in - matches", auroratype.Constraint{Field: "key", Operator: "in", Value: []string{"a", "b", "c"}}, "b", true},
		{"in - no match", auroratype.Constraint{Field: "key", Operator: "in", Value: []string{"a", "b", "c"}}, "d", false},
		{"notIn - matches", auroratype.Constraint{Field: "key", Operator: "notIn", Value: []string{"a", "b", "c"}}, "d", true},
		{"notIn - no match", auroratype.Constraint{Field: "key", Operator: "notIn", Value: []string{"a", "b", "c"}}, "b", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := param([]auroratype.Constraint{tt.constraint})
			attr := NewAttribute()
			attr.Set("key", tt.attrVal)

			result := e.evaluateParameter(context.Background(), "test_param", p, attr)

			assert.Equal(t, tt.expected, result.matched)
			if tt.expected {
				assert.Equal(t, "matched", result.value)
			} else {
				assert.Equal(t, "default", result.value)
			}
		})
	}
}
