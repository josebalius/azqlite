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

// Client represents a queue service.
type Client struct {
	serviceURL azqueue.ServiceURL
}

// Config is a set of configuration options for the Azure Queue Service.
type Config struct {
	AccountName string
	AccountKey  string

	// AzureServiceURL is optional to override the default Azure Service URL.
	AzureServiceURL string
}

// NewClient creates a new Azure Queue Service. It is configured via the supplied Config.
// The config can contain a different AzureServiceURL to use for the service, this is useful for
// testing.
func NewClient(cfg Config) (*Client, error) {
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
	return &Client{serviceURL}, nil
}

// CreateQueue creates a new queue with the given name.
func (c *Client) CreateQueue(ctx context.Context, queueName string) (Queue, error) {
	queueURL := c.serviceURL.NewQueueURL(queueName)
	if _, err := queueURL.Create(ctx, azqueue.Metadata{}); err != nil {
		return nil, err
	}
	return newQueue(queueURL), nil
}

// DeleteQueue deletes the queue with the given name.
func (c *Client) DeleteQueue(ctx context.Context, queueName string) error {
	queueURL := c.serviceURL.NewQueueURL(queueName)
	_, err := queueURL.Delete(ctx)
	return err
}

// GetQueue returns an existing queue with the given name.
func (c *Client) GetQueue(queueName string) Queue {
	queueURL := c.serviceURL.NewQueueURL(queueName)
	return newQueue(queueURL)
}
