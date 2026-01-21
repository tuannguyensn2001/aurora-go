package experiment

import (
	"context"
	"sort"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/core/evaluator"
)

type Evaluation struct {
	ExperimentID string
	VariantKey   string
	Values       map[string]interface{}
	Matched      bool
}

type Engine struct {
	operators map[evaluator.Operator]func(a, b any) bool
}

func NewEngine() *Engine {
	return &Engine{
		operators: make(map[evaluator.Operator]func(a, b any) bool),
	}
}

func (e *Engine) Bootstrap() {
	e.operators = evaluator.DefaultOperators
}

func (e *Engine) Evaluate(
	ctx context.Context,
	experiments []auroratype.Experiment,
	parameterName string,
	attr map[string]any,
) *Evaluation {
	sort.Slice(experiments, func(i, j int) bool {
		return experiments[i].Priority < experiments[j].Priority
	})

	for _, exp := range experiments {
		hasParam := false
		for _, p := range exp.Parameters {
			if p == parameterName {
				hasParam = true
				break
			}
		}
		if !hasParam {
			continue
		}

		if !e.checkStatus(exp) {
			continue
		}

		if !e.checkTime(exp) {
			continue
		}

		if !e.checkPopulation(ctx, exp, attr) {
			continue
		}

		if !e.checkConstraints(ctx, exp, attr) {
			continue
		}

		variant := evaluator.SelectVariantByHash(exp.ID, exp.HashAttribute, attr, exp.Variants)
		if variant == nil {
			continue
		}

		return &Evaluation{
			ExperimentID: exp.ID,
			VariantKey:   variant.Key,
			Values:       variant.Values,
			Matched:      true,
		}
	}

	return &Evaluation{
		Matched: false,
	}
}

func (e *Engine) checkStatus(exp auroratype.Experiment) bool {
	return exp.Status == auroratype.StatusRunning
}

func (e *Engine) checkTime(exp auroratype.Experiment) bool {
	now := time.Now().Unix()

	if exp.StartTime != nil && now < *exp.StartTime {
		return false
	}

	if exp.EndTime != nil && now > *exp.EndTime {
		return false
	}

	return true
}

func (e *Engine) checkPopulation(ctx context.Context, exp auroratype.Experiment, attr map[string]any) bool {
	if exp.PopulationSize <= 0 {
		return false
	}

	if exp.PopulationSize >= 100 {
		return true
	}

	if exp.HashAttribute == "" {
		return false
	}

	hashValue := attr[exp.HashAttribute]
	if hashValue == nil {
		return false
	}

	hash := evaluator.CalculateHash(hashValue, exp.ID)
	return evaluator.IsInPercentageRange(hash, exp.PopulationSize)
}

func (e *Engine) checkConstraints(ctx context.Context, exp auroratype.Experiment, attr map[string]any) bool {
	for _, constraint := range exp.Constraints {
		if !evaluator.EvaluateConstraint(constraint, attr, e.operators) {
			return false
		}
	}
	return true
}
