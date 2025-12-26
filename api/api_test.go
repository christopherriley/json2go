package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockResponseWriter struct{}

func (mockResponseWriter) Header() http.Header {
	return http.Header{}
}

func (mockResponseWriter) Write([]byte) (int, error) {
	return -1, nil
}

func (mockResponseWriter) WriteHeader(statusCode int) {
}

type mockReadCloserWithError struct{}

func (mockReadCloserWithError) Read(p []byte) (n int, err error) {
	return 0, fmt.Errorf("read failed")
}

func (mockReadCloserWithError) Close() error {
	return nil
}

func createTestGoRequest(t *testing.T, url, body string) *http.Request {
	req, err := http.NewRequest("GET", url, strings.NewReader(body))
	require.NoError(t, err)

	return req
}

func TestApi(t *testing.T) {
	t.Run("handleGetGo", func(t *testing.T) {
		w := mockResponseWriter{}
		var requestError RequestError
		var internalError InternalError

		t.Run("invalid json input returns request error", func(t *testing.T) {
			goReq, err := NewGoRequest(w, createTestGoRequest(t, "/go", "dsfsdf"))
			require.NoError(t, err)

			_, err = goReq.generateCode()
			require.Error(t, err)
			assert.ErrorAs(t, err, &requestError)
		})

		t.Run("unreadable body returns internal error", func(t *testing.T) {
			reqWithUnreadableBody, err := http.NewRequest("GET", "/go", mockReadCloserWithError{})
			require.NoError(t, err)

			_, err = NewGoRequest(w, reqWithUnreadableBody)
			require.Error(t, err)
			assert.ErrorAs(t, err, &internalError)
		})

		t.Run("valid json input", func(t *testing.T) {
			t.Run("with defaults", func(t *testing.T) {
				input := `{"Name": "chris"}`

				goReq, err := NewGoRequest(w, createTestGoRequest(t, "/go", input))
				require.NoError(t, err)

				expected := `
// this file was generated
// do not modify

package main

type Anonymous struct {
    Name string
}

var Instance Anonymous = Anonymous{
    Name: "chris",
}
`
				actual, err := goReq.generateCode()
				require.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(actual))
			})

			t.Run("with query params", func(t *testing.T) {
				input := `{"Name": "chris"}`

				goReq, err := NewGoRequest(w, createTestGoRequest(t, `/go?package=somepkg&struct=TestStruct&instance=myVar`, input))
				require.NoError(t, err)

				expected := `
// this file was generated
// do not modify

package somepkg

type TestStruct struct {
    Name string
}

var myVar TestStruct = TestStruct{
    Name: "chris",
}
`

				actual, err := goReq.generateCode()
				require.NoError(t, err)
				assert.Equal(t, strings.TrimSpace(expected), strings.TrimSpace(actual))
			})
		})
	})
}
