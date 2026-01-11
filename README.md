# Aurora-Go

Feature flags and parameter configuration for Go applications. Inspired by AWS AppConfig and LaunchDarkly.

[![Go Reference](https://pkg.go.dev/badge/github.com/tuannguyensn2001/aurora-go.svg)](https://pkg.go.dev/github.com/tuannguyensn2001/aurora-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/tuannguyensn2001/aurora-go)](https://goreportcard.com/report/github.com/tuannguyensn2001/aurora-go)

## Documentation

Full documentation available at https://aurora.tuannguyensn2001aa.workers.dev/

## Installation

```bash
go get github.com/tuannguyensn2001/aurora-go
```

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

## Features

- Feature flags and parameter configuration
- Attribute-based targeting
- Percentage rollouts with consistent hashing
- Multiple fetchers (file, S3)
- Built-in metrics and observability
- Custom operators support
- Strong consistency option

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
