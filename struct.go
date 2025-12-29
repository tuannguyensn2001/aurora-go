package aurora

type parameter struct {
	DefaultValue interface{} `yaml:"defaultValue"`
	Rules        []rule      `yaml:"rules"`
}

type rule struct {
	RolloutValue  interface{}  `yaml:"rolloutValue"`
	Percentage    *int         `yaml:"percentage,omitempty"`
	HashAttribute *string      `yaml:"hashAttribute,omitempty"`
	EffectiveAt   *int64       `yaml:"effectiveAt,omitempty"`
	Constraints   []constraint `yaml:"constraints"`
}

type constraint struct {
	Field    string      `yaml:"field"`
	Operator string      `yaml:"operator"`
	Value    interface{} `yaml:"value"`
}
