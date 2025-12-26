package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (mockResponseWriter) Write([]byte) (int, error) {
	return 0, nil
}

func (mockResponseWriter) WriteHeader(statusCode int) {
}

type mockReadCloser struct {
	err error
}

func (m mockReadCloser) Read(p []byte) (int, error) {
	fmt.Println("Read: len(p): ", len(p))
	if m.err != nil {
		return 0, m.err
	}

	return 0, io.EOF
}

func (mockReadCloser) Close() error {
	return nil
}

func createTestGoRequest(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest("GET", url, strings.NewReader(body))
	require.NoError(t, err)

	return req
}
