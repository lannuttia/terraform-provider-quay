package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

func TestCreateRepoNotification(t *testing.T) {
	created := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/notification/") {
			created = true
			writeJSON(w, http.StatusCreated, client.Notification{
				UUID:        "abc-123",
				Event:       "repo_push",
				Method:      "webhook",
				Config:      map[string]any{"url": "https://example.com"},
				EventConfig: map[string]any{},
			})
		}
	})

	notif, err := c.CreateRepoNotification(context.Background(), "myorg/myrepo", client.CreateNotificationRequest{
		Event:       "repo_push",
		Method:      "webhook",
		Config:      map[string]any{"url": "https://example.com"},
		EventConfig: map[string]any{},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Error("POST notification endpoint was never called")
	}
	if notif.UUID != "abc-123" {
		t.Errorf("expected UUID abc-123, got %q", notif.UUID)
	}
}

func TestGetRepoNotification(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			writeJSON(w, http.StatusOK, client.Notification{
				UUID:   "abc-123",
				Event:  "repo_push",
				Method: "webhook",
				Config: map[string]any{"url": "https://example.com"},
			})
		}
	})

	notif, err := c.GetRepoNotification(context.Background(), "myorg/myrepo", "abc-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if notif.UUID != "abc-123" {
		t.Errorf("expected UUID abc-123, got %q", notif.UUID)
	}
}

func TestDeleteRepoNotification(t *testing.T) {
	deleted := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		}
	})
	if err := c.DeleteRepoNotification(context.Background(), "myorg/myrepo", "abc-123"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("DELETE was never called")
	}
}

func TestGetRepoNotification_NotFound(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "not found"})
	})
	_, err := c.GetRepoNotification(context.Background(), "myorg/myrepo", "missing")
	if !client.IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
}
