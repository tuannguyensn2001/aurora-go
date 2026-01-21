package auroratype

type ValidationError struct {
	Parameter string
	RuleIndex int
	Field     string
	Message   string
}

func (e ValidationError) Error() string {
	if e.RuleIndex >= 0 {
		return "parameter \"" + e.Parameter + "\"." + e.Field + ": " + e.Message
	}
	return "parameter \"" + e.Parameter + "\": " + e.Message
}

type ValidationErrors struct {
	Errors []ValidationError
}

func (e ValidationErrors) Error() string {
	if len(e.Errors) == 0 {
		return "validation errors"
	}
	msg := "validation errors:\n"
	for _, err := range e.Errors {
		msg += "  - " + err.Error() + "\n"
	}
	return msg[:len(msg)-1]
}

func ValidateConfig(config map[string]Parameter) []ValidationError {
	var errors []ValidationError
	for name, param := range config {
		errors = append(errors, validateParameter(name, param)...)
	}
	return errors
}

func validateParameter(name string, param Parameter) []ValidationError {
	var errors []ValidationError

	if name == "" {
		errors = append(errors, ValidationError{
			Parameter: name,
			Field:     "",
			Message:   "parameter name cannot be empty",
		})
	}

	for i, rule := range param.Rules {
		errors = append(errors, validateRule(name, i, rule)...)
	}

	return errors
}

func validateRule(paramName string, ruleIndex int, rule Rule) []ValidationError {
	var errors []ValidationError

	if rule.Percentage != nil {
		if *rule.Percentage < 0 || *rule.Percentage > 100 {
			errors = append(errors, ValidationError{
				Parameter: paramName,
				RuleIndex: ruleIndex,
				Field:     "percentage",
				Message:   "must be between 0 and 100",
			})
		}

		if rule.HashAttribute == nil || *rule.HashAttribute == "" {
			errors = append(errors, ValidationError{
				Parameter: paramName,
				RuleIndex: ruleIndex,
				Field:     "hashAttribute",
				Message:   "is required when percentage is set",
			})
		}
	}

	for i, constraint := range rule.Constraints {
		errors = append(errors, validateConstraint(paramName, ruleIndex, i, constraint)...)
	}

	return errors
}

func validateConstraint(paramName string, ruleIndex, constraintIndex int, constraint Constraint) []ValidationError {
	var errors []ValidationError

	if constraint.Field == "" {
		errors = append(errors, ValidationError{
			Parameter: paramName,
			RuleIndex: ruleIndex,
			Field:     "constraints[" + string(rune('0'+constraintIndex)) + "].field",
			Message:   "cannot be empty",
		})
	}

	if constraint.Operator == "" {
		errors = append(errors, ValidationError{
			Parameter: paramName,
			RuleIndex: ruleIndex,
			Field:     "constraints[" + string(rune('0'+constraintIndex)) + "].operator",
			Message:   "cannot be empty",
		})
	}

	return errors
}
