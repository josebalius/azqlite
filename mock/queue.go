package mock

import (
	"context"
	"time"

	"github.com/josebalius/azqlite"
)

type Queue struct {
	MessageCountFunc    func(ctx context.Context) (int, error)
	MessageCountInvoked bool

	EnqueueFunc    func(ctx context.Context, msg string, timeout, ttl time.Duration) (*azqlite.Message, error)
	EnqueueInvoked bool

	DequeueFunc    func(ctx context.Context, count int, timeout time.Duration) ([]*azqlite.Message, error)
	DequeueInvoked bool

	PeekFunc    func(ctx context.Context, count int) ([]*azqlite.Message, error)
	PeekInvoked bool

	DeleteFunc    func(ctx context.Context, msg *azqlite.Message) error
	DeleteInvoked bool
}

func (q *Queue) MessageCount(ctx context.Context) (int, error) {
	q.MessageCountInvoked = true
	return q.MessageCountFunc(ctx)
}

func (q *Queue) Enqueue(ctx context.Context, msg string, timeout, ttl time.Duration) (*azqlite.Message, error) {
	q.EnqueueInvoked = true
	return q.EnqueueFunc(ctx, msg, timeout, ttl)
}

func (q *Queue) Dequeue(ctx context.Context, count int, timeout time.Duration) ([]*azqlite.Message, error) {
	q.DequeueInvoked = true
	return q.DequeueFunc(ctx, count, timeout)
}

func (q *Queue) Peek(ctx context.Context, count int) ([]*azqlite.Message, error) {
	q.PeekInvoked = true
	return q.PeekFunc(ctx, count)
}

func (q *Queue) Delete(ctx context.Context, msg *azqlite.Message) error {
	q.DeleteInvoked = true
	return q.DeleteFunc(ctx, msg)
}
