package azqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-storage-queue-go/azqueue"
)

var _ Queue = &queue{}

type queue struct {
	queueURL    azqueue.QueueURL
	messagesURL azqueue.MessagesURL
}

func newQueue(queueURL azqueue.QueueURL) *queue {
	messagesURL := queueURL.NewMessagesURL()
	return &queue{queueURL, messagesURL}
}

func (q *queue) MessageCount(ctx context.Context) (int, error) {
	props, err := q.queueURL.GetProperties(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get queue properties: %w", err)
	}
	return int(props.ApproximateMessagesCount()), nil
}

type Message struct {
	ID           string `json:"id" xml:"MessageId"`
	PopReceipt   string `json:"popReceipt" xml:"PopReceipt"`
	DequeueCount int    `json:"dequeueCount" xml:"DequeueCount"`
	Body         string `json:"messageText" xml:"MessageText"`
}

func (q *queue) Dequeue(ctx context.Context, count int, timeout time.Duration) ([]*Message, error) {
	dequeue, err := q.messagesURL.Dequeue(ctx, int32(count), timeout)
	if err != nil {
		return nil, fmt.Errorf("failed to dequeue messages: %w", err)
	}

	numMessages := dequeue.NumMessages()
	if numMessages == 0 {
		return nil, nil
	}

	messages := make([]*Message, numMessages)
	for i := int32(0); i < numMessages; i++ {
		msg := dequeue.Message(i)
		messages[i] = &Message{
			ID:           string(msg.ID),
			PopReceipt:   string(msg.PopReceipt),
			DequeueCount: int(msg.DequeueCount),
			Body:         msg.Text,
		}
	}

	return messages, nil
}

func (q *queue) Enqueue(ctx context.Context, message string, timeout, ttl time.Duration) (*Message, error) {
	msg, err := q.messagesURL.Enqueue(ctx, message, timeout, ttl)
	if err != nil {
		return nil, fmt.Errorf("failed to enqueue message: %w", err)
	}

	return &Message{
		ID:           string(msg.MessageID),
		PopReceipt:   string(msg.PopReceipt),
		DequeueCount: 0,
		Body:         message,
	}, nil
}

func (q *queue) Peek(ctx context.Context, count int) ([]*Message, error) {
	peek, err := q.messagesURL.Peek(ctx, int32(count))
	if err != nil {
		return nil, fmt.Errorf("failed to peek messages: %w", err)
	}

	numMessages := peek.NumMessages()
	if numMessages == 0 {
		return nil, nil
	}

	messages := make([]*Message, numMessages)
	for i := int32(0); i < numMessages; i++ {
		msg := peek.Message(i)
		messages[i] = &Message{
			ID:           string(msg.ID),
			PopReceipt:   "",
			DequeueCount: int(msg.DequeueCount),
			Body:         msg.Text,
		}
	}

	return messages, nil
}

func (q *queue) Delete(ctx context.Context, m *Message) error {
	msgIDURL := q.messagesURL.NewMessageIDURL(azqueue.MessageID(m.ID))
	_, err := msgIDURL.Delete(ctx, azqueue.PopReceipt(m.PopReceipt))
	return err
}
