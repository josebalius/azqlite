package azqlite

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

// Queue represents a queue in Azure Storage. It provides a core set of operations
// to interact with a queue.
type Queue interface {
	// MessageCount returns an approximate number of messages in the queue.
	MessageCount(ctx context.Context) (int, error)

	// Enqueue adds a message to the queue. The timeout represents how long the message
	// should be hidden from consumers. The ttl is how long the message should be kept
	// in the queue. Pass -time.Second for infinite lifetime.
	Enqueue(ctx context.Context, message string, timeout, ttl time.Duration) (*Message, error)

	// Dequeue removes a message from the queue. The timeout represents how long the
	// message should be hidden from consumers. If the message is not deleted, it will
	// reappear in the queue after the timeout.
	Dequeue(ctx context.Context, count int, timeout time.Duration) ([]*Message, error)

	// Peek returns the next set of messages in the queue without removing it. The count
	// represents how many messages to return.
	Peek(ctx context.Context, count int) ([]*Message, error)

	// Delete removes the message from the queue.
	Delete(ctx context.Context, m *Message) error
}

// Service represents a queue service.
type Service struct {
	serviceURL azqueue.ServiceURL
}

// Config is a set of configuration options for the Azure Queue Service.
type Config struct {
	AccountName string
	AccountKey  string

	// AzureServiceURL is optional to override the default Azure Service URL.
	AzureServiceURL string
}

// NewService creates a new Azure Queue Service. It is configured via the supplied Config.
// The config can contain a different AzureServiceURL to use for the service, this is useful for
// testing.
func NewService(cfg Config) (*Service, error) {
	if cfg.AzureServiceURL == "" {
		cfg.AzureServiceURL = "https://%s.queue.core.windows.net"
	}
	credential, err := azqueue.NewSharedKeyCredential(cfg.AccountName, cfg.AccountKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential: %w", err)
	}
	p := azqueue.NewPipeline(credential, azqueue.PipelineOptions{})
	u, err := url.Parse(fmt.Sprintf(cfg.AzureServiceURL, cfg.AccountName))
	if err != nil {
		return nil, fmt.Errorf("failed to parse url: %w", err)
	}
	serviceURL := azqueue.NewServiceURL(*u, p)
	return &Service{serviceURL}, nil
}

// CreateQueue creates a new queue with the given name.
func (s *Service) CreateQueue(ctx context.Context, queueName string) (Queue, error) {
	queueURL := s.serviceURL.NewQueueURL(queueName)
	if _, err := queueURL.Create(ctx, azqueue.Metadata{}); err != nil {
		return nil, err
	}
	return newQueue(queueURL), nil
}

// DeleteQueue deletes the queue with the given name.
func (s *Service) DeleteQueue(ctx context.Context, queueName string) error {
	queueURL := s.serviceURL.NewQueueURL(queueName)
	_, err := queueURL.Delete(ctx)
	return err
}

// NewQueue instantiates an existing queue with the given name.
func (s *Service) NewQueue(queueName string) Queue {
	queueURL := s.serviceURL.NewQueueURL(queueName)
	return newQueue(queueURL)
}
