package core

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/spaolacci/murmur3"
)

type engine struct {
	mu        sync.Mutex
	operators map[Operator]func(a, b any) bool
}

func (e *engine) registerOperator(name Operator, fn func(a, b any) bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
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

func (e *engine) evaluateParameter(ctx context.Context, parameterName string, parameter Parameter, attribute *attribute) *resolvedValue {
	for _, rule := range parameter.Rules {
		match := e.evaluateRule(ctx, parameterName, rule, attribute)
		if match {
			return NewResolvedValue(rule.RolloutValue, match)
		}
	}
	return NewResolvedValue(parameter.DefaultValue, false)
}

func (e *engine) evaluateRule(ctx context.Context, parameterName string, rule Rule, attribute *attribute) bool {
	// Check if rule is effective (time-based check)
	if rule.EffectiveAt != nil {
		currentTime := time.Now().Unix()
		if currentTime < *rule.EffectiveAt {
			// Rule is not yet effective
			return false
		}
	}

	// First, check all constraints
	for _, constraint := range rule.Constraints {
		operator := e.operators[Operator(constraint.Operator)]
		if operator == nil {
			return false
		}
		if !operator(attribute.Get(constraint.Field), constraint.Value) {
			return false
		}
	}

	// If constraints pass, check percentage rollout if configured
	if rule.Percentage != nil && rule.HashAttribute != nil {
		hashValue := attribute.Get(*rule.HashAttribute)
		if hashValue == nil {
			return false
		}

		// Calculate hash percentage
		percentage := calculateHashPercentage(parameterName, hashValue, *rule.Percentage)
		if !percentage {
			return false
		}
	}

	return true
}

// calculateHashPercentage calculates if the hash of the value falls within the percentage range
func calculateHashPercentage(parameterName string, value interface{}, percentage int) bool {
	if percentage <= 0 {
		return false
	}
	if percentage >= 100 {
		return true
	}

	// Convert value to string for hashing
	var valueStr string
	switch v := value.(type) {
	case string:
		valueStr = v
	case int, int8, int16, int32, int64:
		valueStr = fmt.Sprintf("%d", v)
	case uint, uint8, uint16, uint32, uint64:
		valueStr = fmt.Sprintf("%d", v)
	case float32, float64:
		valueStr = fmt.Sprintf("%f", v)
	default:
		valueStr = fmt.Sprintf("%v", v)
	}

	// Combine parameterName and hashAttributeValue for hashing
	// This ensures different parameters have independent rollouts
	combinedStr := parameterName + ":" + valueStr

	// Hash the combined value using Murmur3
	hash := murmur3.Sum32([]byte(combinedStr))

	// Use 10000 buckets for better precision (0.01% granularity)
	const numBuckets = 10000
	hashBucket := int(hash % numBuckets)

	// Convert percentage (0-100) to bucket threshold (0-10000)
	threshold := percentage * (numBuckets / 100)

	// Check if hash bucket is less than the threshold
	return hashBucket < threshold
}
