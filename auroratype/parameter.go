package auroratype

type Parameter struct {
	DefaultValue interface{} `yaml:"defaultValue"`
	Rules        []Rule      `yaml:"rules"`
}

type Rule struct {
	RolloutValue  interface{}  `yaml:"rolloutValue"`
	Percentage    *int         `yaml:"percentage,omitempty"`
	HashAttribute *string      `yaml:"hashAttribute,omitempty"`
	EffectiveAt   *int64       `yaml:"effectiveAt,omitempty"`
	Constraints   []Constraint `yaml:"constraints"`
}

type Constraint struct {
	Field    string      `yaml:"field"`
	Operator string      `yaml:"operator"`
	Value    interface{} `yaml:"value"`
}
