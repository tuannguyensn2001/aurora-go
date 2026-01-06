package s3

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

// Fetcher fetches configuration from an S3 bucket.
type Fetcher struct {
	client *s3.Client
	bucket string
	key    string
}

// Options configures the S3 Fetcher.
type Options struct {
	Client *s3.Client
	Bucket string
	Key    string
}

// NewFetcher creates a new S3-based Fetcher.
func NewFetcher(opts Options) *Fetcher {
	return &Fetcher{
		client: opts.Client,
		bucket: opts.Bucket,
		key:    opts.Key,
	}
}

// Fetch retrieves configuration data from S3.
func (f *Fetcher) Fetch(ctx context.Context) (map[string]auroratype.Parameter, error) {
	output, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(f.key),
	})
	if err != nil {
		return nil, err
	}
	defer output.Body.Close()

	data, err := io.ReadAll(output.Body)
	if err != nil {
		return nil, err
	}

	var config map[string]auroratype.Parameter
	if strings.HasSuffix(f.key, ".yaml") || strings.HasSuffix(f.key, ".yml") {
		err = yaml.Unmarshal(data, &config)
	} else {
		err = json.Unmarshal(data, &config)
	}

	if err != nil {
		return nil, err
	}

	return config, nil
}

func (f *Fetcher) IsStatic() bool {
	return false
}
