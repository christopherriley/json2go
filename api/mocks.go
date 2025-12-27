package api

import (
	"io"
	"net/http"
	"strings"
)

type mockResponseWriter struct {
	buf        *strings.Builder
	statusCode *int
}

func (mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m mockResponseWriter) Write(b []byte) (int, error) {
	if _, err := m.buf.Write(b); err != nil {
		return 0, err
	}
	return len(b), nil
}

func (m mockResponseWriter) WriteHeader(statusCode int) {
	*m.statusCode = statusCode
}

type mockReadCloser struct {
	err error
}

func (m mockReadCloser) Read(p []byte) (int, error) {
	if m.err != nil {
		return 0, m.err
	}

	return len(p), io.EOF
}

func (mockReadCloser) Close() error {
	return nil
}
