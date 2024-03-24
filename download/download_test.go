package downloads

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/st19n/mageutils/mgos"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDownload(t *testing.T) {
	t.Run("not found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(404)
		}))
		defer srv.Close()

		opts := DownloadOptions{
			URL:  srv.URL,
			Name: "mybin",
		}
		err := DownloadBinary("bin", opts)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "404 Not Found")
	})

	t.Run("found", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte("ok"))
		}))
		defer srv.Close()

		dest, err := os.MkdirTemp("", "mageutils")
		require.NoError(t, err)
		defer os.RemoveAll(dest)

		opts := DownloadOptions{
			URL:  srv.URL,
			Name: "mybin",
		}
		err = DownloadBinary(dest, opts)
		require.NoError(t, err)
		assert.FileExists(t, filepath.Join(dest, "mybin"+mgos.FileExt()))
	})
}
