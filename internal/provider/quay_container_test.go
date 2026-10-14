package provider_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	dockercontainer "github.com/moby/moby/api/types/container"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/redis"
	tcnetwork "github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)


const (
	// v3.14.0+ uses SQLite WAL mode, which handles concurrent access from
	// parallel test runs without "database is locked" 500 errors.
	defaultQuayImage = "quay.io/projectquay/quay:v3.17.3"
	quayPort         = "8080/tcp"
	quayAdminUser    = "quayadmin"
	quayAdminPass    = "quayadmin123"
	quayAdminEmail   = "admin@quay.local"
	quayStartTimeout = 5 * time.Minute
)

// quayContainer holds a running Quay instance and the credentials needed to
// interact with it.
type quayContainer struct {
	tc.Container
	accessToken string
	baseURL     string
	tmpDir      string // host temp dir bind-mounted into the container; removed on cleanup
	network     *tc.DockerNetwork
	redis       tc.Container
}

// Cleanup terminates the Quay container, its Redis sidecar, the shared
// network, and the host temp directory.
func (q *quayContainer) Cleanup(ctx context.Context) {
	if q.Container != nil {
		_ = q.Container.Terminate(ctx)
	}
	if q.redis != nil {
		_ = q.redis.Terminate(ctx)
	}
	if q.network != nil {
		_ = q.network.Remove(ctx)
	}
	if q.tmpDir != "" {
		_ = os.RemoveAll(q.tmpDir)
	}
}

// runQuay starts a Quay container with a Redis sidecar on a shared Docker
// network, waits for the health endpoint, then bootstraps an admin user and
// returns an access token.
func runQuay(ctx context.Context, image string) (*quayContainer, error) {
	if image == "" {
		image = defaultQuayImage
	}

	// Isolated network so Quay can reach Redis by its hostname alias.
	net, err := tcnetwork.New(ctx)
	if err != nil {
		return nil, fmt.Errorf("creating docker network: %w", err)
	}

	// Redis must be up before Quay tries to connect.
	redisC, err := redis.Run(ctx, "redis:7-alpine",
		tcnetwork.WithNetwork([]string{"redis"}, net),
	)
	if err != nil {
		_ = net.Remove(ctx)
		return nil, fmt.Errorf("starting redis: %w", err)
	}

	cfg, err := buildQuayConfig()
	if err != nil {
		return nil, fmt.Errorf("building quay config: %w", err)
	}

	// Write config and create storage dir on the host, then bind-mount them
	// into the container. This is the most Podman-compatible approach —
	// CopyToContainer doesn't work on stopped containers under Podman, and
	// ContainerFile trips Podman's path-escapes-parent security check.
	tmpDir, err := os.MkdirTemp("", "quay-test-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp dir: %w", err)
	}
	// tmpDir will be cleaned up when the test process exits; Podman may need
	// it to persist for the container's lifetime so we don't defer removal.

	configDir := filepath.Join(tmpDir, "config")
	storageDir := filepath.Join(tmpDir, "storage")
	for _, d := range []string{configDir, storageDir, filepath.Join(storageDir, "registry")} {
		if err := os.MkdirAll(d, 0o777); err != nil {
			return nil, fmt.Errorf("creating dir %s: %w", d, err)
		}
		// Podman rootless maps host UIDs into a user namespace, so the UID
		// inside the container (Quay runs as 1001) won't match the host UID.
		// Make directories world-writable so any mapped UID can read/write.
		if err := os.Chmod(d, 0o777); err != nil {
			return nil, fmt.Errorf("chmod dir %s: %w", d, err)
		}
	}
	if err := os.WriteFile(filepath.Join(configDir, "config.yaml"), []byte(cfg), 0o644); err != nil {
		return nil, fmt.Errorf("writing config: %w", err)
	}
	// World-readable so the container UID can open it.
	if err := os.Chmod(filepath.Join(configDir, "config.yaml"), 0o644); err != nil {
		return nil, fmt.Errorf("chmod config: %w", err)
	}

	c, err := tc.Run(ctx, image,
		tcnetwork.WithNetwork([]string{"quay"}, net),
		tc.CustomizeRequest(tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				ExposedPorts: []string{quayPort},
				Env: map[string]string{
					"IGNORE_VALIDATION": "true",
				},
				// Use HostConfigModifier to set bind mounts with the :z SELinux
				// label. On Podman/Fedora (SELinux enforcing), the :z flag
				// relabels the host directory so the container process can access
				// it regardless of UID mapping in rootless mode.
				HostConfigModifier: func(hc *dockercontainer.HostConfig) {
					hc.Binds = append(hc.Binds,
						configDir+":"+"/conf/stack"+":z",
						storageDir+":"+"/datastorage"+":z",
					)
				},
				WaitingFor: wait.ForAll(
					wait.ForListeningPort(quayPort).
						WithStartupTimeout(quayStartTimeout),
					wait.ForHTTP("/health/instance").
						WithPort(quayPort).
						WithStatusCodeMatcher(func(status int) bool { return status == http.StatusOK }).
						WithResponseMatcher(func(body io.Reader) bool {
							raw, _ := io.ReadAll(body)
							return strings.Contains(string(raw), `"status_code":200`)
						}).
						WithStartupTimeout(quayStartTimeout),
				),
			},
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("starting quay: %w", err)
	}

	host, err := c.Host(ctx)
	if err != nil {
		return nil, fmt.Errorf("getting quay host: %w", err)
	}
	port, err := c.MappedPort(ctx, quayPort)
	if err != nil {
		return nil, fmt.Errorf("getting quay port: %w", err)
	}
	baseURL := fmt.Sprintf("http://%s:%s", host, port.Port())

	token, err := bootstrapQuayAdmin(ctx, baseURL)
	if err != nil {
		return nil, fmt.Errorf("bootstrapping admin: %w", err)
	}

	return &quayContainer{
		Container:   c,
		accessToken: token,
		baseURL:     baseURL,
		tmpDir:      tmpDir,
		network:     net,
		redis:       redisC,
	}, nil
}

// bootstrapQuayAdmin calls FEATURE_USER_INITIALIZE to create the first admin
// user and returns its OAuth token. This endpoint only works once (while the
// user database is empty).
func bootstrapQuayAdmin(ctx context.Context, baseURL string) (string, error) {
	body, err := json.Marshal(map[string]any{
		"username":     quayAdminUser,
		"password":     quayAdminPass,
		"email":        quayAdminEmail,
		"access_token": true,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		baseURL+"/api/v1/user/initialize",
		bytes.NewReader(body),
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("POST user/initialize: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	raw, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("user/initialize returned %d: %s", resp.StatusCode, raw)
	}

	var result struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("decoding response: %w", err)
	}
	if result.AccessToken == "" {
		return "", fmt.Errorf("no access_token in response: %s", raw)
	}
	return result.AccessToken, nil
}

// quayConfigTmpl is the minimal config.yaml for a test Quay instance.
var quayConfigTmpl = template.Must(template.New("cfg").Parse(`
AUTHENTICATION_TYPE: Database
SERVER_HOSTNAME: localhost:8080
PREFERRED_URL_SCHEME: http
SETUP_COMPLETE: true

SECRET_KEY: "test-secret-key-32-bytes-long-xx"
DATABASE_SECRET_KEY: "test-db-secret-key-32-bytes-xxxx"

DB_URI: "sqlite:////datastorage/quay.db"

BUILDLOGS_REDIS:
  host: redis
  port: 6379

USER_EVENTS_REDIS:
  host: redis
  port: 6379

DISTRIBUTED_STORAGE_CONFIG:
  default:
    - LocalStorage
    - storage_path: /datastorage/registry
DISTRIBUTED_STORAGE_DEFAULT_LOCATIONS: []
DISTRIBUTED_STORAGE_PREFERENCE:
  - default

SUPER_USERS:
  - {{ .AdminUser }}

FEATURE_USER_INITIALIZE: true
BROWSER_API_CALLS_XHR_ONLY: false
FEATURE_MAILING: false
FEATURE_BUILD_SUPPORT: false
FEATURE_SECURITY_SCANNER: false
FEATURE_GARBAGE_COLLECTION: false
FEATURE_DIRECT_LOGIN: true
FEATURE_APP_SPECIFIC_TOKENS: true
`))



func buildQuayConfig() (string, error) {
	var buf strings.Builder
	err := quayConfigTmpl.Execute(&buf, map[string]string{
		"AdminUser": quayAdminUser,
	})
	return buf.String(), err
}
