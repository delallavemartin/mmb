package handlers_test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

    "go.uber.org/zap/zaptest"

    "mllave.com/mllave/mmb/mb/handlers"
    "mllave.com/mllave/mmb/mb/src/model/subscribers"
)

func TestSubscriberHandler(t *testing.T) {
    subs := &subscribers.SubscribersList{Addresses: make([]string, 0, 10)}
    logger := zaptest.NewLogger(t)
    handler := handlers.SubscriberHandler(subs, logger)

    body := bytes.NewBufferString("8081")
    req := httptest.NewRequest(http.MethodPost, "/subscribe", body)
    w := httptest.NewRecorder()

    handler(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusCreated {
        t.Errorf("expected status %d, got %d", http.StatusCreated, resp.StatusCode)
    }
    if len(subs.Addresses) != 1 || subs.Addresses[0] != "8081" {
        t.Errorf("subscriber not added correctly: %+v", subs.Addresses)
    }
}

func TestPublisherHandler(t *testing.T) {
    subs := &subscribers.SubscribersList{Addresses: []string{"8081"}}
    logger := zaptest.NewLogger(t)
    handler := handlers.PublisherHandler(subs, logger)

    body := bytes.NewBufferString("test message")
    req := httptest.NewRequest(http.MethodPost, "/notify", body)
    w := httptest.NewRecorder()

    handler(w, req)

    resp := w.Result()
    if resp.StatusCode != http.StatusAccepted {
        t.Errorf("expected status %d, got %d", http.StatusAccepted, resp.StatusCode)
    }
}