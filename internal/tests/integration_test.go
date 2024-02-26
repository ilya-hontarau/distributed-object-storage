// TODO: use build tags
package tests

import (
	"io"
	"math/rand"
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUploadDownload(t *testing.T) {
	t.Run("happy path", func(t *testing.T) {
		id := generateCorrectID()
		wantBody := "test"

		putObjectBody(t, id, wantBody, http.StatusOK)

		body := getObjectBody(t, id, http.StatusOK)
		assert.Equal(t, wantBody, body)
	})
	t.Run("empty upload id", func(t *testing.T) {
		putObjectBody(t, " ", "", http.StatusBadRequest)
	})
	t.Run("too long upload id", func(t *testing.T) {
		putObjectBody(t, uuid.NewString(), "", http.StatusBadRequest)
	})
	t.Run("empty download id", func(t *testing.T) {
		getObjectBody(t, " ", http.StatusBadRequest)
	})
	t.Run("too long download id", func(t *testing.T) {
		getObjectBody(t, uuid.NewString(), http.StatusBadRequest)
	})

	t.Run("not existing id", func(t *testing.T) {
		getObjectBody(t, generateCorrectID(), http.StatusNotFound)
	})
	t.Run("upload big body", func(t *testing.T) {
		body50mb := randomString(1024 * 1024 * 50)
		id := generateCorrectID()

		putObjectBody(t, id, body50mb, http.StatusOK)

		body := getObjectBody(t, id, http.StatusOK)
		assert.Equal(t, body50mb, body)
	})
	t.Run("upload too big body", func(t *testing.T) {
		body60mb := randomString(1024 * 1024 * 60)
		id := generateCorrectID()

		putObjectBody(t, id, body60mb, http.StatusBadRequest)
	})
}

func BenchmarkUploadDownloadSmallBody(b *testing.B) {
	body1kb := randomString(1024)
	for i := 0; i < b.N; i++ {
		id := generateCorrectID()
		putObjectBody(b, id, body1kb, http.StatusOK)

		body := getObjectBody(b, id, http.StatusOK)
		assert.Equal(b, body1kb, body)
	}
}

func BenchmarkUploadDownloadBifBody(b *testing.B) {
	body1kb := randomString(1024 * 1024 * 50)
	for i := 0; i < b.N; i++ {
		id := generateCorrectID()
		putObjectBody(b, id, body1kb, http.StatusOK)

		body := getObjectBody(b, id, http.StatusOK)
		assert.Equal(b, body1kb, body)
	}
}

func generateCorrectID() string {
	return uuid.NewString()[:32]
}

func putObjectBody(t require.TestingT, id string, body string, expStatusCode int) {
	req, err := http.NewRequest(http.MethodPut, "http://localhost:3000/object/"+id, strings.NewReader(body))
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, expStatusCode, resp.StatusCode)
	resp.Body.Close()
}

func getObjectBody(t require.TestingT, id string, expStatusCode int) string {

	resp, err := http.Get("http://localhost:3000/object/" + id)
	require.NoError(t, err)
	defer resp.Body.Close()

	assert.Equal(t, expStatusCode, resp.StatusCode)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	return string(respBody)
}

func randomString(length int) string {
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
