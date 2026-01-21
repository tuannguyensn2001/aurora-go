package experiment

import (
	"context"
	"testing"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
)

func TestEngine_Evaluate(t *testing.T) {
	engine := NewEngine()
	engine.Bootstrap()

	experiments := []auroratype.Experiment{
		{
			ID:             "exp_001",
			Name:           "Test Experiment",
			Parameters:     []string{"buttonColor"},
			HashAttribute:  "userID",
			PopulationSize: 100,
			Priority:       1,
			Status:         auroratype.StatusRunning,
			Variants: []auroratype.Variant{
				{
					Key:     "control",
					Rollout: 50,
					Values: map[string]interface{}{
						"buttonColor": "blue",
					},
				},
				{
					Key:     "treatment",
					Rollout: 50,
					Values: map[string]interface{}{
						"buttonColor": "green",
					},
				},
			},
		},
	}

	ctx := context.Background()
	attr := map[string]any{
		"userID": "user123",
	}

	result := engine.Evaluate(ctx, experiments, "buttonColor", attr)

	if !result.Matched {
		t.Error("Expected experiment to match")
	}

	if result.ExperimentID != "exp_001" {
		t.Errorf("Expected experiment ID exp_001, got %s", result.ExperimentID)
	}

	if result.Values["buttonColor"] == nil {
		t.Error("Expected buttonColor value to be set")
	}
}

func TestEngine_NoExperimentForParameter(t *testing.T) {
	engine := NewEngine()
	engine.Bootstrap()

	experiments := []auroratype.Experiment{
		{
			ID:         "exp_001",
			Name:       "Test Experiment",
			Parameters: []string{"otherParam"},
			Priority:   1,
			Status:     auroratype.StatusRunning,
			Variants: []auroratype.Variant{
				{
					Key:     "control",
					Rollout: 100,
					Values: map[string]interface{}{
						"otherParam": "value",
					},
				},
			},
		},
	}

	ctx := context.Background()
	attr := map[string]any{
		"userID": "user123",
	}

	result := engine.Evaluate(ctx, experiments, "buttonColor", attr)

	if result.Matched {
		t.Error("Expected experiment not to match")
	}
}

func TestEngine_PopulationSize(t *testing.T) {
	engine := NewEngine()
	engine.Bootstrap()

	experiments := []auroratype.Experiment{
		{
			ID:             "exp_001",
			Name:           "Test Experiment",
			Parameters:     []string{"buttonColor"},
			HashAttribute:  "userID",
			PopulationSize: 0,
			Priority:       1,
			Status:         auroratype.StatusRunning,
			Variants: []auroratype.Variant{
				{
					Key:     "control",
					Rollout: 100,
					Values: map[string]interface{}{
						"buttonColor": "blue",
					},
				},
			},
		},
	}

	ctx := context.Background()
	attr := map[string]any{
		"userID": "user123",
	}

	result := engine.Evaluate(ctx, experiments, "buttonColor", attr)

	if result.Matched {
		t.Error("Expected experiment not to match with populationSize 0")
	}
}

func TestEngine_StatusNotRunning(t *testing.T) {
	engine := NewEngine()
	engine.Bootstrap()

	experiments := []auroratype.Experiment{
		{
			ID:         "exp_001",
			Name:       "Test Experiment",
			Parameters: []string{"buttonColor"},
			Priority:   1,
			Status:     auroratype.StatusScheduled,
			Variants: []auroratype.Variant{
				{
					Key:     "control",
					Rollout: 100,
					Values: map[string]interface{}{
						"buttonColor": "blue",
					},
				},
			},
		},
	}

	ctx := context.Background()
	attr := map[string]any{
		"userID": "user123",
	}

	result := engine.Evaluate(ctx, experiments, "buttonColor", attr)

	if result.Matched {
		t.Error("Expected experiment not to match when status is not running")
	}
}
