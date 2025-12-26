package api

import (
	"io"
	"net/http"
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
	if m.err != nil {
		return 0, m.err
	}

	return len(p), io.EOF
}

func (mockReadCloser) Close() error {
	return nil
}
