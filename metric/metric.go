// Package metric provides a custom CloudWatch metric for filesystem storage.
package metric

import (
	"errors"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

var (
	ErrMissingName      = errors.New("missing metric name")
	ErrMissingNamespace = errors.New("missing metric namespace")
	ErrMissingRegion    = errors.New("missing metric region")
)

// Identifier represents the identity of a datapoint.
type Identifier struct {
	Filesystem string
	InstanceID string
}

// Metric represents a CloudWatch metric.
type Metric struct {
	client    *cloudwatch.CloudWatch
	Name      string
	Namespace string
	Region    string
}

type storage struct {
	available uint64
	total     uint64
	used      uint64
}

// New creates and returns a metric.
func New(name, namespace, region string) (*Metric, error) {
	if len(name) < 1 {
		return nil, ErrMissingName
	}
	if len(namespace) < 1 {
		return nil, ErrMissingNamespace
	}
	if len(region) < 1 {
		return nil, ErrMissingRegion
	}

	return &Metric{
		client:    cloudwatch.New(session.New(), &aws.Config{Region: &region}),
		Name:      name,
		Namespace: namespace,
		Region:    region,
	}, nil
}

// Publish publishes a datapoint to a metric.
func (m *Metric) Publish(identifier *Identifier) (uint64, error) {
	r, err := statfs(identifier.Filesystem)
	if err != nil {
		return 0, err
	}

	input := &cloudwatch.PutMetricDataInput{
		MetricData: []*cloudwatch.MetricDatum{
			{
				Dimensions: []*cloudwatch.Dimension{
					{
						Name:  aws.String("InstanceId"),
						Value: &identifier.InstanceID,
					},
					{
						Name:  aws.String("Filesystem"),
						Value: &identifier.Filesystem,
					},
				},
				MetricName: &m.Name,
				Unit:       aws.String(string(cloudwatch.StandardUnitBytes)),
				Value:      aws.Float64(float64(r.available)),
			},
		},
		Namespace: &m.Namespace,
	}

	if _, err := m.client.PutMetricData(input); err != nil {
		return 0, err
	}

	return r.available, nil
}

func statfs(fs string) (*storage, error) {
	var s syscall.Statfs_t
	if err := syscall.Statfs(fs, &s); err != nil {
		return nil, err
	}

	return &storage{
		available: s.Bavail * uint64(s.Bsize),
		total:     s.Blocks * uint64(s.Bsize),
		used:      s.Blocks*uint64(s.Bsize) - s.Bavail*uint64(s.Bsize),
	}, nil
}
