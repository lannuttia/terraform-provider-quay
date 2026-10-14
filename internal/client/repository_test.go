package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

func TestCreateRepository(t *testing.T) {
	created := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/repository"):
			created = true
			writeJSON(w, http.StatusCreated, map[string]string{})
		case r.Method == http.MethodGet:
			writeJSON(w, http.StatusOK, client.Repository{
				Namespace: "myorg",
				Name:      "myrepo",
				IsPublic:  false,
			})
		}
	})

	repo, err := c.CreateRepository(context.Background(), client.CreateRepositoryRequest{
		Repository:  "myrepo",
		Namespace:   "myorg",
		Visibility:  "private",
		Description: "test",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Error("POST /repository was never called")
	}
	if repo.Name != "myrepo" {
		t.Errorf("expected name myrepo, got %q", repo.Name)
	}
	if repo.IsPublic {
		t.Error("expected IsPublic=false")
	}
}

func TestGetRepository_NotFound(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "not found"})
	})
	_, err := c.GetRepository(context.Background(), "myorg/missing")
	if !client.IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
}

func TestUpdateRepository(t *testing.T) {
	updated := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updated = true
			w.WriteHeader(http.StatusOK)
		case http.MethodGet:
			writeJSON(w, http.StatusOK, client.Repository{
				Namespace:   "myorg",
				Name:        "myrepo",
				Description: "updated desc",
			})
		}
	})
	repo, err := c.UpdateRepository(context.Background(), "myorg/myrepo", client.UpdateRepositoryRequest{
		Description: "updated desc",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Error("PUT was never called")
	}
	if repo.Description != "updated desc" {
		t.Errorf("expected description 'updated desc', got %q", repo.Description)
	}
}

func TestChangeRepositoryVisibility(t *testing.T) {
	changed := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "changevisibility") {
			changed = true
			w.WriteHeader(http.StatusOK)
		}
	})
	if err := c.ChangeRepositoryVisibility(context.Background(), "myorg/myrepo", client.ChangeVisibilityRequest{Visibility: "public"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("changevisibility endpoint was never called")
	}
}

func TestDeleteRepository(t *testing.T) {
	deleted := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		}
	})
	if err := c.DeleteRepository(context.Background(), "myorg/myrepo"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("DELETE was never called")
	}
}
