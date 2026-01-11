# Aurora-Go

Feature flags and parameter configuration for Go applications. Inspired by AWS AppConfig and LaunchDarkly.

[![Go Reference](https://pkg.go.dev/badge/github.com/tuannguyensn2001/aurora-go.svg)](https://pkg.go.dev/github.com/tuannguyensn2001/aurora-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/tuannguyensn2001/aurora-go)](https://goreportcard.com/report/github.com/tuannguyensn2001/aurora-go)

## Table of Contents

- [Quick Start](#quick-start)
- [Configuration](#configuration)
- [Usage Examples](#usage-examples)
- [Fetchers](#fetchers)
- [API Reference](#api-reference)
- [Examples](#examples)
- [Deep Dive](#deep-dive)
  - [How It Works](#how-it-works)
  - [Build Your Own Management System](#build-your-own-management-system)
- [Contributing](#contributing)
- [License](#license)

## Quick Start

```go
package main

import (
	"context"
	"log"

	"github.com/tuannguyensn2001/aurora-go"
	"github.com/tuannguyensn2001/aurora-go/fetcher/file"
)

func main() {
	fetcher := file.New("/path/to/parameters.yaml")
	storage := aurora.NewStorage(fetcher)
	client := aurora.NewClient(storage, aurora.ClientOptions{})

	if err := client.Start(context.Background()); err != nil {
		log.Fatal(err)
	}

	attrs := aurora.NewAttribute()
	attrs.Set("country", "US")
	attrs.Set("plan", "premium")

	result := client.GetParameter(context.Background(), "newFeature", attrs)
	if enabled := result.Boolean(false); enabled {
		// Feature enabled for this user
	}
}
```

## Configuration

Create a `parameters.yaml` file:

```yaml
newFeature:
  defaultValue: false
  rules:
    - rolloutValue: true
      percentage: 10
      hashAttribute: userID
      constraints:
        - field: country
          operator: equal
          value: "US"

featureThreshold:
  defaultValue: 100
  rules:
    - rolloutValue: 50
      effectiveAt: 1704067200
      constraints:
        - field: plan
          operator: equal
          value: "free"
```

### Parameter Fields

| Field | Type | Description |
|-------|------|-------------|
| `defaultValue` | any | Fallback when no rules match |
| `rules` | []Rule | Rules evaluated in order |

### Rule Fields

| Field | Type | Description |
|-------|------|-------------|
| `rolloutValue` | any | Value returned when matched |
| `percentage` | number | 0-100 for gradual rollout |
| `hashAttribute` | string | Attribute for consistent hashing |
| `effectiveAt` | number | Unix timestamp for scheduled release |
| `constraints` | []Constraint | Conditions to match |

### Built-in Operators

- `equal` / `notEqual`
- `greaterThan` / `lessThan`
- `greaterThanOrEqual` / `lessThanOrEqual`
- `contains` / `notContains`
- `in` / `notIn`
- `startsWith` / `endsWith`
- `matchesRegex`

## Usage Examples

### Attribute-Based Targeting

```go
attrs := aurora.NewAttribute()
attrs.Set("country", "CA")
attrs.Set("plan", "enterprise")
attrs.Set("age", 25)

result := client.GetParameter(ctx, "featureX", attrs)
value := result.String("default")
```

### Percentage Rollouts

```go
// Enable feature for 25% of users
rule := aurora.Rule{
    RolloutValue: true,
    Percentage:   25,
    HashAttribute: "userID",
}
```

### Custom Operators

```go
client.RegisterOperator("startsWith", func(a, b any) bool {
    return strings.HasPrefix(fmt.Sprint(a), fmt.Sprint(b))
})

// Use in config:
// - field: email
//   operator: startsWith
//   value: "admin@"
```

### Type-Safe Retrieval

```go
result := client.GetParameter(ctx, "myParam", attrs)

boolVal := result.Boolean(false)      // false fallback
intVal := result.Int(0)               // 0 fallback
stringVal := result.String("default") // "default" fallback
```

## Fetchers

### File Fetcher

```go
fetcher := file.New("/path/to/config.yaml")
```

### S3 Fetcher

```go
fetcher := s3.New(ctx, s3.Config{
    Bucket: "my-bucket",
    Key:    "config/parameters.yaml",
    Region: "us-east-1",
})
// Automatically polls for updates
```

## API Reference

### NewClient

```go
client := aurora.NewClient(storage, aurora.ClientOptions{})
```

### Start

```go
if err := client.Start(ctx); err != nil {
    log.Fatal(err)
}
```

Initializes the client. Call before using `GetParameter`.

### GetParameter

```go
result := client.GetParameter(ctx, "parameterName", attrs)
```

Retrieves a parameter value based on matching rules and user attributes.

### RegisterOperator

```go
client.RegisterOperator("myOperator", func(a, b any) bool {
    // your logic
})
```

Registers a custom operator for use in rule constraints.

## Examples

Check out the [examples](/examples) directory:

- [Basic usage](/examples/main.go) - Get started quickly
- [File fetcher](/examples/file_example.go) - Load from local files
- [S3 fetcher](/examples/s3_example.go) - Load from AWS S3

## Deep Dive

### How It Works

Aurora-Go is built around a simple but powerful concept: **evaluate rules to find the right value for each user**.

#### Core Components

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           Aurora-Go SDK                                 │
├─────────────────────────────────────────────────────────────────────────┤
│                                                                         │
│   ┌─────────────┐                                                      │
│   │   Fetcher   │  Reads configuration from source                     │
│   │ (File/S3)   │  Supports YAML and JSON formats                      │
│   └──────┬──────┘                                                      │
│          │                                                             │
│          ▼                                                             │
│   ┌─────────────┐     ┌─────────────────────────────────────────┐     │
│   │   Storage   │ ──▶ │  In-memory cache with automatic        │     │
│   │             │     │  polling and thread-safe access        │     │
│   └──────┬──────┘     └─────────────────────────────────────────┘     │
│          │                                                             │
│          ▼                                                             │
│   ┌─────────────┐     ┌─────────────────────────────────────────┐     │
│   │   Engine    │ ──▶ │  Rule evaluation engine:               │     │
│   │             │     │  - Match constraints against attrs     │     │
│   │             │     │  - Calculate percentage rollouts       │     │
│   │             │     │  - Apply time-based rules              │     │
│   │             │     │  - Support custom operators            │     │
│   └─────────────┘     └─────────────────────────────────────────┘     │
│                                                                         │
└─────────────────────────────────────────────────────────────────────────┘
```

#### Evaluation Flow

When you call `GetParameter(name, attributes)`:

```
1. Check in-memory cache for parameter config
2. If not cached or stale, fetch from source (file/S3)
3. Find the parameter by name
4. Evaluate rules in order (first matching rule wins):
   a. Check all constraints match user attributes
   b. If percentage rollout: hash(userAttribute) % 100 < percentage
   c. If effectiveAt: check if current time >= scheduled time
5. Return rolloutValue if rule matches, otherwise try next rule
6. Return defaultValue if no rules match
```

#### Rule Evaluation Order

Rules are evaluated **sequentially** - the first matching rule wins:

```yaml
featureX:
  defaultValue: false
  rules:
    - # Rule 1: Checked first (highest priority)
      rolloutValue: true
      constraints:
        - field: plan
          operator: equal
          value: "premium"

    - # Rule 2: Checked second
      rolloutValue: true
      percentage: 50
      hashAttribute: userID
      constraints:
        - field: country
          operator: equal
          value: "US"

    - # Rule 3: Checked last (fallback for matched constraints)
      rolloutValue: false
```

#### Percentage Rollouts

Uses **Murmur3 hashing** for consistent, deterministic rollouts:

```go
// Same user always gets the same result (no flickering)
hash := murmur3.Sum64([]string(userID, "featureX"))
bucket := hash % 100  // 0-99

// If bucket < percentage, user gets the rollout value
if bucket < 25 {  // 25% of users
    // Feature enabled
}
```

**Benefits:**
- Same user always gets same result (consistent experience)
- Even distribution across users
- 0.01% granularity (1 in 10,000 users)
- No server-side state needed

#### Time-Based Rules

Rules can be scheduled for future activation:

```yaml
featureX:
  defaultValue: false
  rules:
    - rolloutValue: true
      effectiveAt: 1704067200  # Jan 1, 2024 00:00:00 UTC
      constraints:
        - field: plan
          operator: equal
          value: "beta"
```

The rule only activates when `current_time >= effectiveAt`. This enables:
- Scheduled feature releases
- Marketing campaign timing
- Compliance with launch dates
- Gradual expansion to more users

---

### Build Your Own Management System

Aurora-Go is designed as a lightweight SDK that handles rule evaluation. You can build a complete feature management platform on top of it:

```
┌─────────────────────────────────────────────────────────────────────────────────────┐
│                         Your Feature Flag Platform                                   │
├─────────────────────────────────────────────────────────────────────────────────────┤
│                                                                                      │
│  ┌─────────────┐     ┌─────────────┐     ┌─────────────────────────────┐           │
│  │   Web UI    │     │   Admin     │     │   REST API                  │           │
│  │  Dashboard  │     │    API      │     │   (CRUD Operations)         │           │
│  └──────┬──────┘     └──────┬──────┘     └─────────────┬───────────────┘           │
│         │                   │                           │                            │
│         └───────────────────┼───────────────────────────┘                            │
│                             ▼                                                        │
│                   ┌─────────────────┐                                                │
│                   │   Database      │  (PostgreSQL, MySQL, DynamoDB, etc.)          │
│                   │  (Authoritative │                                                │
│                   │   Source of     │                                                │
│                   │   Truth)        │                                                │
│                   └────────┬────────┘                                                │
│                            │                                                         │
│                            ▼                                                         │
│                   ┌─────────────────┐     ┌───────────────────────────────────┐     │
│                   │ Config Generator│     │  Audit Logger                     │     │
│                   │ (Export to YAML)│     │  (Track all changes)              │     │
│                   └────────┬────────┘     └───────────────────────────────────┘     │
│                            │                                                          │
│                            ▼                                                          │
│                   ┌─────────────────┐                                                │
│                   │  AWS S3 / GCS   │  (Recommended for production)                │
│                   │  (Config Store) │                                                │
│                   └────────┬────────┘                                                │
│                            │                                                         │
│                            │    ┌───────────────────────────────────────────┐        │
│                            │    │  CDN (Optional)                           │        │
│                            │    │  CloudFront, Cloudflare, etc.            │        │
│                            │    └─────────────────┬─────────────────────────┘        │
│                            │                      │                                 │
│                            └──────────────────────┼───────────────────────┐         │
│                                                   │                       ▼         │
│                                   ┌───────────────┴───────────────┐    ┌─────────┐ │
│                                   │      Aurora-Go SDK            │    │Your App │ │
│                                   │  (Rule Evaluation + Caching)  │    │         │ │
│                                   └───────────────────────────────┘    └─────────┘ │
│                                                                                      │
└─────────────────────────────────────────────────────────────────────────────────────┘
```

#### Why Build Your Own?

Aurora-Go focuses on what it does best: **fast, reliable rule evaluation**. By building your own management system, you get:

| What You Control | Why It Matters |
|------------------|----------------|
| **Management UI** | User-friendly interface for your team |
| **API Endpoints** | Integrate with CI/CD, chatbots, scripts |
| **Database** | Store history, audit logs, user segments |
| **Authentication** | Role-based access control |
| **Webhooks** | Notify external services on changes |
| **Analytics** | Track flag usage and experiment results |

#### Integration Options

##### Option 1: AWS S3 (Recommended for Production)

Upload YAML configs to S3. The SDK polls S3 for updates automatically.

```go
import "github.com/tuannguyensn2001/aurora-go/fetcher/s3"

fetcher := s3.New(ctx, s3.Config{
    Bucket: "my-feature-flags",
    Key:    "production/parameters.yaml",
    Region: "us-east-1",
})

storage := aurora.NewStorage(fetcher)
client := aurora.NewClient(storage, aurora.ClientOptions{})
```

**Benefits:**
- Automatic polling for hot-reload (no restarts needed)
- Scalable and highly available (S3 guarantees 99.99% availability)
- Versioning support (revert mistakes instantly)
- Works across multiple instances/regions
- Low latency with S3 + CloudFront

##### Option 2: Static Files

For simpler setups, testing, or when you control deployment:

```go
import "github.com/tuannguyensn2001/aurora-go/fetcher/file"

fetcher := file.New("/path/to/config.yaml")

storage := aurora.NewStorage(fetcher)
client := aurora.NewClient(storage, aurora.ClientOptions{})
```

**Use cases:**
- Local development
- Testing environments
- Air-gapped systems
- Simple applications

#### Scalability Strategies

| Component | Strategy | Implementation |
|-----------|----------|----------------|
| **Config Storage** | Use S3 | Unlimited scalability, 99.99% availability |
| **SDK Cache** | In-memory | Built-in with automatic polling |
| **Multiple Instances** | Share config source | All instances use same S3 bucket |
| **Global Distribution** | S3 + CDN | CloudFront for low-latency access |
| **High Traffic** | Local caching | SDK caches in-memory, reduces S3 calls |

**Scaling Architecture:**

```
                              ┌─────────────────────────────────────┐
                              │            AWS S3                   │
                              │  (Single Source of Truth)           │
                              └──────────────┬──────────────────────┘
                                             │
                              ┌──────────────┼──────────────┐
                              │              │              │
                              ▼              ▼              ▼
                       ┌───────────┐  ┌───────────┐  ┌───────────┐
                       │CloudFront │  │CloudFront │  │CloudFront │
                       │  (Edge)   │  │  (Edge)   │  │  (Edge)   │
                       └─────┬─────┘  └─────┬─────┘  └─────┬─────┘
                             │              │              │
                             └──────────────┼──────────────┘
                                            │
                       ┌────────────────────┼────────────────────┐
                       │                    │                    │
                       ▼                    ▼                    ▼
                ┌────────────┐      ┌────────────┐      ┌────────────┐
                │   App 1    │      │   App 2    │      │   App N    │
                │ (us-east)  │      │ (us-west)  │      │ (eu-west)  │
                │            │      │            │      │            │
                │ In-memory  │      │ In-memory  │      │ In-memory  │
                │   Cache    │      │   Cache    │      │   Cache    │
                └────────────┘      └────────────┘      └────────────┘
```

**Best Practices:**

1. **Use S3 as source of truth** - Never store configs in database
2. **Enable S3 versioning** - Instant rollback when mistakes happen
3. **Use CloudFront CDN** - Reduce latency for global applications
4. **Cache aggressively** - SDK caches in-memory, CDN caches at edge
5. **Separate environments** - Different S3 paths for dev/staging/prod

#### Reliability Strategies

| Risk | Mitigation |
|------|------------|
| **S3 outage** | SDK uses cached config; graceful degradation |
| **Config corruption** | Validate YAML before upload; use versioned uploads |
| **Network latency** | In-memory cache + CDN reduces fetch frequency |
| **Rollout errors** | Test in staging first; use gradual percentage rollouts |
| **Permission issues** | Use IAM roles for EC2/ECS; assume role for cross-account |

**Resilience Pattern:**

```go
// Handle config fetch failures gracefully
fetcher := s3.New(ctx, s3.Config{
    Bucket: "my-feature-flags",
    Key:    "production/parameters.yaml",
    Region: "us-east-1",
})

storage := aurora.NewStorage(fetcher)
client := aurora.NewClient(storage, aurora.ClientOptions{})

if err := client.Start(ctx); err != nil {
    // Log error but don't crash
    // SDK will use empty config (all defaultValues)
    log.Printf("Warning: Failed to load feature flags: %v", err)
}

// Default values are used when:
// - Config fetch fails
- Parameter not found
- No rules match
```

#### Example: Generate and Upload Config

```go
package main

import (
	"bytes"
	"context"
	"fmt"

	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v2"
)

type ConfigGenerator struct {
	s3Client *s3.Client
	bucket   string
}

func NewConfigGenerator(bucket, region string) (*ConfigGenerator, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, fmt.Errorf("load aws config: %w", err)
	}

	return &ConfigGenerator{
		s3Client: s3.NewFromConfig(cfg),
		bucket:   bucket,
	}, nil
}

func (g *ConfigGenerator) CreateFeatureFlag(
	ctx context.Context,
	name string,
	rules []auroratype.Rule,
	defaultValue any,
) error {
	params := map[string]auroratype.Parameter{
		name: {
			DefaultValue: defaultValue,
			Rules:        rules,
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(params)
	if err != nil {
		return fmt.Errorf("marshal yaml: %w", err)
	}

	// Upload to S3
	_, err = g.s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:               aws.String(g.bucket),
		Key:                  aws.String(fmt.Sprintf("features/%s.yaml", name)),
		Body:                 bytes.NewReader(data),
		ContentType:          aws.String("application/yaml"),
		ServerSideEncryption: aws.String("AES256"),
	})
	if err != nil {
		return fmt.Errorf("upload to s3: %w", err)
	}

	return nil
}

// Usage
func main() {
	gen, _ := NewConfigGenerator("my-feature-flags", "us-east-1")

	rules := []auroratype.Rule{
		{
			RolloutValue: true,
			Constraints: []auroratype.Constraint{
				{
					Field:    "plan",
					Operator: auroratype.OperatorEqual,
					Value:    "premium",
				},
			},
		},
	}

	gen.CreateFeatureFlag(context.Background(), "newCheckoutFlow", rules, false)
}
```

#### Example: Management API with Go + Gin

```go
package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

type CreateFlagRequest struct {
	Name          string                  `json:"name" binding:"required"`
	DefaultValue  any                     `json:"defaultValue"`
	Rules         []auroratype.Rule       `json:"rules"`
}

type ManagementAPI struct {
	generator *ConfigGenerator
}

func (api *ManagementAPI) CreateFlag(c *gin.Context) {
	var req CreateFlagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate rules
	for i, rule := range req.Rules {
		if err := validateRule(rule); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": fmt.Sprintf("rule %d: %v", i, err),
			})
			return
		}
	}

	// Create feature flag
	err := api.generator.CreateFeatureFlag(
		c.Request.Context(),
		req.Name,
		req.Rules,
		req.DefaultValue,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "flag created"})
}

func (api *ManagementAPI) ListFlags(c *gin.Context) {
	// List all flags from S3
	c.JSON(http.StatusOK, gin.H{"flags": []string{}})
}

func validateRule(rule auroratype.Rule) error {
	if rule.Percentage > 0 && rule.HashAttribute == "" {
		return ErrHashAttributeRequired
	}
	return nil
}

func main() {
	r := gin.Default()

	api := &ManagementAPI{
		generator: NewConfigGenerator("my-feature-flags", "us-east-1"),
	}

	r.POST("/api/flags", api.CreateFlag)
	r.GET("/api/flags", api.ListFlags)

	r.Run(":8080")
}
```

#### Monitoring and Observability

Track the health of your feature flag system:

```go
// Custom metrics for your management system
type Metrics struct {
	flagsCreated    prometheus.Counter
	flagsUpdated    prometheus.Counter
	configFetchLatency prometheus.Histogram
}

func (m *Metrics) RecordConfigFetch(durationMs float64) {
	m.configFetchLatency.Observe(durationMs)
}

// Aurora-Go SDK logging
client := aurora.NewClient(storage, aurora.ClientOptions{
	Logger: slog.New(slog.NewJSONHandler(os.Stdout, nil)),
})
```

**Key Metrics to Monitor:**

| Metric | Description | Alert Threshold |
|--------|-------------|-----------------|
| Config fetch latency | Time to fetch from S3 | > 1s |
| Cache hit rate | Percentage of requests served from cache | < 95% |
| Flag evaluation errors | Errors during rule evaluation | > 0 |
| SDK initialization time | Time to start client | > 5s |

#### Security Considerations

| Concern | Mitigation |
|---------|------------|
| **Unauthorized flag changes** | Use IAM roles with least privilege |
| **Config access from app** | Restrict S3 bucket access to specific roles |
| **Audit trail** | Log all flag changes with user and timestamp |
| **Sensitive data** | Don't put secrets in feature flags |
| **YAML injection** | Validate and sanitize all inputs |

**IAM Policy Example:**

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": [
        "s3:GetObject",
        "s3:PutObject"
      ],
      "Resource": "arn:aws:s3:::my-feature-flags/*"
    }
  ]
}
```

#### Recommended Workflow

```
1. Create flag in UI → Saves to database
2. Review changes in UI → Approval workflow (optional)
3. Deploy to staging → Test with real users
4. Roll out to production → Start with 1%, gradually increase
5. Monitor metrics → Watch for issues
6. Full rollout → 100% of users
7. Cleanup → Remove old flags, document learnings
```

This architecture gives you full control over your feature management workflow while keeping the Aurora-Go SDK simple, fast, and focused on rule evaluation.

## Best Practices

### Naming Conventions

Use clear, descriptive names for parameters and attributes:

```yaml
# Good examples
newCheckoutFlow:
darkModeEnabled:
maxUploadSizeMB:
apiRateLimit:

# Avoid
feature1:
flagX:
test:
```

**Parameter naming tips:**
- Use kebab-case for parameter names (`new-feature` not `newFeature`)
- Use descriptive prefixes: `enable-*`, `disable-*`, `max-*`, `min-*`
- Include units in names when applicable (`timeoutSeconds` not `timeout`)

**Attribute naming tips:**
- Use camelCase for attributes (`userId` not `userID`)
- Use consistent prefixes: `user.*`, `device.*`, `session.*`

### Configuration Structure

Organize your configuration files logically:

```yaml
# Group related flags together
feature_flags:
  newCheckoutFlow:
    defaultValue: false
    rules:
      - rolloutValue: true
        constraints:
          - field: plan
            operator: equal
            value: "premium"

experiments:
  pricingTest:
    defaultValue: "original"
    rules:
      - rolloutValue: "variantA"
        percentage: 50
        hashAttribute: userId

runtime_config:
  maxUploadSizeMB:
    defaultValue: 100
  sessionTimeoutMinutes:
    defaultValue: 30
```

### Rule Design

Keep rules simple and focused:

```yaml
# Good: One clear condition per rule
featureX:
  defaultValue: false
  rules:
    - rolloutValue: true
      constraints:
        - field: plan
          operator: equal
          value: "premium"

# Avoid: Overly complex rules
featureX:
  defaultValue: false
  rules:
    - rolloutValue: true
      constraints:
        - field: plan
          operator: equal
          value: "premium"
        - field: country
          operator: in
          value: ["US", "CA", "UK"]
        - field: age
          operator: greaterThan
          value: 18
        - field: device
          operator: in
          value: ["ios", "android"]
```

**Tips for rule design:**
- Order rules from most specific to most general
- Use percentage rollouts for gradual rollout
- Always have a default value
- Avoid combining too many constraints

### Gradual Rollouts

Always start with small percentages and increase gradually:

```yaml
newFeature:
  defaultValue: false
  rules:
    # Phase 1: Internal team (100% for specific email domain)
    - rolloutValue: true
      constraints:
        - field: email
          operator: matchesRegex
          value: "@company\\.com$"

    # Phase 2: 1% of users
    - rolloutValue: true
      percentage: 1
      hashAttribute: userId

    # Phase 3: 10% of users
    - rolloutValue: true
      percentage: 10
      hashAttribute: userId

    # Phase 4: 50% of users
    - rolloutValue: true
      percentage: 50
      hashAttribute: userId
```

### Default Values

Always define explicit default values:

```yaml
# Good: Explicit defaults
featureX:
  defaultValue: false  # Clear boolean default
  rules:
    - rolloutValue: true
      percentage: 10
      hashAttribute: userId

# Avoid: Implicit defaults
featureX:
  # No defaultValue - can cause type issues
  rules:
    - rolloutValue: true
```

**Default value types:**
- Boolean features: `defaultValue: false`
- Numeric values: `defaultValue: 100` (not `0` for limit values)
- Strings: `defaultValue: "standard"` (not empty string)
- Arrays: `defaultValue: []` (empty array, not null)

### Testing Flags

Test your flags before rolling out:

```go
// Test in your test suite
func TestFeatureFlag(t *testing.T) {
    client := setupTestClient()

    tests := []struct {
        name      string
        attribute map[string]any
        expected  bool
    }{
        {
            name:      "premium user gets feature",
            attribute: map[string]any{"plan": "premium"},
            expected:  true,
        },
        {
            name:      "free user does not get feature",
            attribute: map[string]any{"plan": "free"},
            expected:  false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            attrs := aurora.NewAttribute()
            for k, v := range tt.attribute {
                attrs.Set(k, v)
            }
            result := client.GetParameter(context.Background(), "featureX", attrs)
            if result.Boolean(false) != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result.Boolean(false))
            }
        })
    }
}
```

### Error Handling

Handle SDK initialization gracefully:

```go
func main() {
    fetcher := s3.New(ctx, s3.Config{
        Bucket: "my-feature-flags",
        Key:    "production/parameters.yaml",
        Region: "us-east-1",
    })

    storage := aurora.NewStorage(fetcher)
    client := aurora.NewClient(storage, aurora.ClientOptions{})

    // Try to start, but don't crash if it fails
    if err := client.Start(ctx); err != nil {
        // Log the error
        log.Printf("Warning: Failed to initialize feature flags: %v", err)
        // Application continues with default values
    }

    // Use flag with fallback
    result := client.GetParameter(ctx, "newFeature", attrs)
    enabled := result.Boolean(false) // Always safe fallback
    if enabled {
        // Feature logic
    }
}
```

### Performance

Optimize for performance in high-traffic applications:

```go
// 1. Create client once at startup, reuse it
var client *aurora.Client

func init() {
    fetcher := s3.New(ctx, s3.Config{
        Bucket: "my-feature-flags",
        Key:    "production/parameters.yaml",
        Region: "us-east-1",
    })
    storage := aurora.NewStorage(fetcher)
    client = aurora.NewClient(storage, aurora.ClientOptions{})
}

// 2. Reuse attribute objects
var attrs *aurora.Attribute

func handleRequest(user User) {
    attrs = aurora.NewAttribute()
    attrs.Set("userId", user.ID)
    attrs.Set("plan", user.Plan)
    attrs.Set("country", user.Country)

    result := client.GetParameter(ctx, "featureX", attrs)
    // Use result
}

// 3. Cache results if appropriate
var flagCache = struct {
    sync.RWMutex
    values map[string]bool
}{values: make(map[string]bool)}

func getCachedFlag(flagName string, attrs *aurora.Attribute) bool {
    cacheKey := flagName + ":" + attrs.Get("userId")
    flagCache.RLock()
    if enabled, ok := flagCache.values[cacheKey]; ok {
        flagCache.RUnlock()
        return enabled
    }
    flagCache.RUnlock()

    result := client.GetParameter(ctx, flagName, attrs)
    enabled := result.Boolean(false)

    flagCache.Lock()
    flagCache.values[cacheKey] = enabled
    flagCache.Unlock()

    return enabled
}
```

### Organization

Structure your feature flags across environments:

```
s3://my-feature-flags/
├── development/
│   └── parameters.yaml
├── staging/
│   └── parameters.yaml
└── production/
    └── parameters.yaml
```

**Environment separation tips:**
- Use different S3 paths for each environment
- Never share production flags with development
- Use same flag names across environments for consistency
- Test in staging before production

### Documentation

Document your feature flags:

```yaml
# This flag controls the new checkout flow experience
# Created by: @john
# Date: 2024-01-15
# Rollout: Start with 1% on 2024-01-20, increase to 10% on 2024-01-25
# Related PR: https://github.com/company/repo/pull/123
newCheckoutFlow:
  defaultValue: false
  rules:
    - rolloutValue: true
      percentage: 1
      hashAttribute: userId
```

### Cleanup

Remove old flags to prevent confusion:

```go
// Cleanup script - run periodically
func CleanupOldFlags(bucket, prefix string) error {
    client := s3.NewFromConfig(cfg)

    paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
        Bucket: aws.String(bucket),
        Prefix: aws.String(prefix),
    })

    for paginator.HasMorePages() {
        page, err := paginator.NextPage(ctx)
        if err != nil {
            return err
        }

        for _, obj := range page.Contents {
            // Check if flag is older than 90 days
            if time.Since(*obj.LastModified) > 90*24*time.Hour {
                // Check if flag is still referenced
                if !isFlagInUse(*obj.Key) {
                    // Archive and delete
                    archiveFlag(*obj.Key)
                    client.DeleteObject(ctx, &s3.DeleteObjectInput{
                        Bucket: aws.String(bucket),
                        Key:    obj.Key,
                    })
                }
            }
        }
    }
    return nil
}
```

### Monitoring

Monitor your feature flags in production:

```go
// Add metrics to your application
import "github.com/prometheus/client_golang/prometheus"

var (
    flagEvaluations = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "feature_flag_evaluations_total",
            Help: "Total number of feature flag evaluations",
        },
        []string{"flag_name", "result"},
    )

    configFetchDuration = prometheus.NewHistogram(
        prometheus.HistogramOpts{
            Name:    "config_fetch_duration_seconds",
            Help:    "Time to fetch configuration",
            Buckets: []float64{0.01, 0.05, 0.1, 0.5, 1},
        },
    )
)

func init() {
    prometheus.MustRegister(flagEvaluations)
    prometheus.MustRegister(configFetchDuration)
}

func evaluateFlag(flagName string, attrs *aurora.Attribute) bool {
    start := time.Now()
    result := client.GetParameter(ctx, flagName, attrs)
    duration := time.Since(start).Seconds()

    configFetchDuration.Observe(duration)

    enabled := result.Boolean(false)
    flagEvaluations.WithLabelValues(flagName, strconv.FormatBool(enabled)).Inc()

    return enabled
}
```

**Key metrics to track:**
- Flag evaluation rate
- Config fetch latency
- Percentage of flags enabled
- Errors and failures

### Common Pitfalls

Avoid these common mistakes:

| Pitfall | Solution |
|---------|----------|
| Forgetting `hashAttribute` in percentage rollout | Always include `hashAttribute` when using `percentage` |
| Overlapping rules | Design rules to be mutually exclusive |
| Hardcoding values in rules | Use attributes instead of hardcoded values |
| No default value | Always specify `defaultValue` |
| Testing in production only | Test in staging first |
| Not versioning configs | Enable S3 versioning |
| Ignoring SDK errors | Log and handle errors gracefully |
| Too many constraints per rule | Keep rules simple and focused |

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
