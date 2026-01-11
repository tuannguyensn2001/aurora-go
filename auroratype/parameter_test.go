package auroratype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParameter(t *testing.T) {
	t.Run("default value only", func(t *testing.T) {
		p := Parameter{
			DefaultValue: "test_value",
			Rules:        nil,
		}

		assert.Equal(t, "test_value", p.DefaultValue)
		assert.Nil(t, p.Rules)
	})

	t.Run("with rules", func(t *testing.T) {
		p := Parameter{
			DefaultValue: "default",
			Rules: []Rule{
				{
					RolloutValue: "value1",
					Constraints:  []Constraint{},
				},
			},
		}

		assert.Equal(t, "default", p.DefaultValue)
		assert.Len(t, p.Rules, 1)
	})
}

func TestRule(t *testing.T) {
	t.Run("basic rule", func(t *testing.T) {
		r := Rule{
			RolloutValue: "rollout",
			Percentage:   nil,
			EffectiveAt:  nil,
			Constraints: []Constraint{
				{Field: "env", Operator: "equal", Value: "prod"},
			},
		}

		assert.Equal(t, "rollout", r.RolloutValue)
		assert.Nil(t, r.Percentage)
		assert.Nil(t, r.EffectiveAt)
		assert.Len(t, r.Constraints, 1)
	})

	t.Run("percentage rollout", func(t *testing.T) {
		percentage := 50
		hashAttr := "user_id"
		r := Rule{
			RolloutValue:  "rollout",
			Percentage:    &percentage,
			HashAttribute: &hashAttr,
			Constraints:   []Constraint{},
		}

		assert.Equal(t, 50, *r.Percentage)
		assert.Equal(t, "user_id", *r.HashAttribute)
	})

	t.Run("time-based rule", func(t *testing.T) {
		effectiveAt := int64(1700000000)
		r := Rule{
			RolloutValue: "rollout",
			EffectiveAt:  &effectiveAt,
			Constraints:  []Constraint{},
		}

		assert.Equal(t, int64(1700000000), *r.EffectiveAt)
	})
}

func TestConstraint(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		operator string
		value    any
	}{
		{"string constraint", "env", "equal", "production"},
		{"numeric constraint", "version", "greaterThan", 2},
		{"array constraint", "regions", "in", []string{"us", "eu"}},
		{"boolean constraint", "enabled", "equal", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := Constraint{
				Field:    tt.field,
				Operator: tt.operator,
				Value:    tt.value,
			}

			assert.Equal(t, tt.field, c.Field)
			assert.Equal(t, tt.operator, c.Operator)
			assert.Equal(t, tt.value, c.Value)
		})
	}
}
