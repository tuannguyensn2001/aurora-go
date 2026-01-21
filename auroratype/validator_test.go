package auroratype

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateConfig(t *testing.T) {
	t.Run("valid config", func(t *testing.T) {
		config := map[string]Parameter{
			"feature1": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Constraints: []Constraint{
							{Field: "country", Operator: "equal", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Empty(t, errs)
	})

	t.Run("empty config is valid", func(t *testing.T) {
		config := map[string]Parameter{}
		errs := ValidateConfig(config)
		assert.Empty(t, errs)
	})
}

func TestValidateParameter(t *testing.T) {
	t.Run("empty parameter name", func(t *testing.T) {
		config := map[string]Parameter{
			"": {
				DefaultValue: false,
				Rules:        []Rule{},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "parameter name cannot be empty")
	})
}

func TestValidateRule(t *testing.T) {
	t.Run("percentage out of range", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue:  true,
						Percentage:    intPtr(150),
						HashAttribute: strPtr("userID"),
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "must be between 0 and 100")
		assert.Equal(t, "testParam", errs[0].Parameter)
		assert.Equal(t, 0, errs[0].RuleIndex)
		assert.Equal(t, "percentage", errs[0].Field)
	})

	t.Run("percentage below zero", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue:  true,
						Percentage:    intPtr(-10),
						HashAttribute: strPtr("userID"),
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "must be between 0 and 100")
	})

	t.Run("percentage without hashAttribute", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Percentage:   intPtr(50),
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "is required when percentage is set")
		assert.Equal(t, "hashAttribute", errs[0].Field)
	})

	t.Run("percentage with empty hashAttribute", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue:  true,
						Percentage:    intPtr(50),
						HashAttribute: strPtr(""),
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "is required when percentage is set")
		assert.Equal(t, "hashAttribute", errs[0].Field)
	})
}

func TestValidateConstraint(t *testing.T) {
	t.Run("empty constraint field", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Constraints: []Constraint{
							{Field: "", Operator: "equal", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "cannot be empty")
		assert.Contains(t, errs[0].Field, "constraints[0].field")
	})

	t.Run("empty constraint operator", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Constraints: []Constraint{
							{Field: "country", Operator: "", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 1)
		assert.Contains(t, errs[0].Message, "cannot be empty")
		assert.Contains(t, errs[0].Field, "constraints[0].operator")
	})

	t.Run("multiple constraint errors", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Constraints: []Constraint{
							{Field: "", Operator: "", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 2)
	})
}

func TestValidateMultipleErrors(t *testing.T) {
	t.Run("multiple parameters with errors", func(t *testing.T) {
		config := map[string]Parameter{
			"param1": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Percentage:   intPtr(150),
					},
				},
			},
			"param2": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue: true,
						Constraints: []Constraint{
							{Field: "", Operator: "equal", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 3)
	})

	t.Run("multiple errors in single parameter", func(t *testing.T) {
		config := map[string]Parameter{
			"testParam": {
				DefaultValue: false,
				Rules: []Rule{
					{
						RolloutValue:  true,
						Percentage:    intPtr(150),
						HashAttribute: strPtr(""),
						Constraints: []Constraint{
							{Field: "", Operator: "", Value: "US"},
						},
					},
				},
			},
		}
		errs := ValidateConfig(config)
		assert.Len(t, errs, 4)
	})
}

func TestValidationErrorMessages(t *testing.T) {
	config := map[string]Parameter{
		"featureFlag": {
			DefaultValue: false,
			Rules: []Rule{
				{
					RolloutValue: true,
					Percentage:   intPtr(150),
					Constraints: []Constraint{
						{Field: "", Operator: "equal", Value: "US"},
					},
				},
			},
		},
	}
	errs := ValidateConfig(config)
	assert.Len(t, errs, 3)

	var errorMsgs []string
	for _, err := range errs {
		errorMsgs = append(errorMsgs, err.Error())
	}

	for _, msg := range errorMsgs {
		assert.Contains(t, msg, "featureFlag")
	}
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
