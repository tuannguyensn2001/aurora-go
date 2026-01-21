package s3

import (
	"context"
	"encoding/json"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/tuannguyensn2001/aurora-go/auroratype"
	"gopkg.in/yaml.v2"
)

const (
	MetricS3FetchLatency            = "s3_fetch_latency"
	MetricS3FetchTotal              = "s3_fetch_total"
	MetricS3FetchExperimentsLatency = "s3_fetch_experiments_latency"
	MetricS3FetchExperimentsTotal   = "s3_fetch_experiments_total"
)

type MetricsRecorder interface {
	Count(metricName string, count int, tags []string)
	Histogram(metricName string, value float64, tags []string)
}

// Fetcher fetches configuration from an S3 bucket.
type Fetcher struct {
	client         *s3.Client
	bucket         string
	key            string
	experimentsKey string
	recorder       MetricsRecorder
}

// Options configures the S3 Fetcher.
type Options struct {
	Client          *s3.Client
	Bucket          string
	Key             string
	ExperimentsKey  string
	MetricsRecorder MetricsRecorder
}

// NewFetcher creates a new S3-based Fetcher.
func NewFetcher(opts Options) *Fetcher {
	recorder := opts.MetricsRecorder
	if recorder == nil {
		recorder = &noopRecorder{}
	}
	return &Fetcher{
		client:         opts.Client,
		bucket:         opts.Bucket,
		key:            opts.Key,
		experimentsKey: opts.ExperimentsKey,
		recorder:       recorder,
	}
}

type noopRecorder struct{}

func (n *noopRecorder) Count(metricName string, count int, tags []string)         {}
func (n *noopRecorder) Histogram(metricName string, value float64, tags []string) {}

// Fetch retrieves configuration data from S3.
func (f *Fetcher) Fetch(ctx context.Context) (map[string]auroratype.Parameter, error) {
	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Microseconds())
		f.recorder.Histogram(MetricS3FetchLatency, duration, []string{"unit:microseconds"})
	}()

	output, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(f.key),
	})
	if err != nil {
		f.recorder.Count(MetricS3FetchTotal, 1, []string{"status:error"})
		return nil, err
	}

	data, err := io.ReadAll(output.Body)
	output.Body.Close()

	if err != nil {
		f.recorder.Count(MetricS3FetchTotal, 1, []string{"status:error"})
		return nil, err
	}

	var config map[string]auroratype.Parameter
	ext := strings.ToLower(filepath.Ext(f.key))
	if ext == ".yaml" || ext == ".yml" {
		err = yaml.Unmarshal(data, &config)
	} else {
		err = json.Unmarshal(data, &config)
	}

	if err != nil {
		f.recorder.Count(MetricS3FetchTotal, 1, []string{"status:error"})
		return nil, err
	}

	f.recorder.Count(MetricS3FetchTotal, 1, []string{"status:success"})
	return config, nil
}

func (f *Fetcher) IsStatic() bool {
	return false
}

func (f *Fetcher) FetchExperiments(ctx context.Context) ([]auroratype.Experiment, error) {
	if f.experimentsKey == "" {
		return nil, nil
	}

	start := time.Now()
	defer func() {
		duration := float64(time.Since(start).Microseconds())
		f.recorder.Histogram(MetricS3FetchExperimentsLatency, duration, []string{"unit:microseconds"})
	}()

	output, err := f.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(f.bucket),
		Key:    aws.String(f.experimentsKey),
	})
	if err != nil {
		f.recorder.Count(MetricS3FetchExperimentsTotal, 1, []string{"status:error"})
		return nil, err
	}

	data, err := io.ReadAll(output.Body)
	output.Body.Close()

	if err != nil {
		f.recorder.Count(MetricS3FetchExperimentsTotal, 1, []string{"status:error"})
		return nil, err
	}

	var config struct {
		Experiments []auroratype.Experiment `yaml:"experiments"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		f.recorder.Count(MetricS3FetchExperimentsTotal, 1, []string{"status:error"})
		return nil, err
	}

	f.recorder.Count(MetricS3FetchExperimentsTotal, 1, []string{"status:success"})
	return config.Experiments, nil
}
