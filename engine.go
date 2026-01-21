package core

import (
	"context"
	"sync"
	"time"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/tuannguyensn2001/aurora-go/core/evaluator"
)

type engine struct {
	mu        sync.Mutex
	operators map[Operator]func(a, b any) bool
}

func (e *engine) registerOperator(name Operator, fn func(a, b any) bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	if fn == nil {
		return
	}
	e.operators[name] = fn
}

func newEngine() *engine {
	return &engine{
		operators: make(map[Operator]func(a, b any) bool),
	}
}

func (e *engine) bootstrap() {
	e.registerOperator(Equal, equalOperator)
	e.registerOperator(NotEqual, notEqualOperator)
	e.registerOperator(GreaterThan, greaterThanOperator)
	e.registerOperator(LessThan, lessThanOperator)
	e.registerOperator(GreaterThanOrEqual, greaterThanOrEqualOperator)
	e.registerOperator(LessThanOrEqual, lessThanOrEqualOperator)
	e.registerOperator(Contains, containsOperator)
	e.registerOperator(In, inOperator)
	e.registerOperator(NotIn, notInOperator)
}

func (e *engine) evaluateParameter(ctx context.Context, parameterName string, parameter auroratype.Parameter, attribute *attribute) *resolvedValue {
	for _, rule := range parameter.Rules {
		match := e.evaluateRule(ctx, parameterName, rule, attribute)
		if match {
			return NewResolvedValue(rule.RolloutValue, match)
		}
	}
	return NewResolvedValue(parameter.DefaultValue, false)
}

func (e *engine) evaluateRule(ctx context.Context, parameterName string, rule auroratype.Rule, attribute *attribute) bool {
	if rule.EffectiveAt != nil {
		currentTime := time.Now().Unix()
		if currentTime < *rule.EffectiveAt {
			return false
		}
	}

	for _, constraint := range rule.Constraints {
		op := e.operators[Operator(constraint.Operator)]
		if op == nil {
			return false
		}
		if !op(attribute.Get(constraint.Field), constraint.Value) {
			return false
		}
	}

	if rule.Percentage != nil && rule.HashAttribute != nil {
		hashValue := attribute.Get(*rule.HashAttribute)
		if hashValue == nil {
			return false
		}

		hash := evaluator.CalculateHash(hashValue, parameterName)
		if !evaluator.IsInPercentageRange(hash, *rule.Percentage) {
			return false
		}
	}

	return true
}
