package auroratype

type ExperimentStatus string

const (
	StatusScheduled ExperimentStatus = "scheduled"
	StatusRunning   ExperimentStatus = "running"
	StatusAborted   ExperimentStatus = "aborted"
	StatusFinished  ExperimentStatus = "finished"
)

type Variant struct {
	Key     string                 `yaml:"key"`
	Rollout int                    `yaml:"rollout"`
	Values  map[string]interface{} `yaml:"values"`
}

type Experiment struct {
	ID             string           `yaml:"id"`
	Name           string           `yaml:"name"`
	Parameters     []string         `yaml:"parameters"`
	HashAttribute  string           `yaml:"hashAttribute"`
	PopulationSize int              `yaml:"populationSize"`
	Priority       int              `yaml:"priority"`
	Status         ExperimentStatus `yaml:"status"`
	StartTime      *int64           `yaml:"startTime,omitempty"`
	EndTime        *int64           `yaml:"endTime,omitempty"`
	Constraints    []Constraint     `yaml:"constraints"`
	Variants       []Variant        `yaml:"variants"`
}
