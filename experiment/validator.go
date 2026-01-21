package experiment

import (
	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

type ValidationError struct {
	Experiment string
	Field      string
	Message    string
}

func (e ValidationError) Error() string {
	return "experiment \"" + e.Experiment + "\"." + e.Field + ": " + e.Message
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

func ValidateExperiments(experiments []auroratype.Experiment) []ValidationError {
	var errors []ValidationError
	for _, exp := range experiments {
		errors = append(errors, validateExperiment(exp)...)
	}
	return errors
}

func validateExperiment(exp auroratype.Experiment) []ValidationError {
	var errors []ValidationError

	if exp.ID == "" {
		errors = append(errors, ValidationError{
			Experiment: exp.Name,
			Field:      "id",
			Message:    "cannot be empty",
		})
	}

	if exp.Name == "" {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "name",
			Message:    "cannot be empty",
		})
	}

	if len(exp.Parameters) == 0 {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "parameters",
			Message:    "cannot be empty",
		})
	}

	if exp.HashAttribute == "" {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "hashAttribute",
			Message:    "cannot be empty",
		})
	}

	if exp.PopulationSize < 0 || exp.PopulationSize > 100 {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "populationSize",
			Message:    "must be between 0 and 100",
		})
	}

	if len(exp.Variants) == 0 {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "variants",
			Message:    "cannot be empty",
		})
	}

	totalRollout := 0
	for i, variant := range exp.Variants {
		totalRollout += variant.Rollout

		if variant.Key == "" {
			errors = append(errors, ValidationError{
				Experiment: exp.ID,
				Field:      "variants[" + string(rune('0'+i)) + "].key",
				Message:    "cannot be empty",
			})
		}

		if variant.Rollout < 0 || variant.Rollout > 100 {
			errors = append(errors, ValidationError{
				Experiment: exp.ID,
				Field:      "variants[" + string(rune('0'+i)) + "].rollout",
				Message:    "must be between 0 and 100",
			})
		}

		if len(variant.Values) == 0 {
			errors = append(errors, ValidationError{
				Experiment: exp.ID,
				Field:      "variants[" + string(rune('0'+i)) + "].values",
				Message:    "cannot be empty",
			})
		}
	}

	if totalRollout != 100 {
		errors = append(errors, ValidationError{
			Experiment: exp.ID,
			Field:      "variants",
			Message:    "total rollout must equal 100, got " + string(rune('0'+totalRollout)),
		})
	}

	return errors
}
