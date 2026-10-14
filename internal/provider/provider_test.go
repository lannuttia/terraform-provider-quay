package provider_test

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/lannuttia/terraform-provider-quay/internal/provider"
)

// testState holds the running Quay instance URL and token, shared across all
// acceptance tests in this package.
var testState struct {
	quayURL   string
	quayToken string
	// runID is a short random suffix appended to all resource names so that
	// parallel matrix jobs (e.g. Terraform 1.13 vs 1.14) never collide on
	// the same org/repo names in a shared Quay instance.
	runID string
}

// rname returns a resource name unique to this test run.
func rname(base string) string {
	return fmt.Sprintf("%s-%s", base, testState.runID)
}

// TestMain starts a Quay testcontainer once for the entire package when
// TF_ACC is set and no external QUAY_URL/QUAY_TOKEN env vars are provided.
func TestMain(m *testing.M) {
	if os.Getenv("TF_ACC") == "" {
		// No acceptance testing — run unit tests (if any) and exit.
		os.Exit(m.Run())
	}

	// Allow overriding with a pre-existing Quay instance.
	if u := os.Getenv("QUAY_URL"); u != "" {
		testState.quayURL = u
		testState.quayToken = os.Getenv("QUAY_TOKEN")
		testState.runID = fmt.Sprintf("%06x", rand.Intn(0xffffff))
		os.Exit(m.Run())
	}

	ctx := context.Background()
	quayC, err := runQuay(ctx, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: could not start Quay testcontainer: %v\n", err)
		os.Exit(1)
	}
	defer quayC.Terminate(ctx) //nolint:errcheck

	testState.quayURL = quayC.baseURL
	testState.quayToken = quayC.accessToken
	testState.runID = fmt.Sprintf("%06x", rand.Intn(0xffffff))

	// Make credentials available to child processes spawned by the test framework.
	if err := os.Setenv("QUAY_URL", testState.quayURL); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: could not set QUAY_URL: %v\n", err)
		os.Exit(1)
	}
	if err := os.Setenv("QUAY_TOKEN", testState.quayToken); err != nil {
		fmt.Fprintf(os.Stderr, "FATAL: could not set QUAY_TOKEN: %v\n", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

// providerFactories returns the provider factory map used by resource.Test.
func providerFactories() map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"quay": providerserver.NewProtocol6WithError(provider.New("test")()),
	}
}

// providerConfig returns an HCL provider block targeting the shared test container.
func providerConfig() string {
	return fmt.Sprintf(`
provider "quay" {
  url   = %q
  token = %q
}
`, testState.quayURL, testState.quayToken)
}

// skipIfNoAcc skips the test if TF_ACC is not set.
func skipIfNoAcc(t *testing.T) {
	t.Helper()
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC not set — skipping acceptance test")
	}
}

// ----- Provider smoke test ------------------------------------------------

func TestProvider_Configure(t *testing.T) {
	skipIfNoAcc(t)
	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: providerFactories(),
		Steps: []resource.TestStep{
			{
				// A minimal config — just verify the provider configures without error.
				Config: providerConfig(),
			},
		},
	})
}
