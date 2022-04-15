package azqlite

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestMessageCount(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("x-ms-approximate-messages-count", "1")
		w.WriteHeader(http.StatusOK)
	}))
	defer testServer.Close()
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	serviceURL := testServer.URL + "/%s"
	s, err := NewClient(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	ctx := context.Background()
	c, err := q.MessageCount(ctx)
	if err != nil {
		t.Fatalf("MessageCount returned error: %v", err)
	}
	if c != 1 {
		t.Fatalf("MessageCount returned %d, expected 0", c)
	}
}

func TestEnqueue(t *testing.T) {
	body := "message-body"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/xml")
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`
			<QueueMessagesList>
				<QueueMessage>
					<MessageId>1</MessageId>
					<InsertionTime>Mon, 02 Jan 2006 15:04:05 MST</InsertionTime>
					<ExpirationTime>Mon, 02 Jan 2006 15:04:05 MST</ExpirationTime>
					<PopReceipt>1</PopReceipt>
					<TimeNextVisible>Mon, 02 Jan 2006 15:04:05 MST</TimeNextVisible>
					<DequeueCount>0</DequeueCount>
					<MessageText>` + body + `</MessageText>
				</QueueMessage>
			</QueueMessagesList>
		`))
	}))
	defer testServer.Close()
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	serviceURL := testServer.URL + "/%s"
	s, err := NewClient(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	ctx := context.Background()
	m, err := q.Enqueue(ctx, body, 1*time.Second, -time.Second)
	if err != nil {
		t.Fatalf("Enqueue returned error: %v", err)
	}
	if m == nil {
		t.Fatal("Enqueue returned nil")
	}
	if m.ID != "1" {
		t.Fatalf("Enqueue returned ID %q, expected 1", m.ID)
	}
	if m.Body != body {
		t.Fatalf("Enqueue returned Body %q, expected %q", m.Body, body)
	}
}

func TestDequeue(t *testing.T) {
	body, body2 := "message-body", "message-body2"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<QueueMessagesList>
				<QueueMessage>
					<MessageId>1</MessageId>
					<InsertionTime>Mon, 02 Jan 2006 15:04:05 MST</InsertionTime>
					<ExpirationTime>Mon, 02 Jan 2006 15:04:05 MST</ExpirationTime>
					<PopReceipt>1</PopReceipt>
					<TimeNextVisible>Mon, 02 Jan 2006 15:04:05 MST</TimeNextVisible>
					<DequeueCount>0</DequeueCount>
					<MessageText>` + body + `</MessageText>
				</QueueMessage>
				<QueueMessage>
					<MessageId>2</MessageId>
					<InsertionTime>Mon, 02 Jan 2006 15:04:05 MST</InsertionTime>
					<ExpirationTime>Mon, 02 Jan 2006 15:04:05 MST</ExpirationTime>
					<PopReceipt>1</PopReceipt>
					<TimeNextVisible>Mon, 02 Jan 2006 15:04:05 MST</TimeNextVisible>
					<DequeueCount>0</DequeueCount>
					<MessageText>` + body2 + `</MessageText>
				</QueueMessage>
			</QueueMessagesList>
		`))
	}))
	defer testServer.Close()
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	serviceURL := testServer.URL + "/%s"
	s, err := NewClient(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	ctx := context.Background()
	messages, err := q.Dequeue(ctx, 30, 1*time.Second)
	if err != nil {
		t.Fatalf("Dequeue returned error: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("Dequeue returned %d messages, expected 2", len(messages))
	}
	if messages[0].ID != "1" {
		t.Fatalf("Dequeue returned ID %q, expected 1", messages[0].ID)
	}
	if messages[0].Body != body {
		t.Fatalf("Dequeue returned Body %q, expected %q", messages[0].Body, body)
	}
	if messages[1].ID != "2" {
		t.Fatalf("Dequeue returned ID %q, expected 2", messages[1].ID)
	}
	if messages[1].Body != body2 {
		t.Fatalf("Dequeue returned Body %q, expected %q", messages[1].Body, body2)
	}
}

func TestPeek(t *testing.T) {
	body := "message-body"
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<QueueMessagesList>
				<QueueMessage>
					<MessageId>1</MessageId>
					<InsertionTime>Mon, 02 Jan 2006 15:04:05 MST</InsertionTime>
					<ExpirationTime>Mon, 02 Jan 2006 15:04:05 MST</ExpirationTime>
					<PopReceipt>1</PopReceipt>
					<TimeNextVisible>Mon, 02 Jan 2006 15:04:05 MST</TimeNextVisible>
					<DequeueCount>0</DequeueCount>
					<MessageText>` + body + `</MessageText>
				</QueueMessage>
			</QueueMessagesList>
		`))
	}))
	defer testServer.Close()
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	serviceURL := testServer.URL + "/%s"
	s, err := NewClient(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	ctx := context.Background()
	messages, err := q.Peek(ctx, 30)
	if err != nil {
		t.Fatalf("Peek returned error: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("Peek returned %d messages, expected 1", len(messages))
	}
	if messages[0].ID != "1" {
		t.Fatalf("Peek returned ID %q, expected 1", messages[0].ID)
	}
	if messages[0].Body != body {
		t.Fatalf("Peek returned Body %q, expected %q", messages[0].Body, body)
	}
}

func TestDelete(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`
			<QueueMessagesList />
		`))
	}))
	defer testServer.Close()
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	serviceURL := testServer.URL + "/%s"
	s, err := NewClient(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
	ctx := context.Background()
	err = q.Delete(ctx, &Message{ID: "1"})
	if err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}
