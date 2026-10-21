package client_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/lannuttia/terraform-provider-quay/internal/client"
)

func TestCreateOrganization(t *testing.T) {
	created := false
	srv, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/organization/"):
			created = true
			writeJSON(w, http.StatusCreated, map[string]string{})
		case r.Method == http.MethodGet && strings.Contains(r.URL.Path, "/organization/myorg"):
			writeJSON(w, http.StatusOK, client.Organization{Name: "myorg", Email: "a@b.com"})
		default:
			http.NotFound(w, r)
		}
	})
	_ = srv

	org, err := c.CreateOrganization(context.Background(), client.CreateOrganizationRequest{
		Name:  "myorg",
		Email: "a@b.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !created {
		t.Error("POST /organization/ was never called")
	}
	if org.Name != "myorg" {
		t.Errorf("expected name myorg, got %q", org.Name)
	}
}

func TestGetOrganization_NotFound(t *testing.T) {
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusNotFound, map[string]string{"message": "not found"})
	})
	_, err := c.GetOrganization(context.Background(), "missing")
	if err == nil {
		t.Fatal("expected error")
	}
	if !client.IsNotFound(err) {
		t.Errorf("expected IsNotFound, got %v", err)
	}
}

func TestUpdateOrganization(t *testing.T) {
	updated := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			updated = true
			writeJSON(w, http.StatusOK, map[string]string{})
		case http.MethodGet:
			writeJSON(w, http.StatusOK, client.Organization{Name: "myorg", Email: "new@b.com"})
		}
	})
	org, err := c.UpdateOrganization(context.Background(), "myorg", client.UpdateOrganizationRequest{
		Email: "new@b.com",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !updated {
		t.Error("PUT was never called")
	}
	if org.Email != "new@b.com" {
		t.Errorf("expected email new@b.com, got %q", org.Email)
	}
}

func TestDeleteOrganization(t *testing.T) {
	deleted := false
	_, c := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
		}
	})
	if err := c.DeleteOrganization(context.Background(), "myorg"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !deleted {
		t.Error("DELETE was never called")
	}
}
