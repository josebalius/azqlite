package azqlite

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNewService(t *testing.T) {
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	s, err := NewService(Config{
		AccountName: accountName,
		AccountKey:  accountKey,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	if s == nil {
		t.Fatal("NewService returned nil")
	}
	if s.serviceURL.String() != "https://name.queue.core.windows.net" {
		t.Fatalf("NewService returned incorrect serviceURL: %v", s.serviceURL)
	}
}

func TestCreateQueue(t *testing.T) {
	serverCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.WriteHeader(http.StatusCreated)
	}))
	defer testServer.Close()

	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))

	serviceURL := testServer.URL + "/%s"
	s, err := NewService(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}

	ctx := context.Background()
	q, err := s.CreateQueue(ctx, "test")
	if err != nil {
		t.Fatalf("CreateQueue returned error: %v", err)
	}
	if q == nil {
		t.Fatal("CreateQueue returned nil")
	}
	if !serverCalled {
		t.Fatal("CreateQueue did not make HTTP call")
	}
}

func TestDeleteQueue(t *testing.T) {
	serverCalled := false
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serverCalled = true
		w.WriteHeader(http.StatusNoContent)
	}))
	defer testServer.Close()

	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))

	serviceURL := testServer.URL + "/%s"
	s, err := NewService(Config{
		AccountName:     accountName,
		AccountKey:      accountKey,
		AzureServiceURL: serviceURL,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}

	ctx := context.Background()
	err = s.DeleteQueue(ctx, "test")
	if err != nil {
		t.Fatalf("DeleteQueue returned error: %v", err)
	}
	if !serverCalled {
		t.Fatal("DeleteQueue did not make HTTP call")
	}
}

func TestGetQueue(t *testing.T) {
	accountName, accountKey := "name", "key"
	accountKey = base64.StdEncoding.EncodeToString([]byte(accountKey))
	s, err := NewService(Config{
		AccountName: accountName,
		AccountKey:  accountKey,
	})
	if err != nil {
		t.Fatalf("NewService returned error: %v", err)
	}
	q := s.GetQueue("test")
	if q == nil {
		t.Fatal("NewQueue returned nil")
	}
}
