package lambdahttp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlers(t *testing.T) {
	files, err := filepath.Glob("requests/*.json")
	require.NoError(t, err)
	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			require.FileExists(t, file)
			//nolint:gosec // This is a test file reading fixtures
			fd, err := os.Open(file)
			require.NoError(t, err)
			defer func() {
				assert.NoError(t, fd.Close())
			}()

			req := LambdaRequest{}
			require.NoError(t, json.NewDecoder(fd).Decode(&req))

			handler := NewAWSLambdaHTTPHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, "POST", r.Method)
				assert.Contains(t, []string{"/v1/some/path", "/v2/some/path"}, r.URL.Path)

				if !slices.Contains([]string{"requests/alb_url-v1.json", "requests/alb_url-v2.json"}, file) {
					// The case `curl -H "some-header: value-1" -H "some-header: value-2"`
					// is not supported by the ALB handler (i.e. value-1 is not present)
					assert.Contains(t, r.Header.Values("some-header"), "value-1")
				}
				assert.Contains(t, r.Header.Values("some-header"), "value-2")
				assert.Equal(t, "query", r.URL.Query().Get("some"))
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("X-Test", "old")
				w.Header().Set("X-Test", "test")
				w.Header().Add("X-Test", "test-2")
				w.WriteHeader(200)
				_, err := w.Write([]byte(`{"message": "hello"}`))
				require.NoError(t, err)
			}))
			response, err := handler(context.Background(), req)
			require.NoError(t, err)
			assert.Equal(t, 200, response.GetStatusCode())
			assert.Equal(t, "application/json", response.GetHeaders().Get("Content-Type"))
			assert.Equal(t, []string{"test", "test-2"}, response.GetHeaders()["X-Test"])
			assert.Equal(t, `{"message": "hello"}`, response.GetBody())
		})
	}
}
