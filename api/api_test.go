package api

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type validJsonTestCase struct {
	testName  string
	inputJson string
	inputUri  string
	expect    string
}

var validJsonWithDefaultsTest = validJsonTestCase{
	testName:  "with defaults",
	inputJson: `{"Name": "chris"}`,
	inputUri:  "/go",
	expect: `
// this file was generated
// do not modify

package main

type Anonymous struct {
    Name string
}

var Instance Anonymous = Anonymous{
    Name: "chris",
}
`,
}

var validJsonWithQueryParamsTest = validJsonTestCase{
	testName:  "with query params",
	inputJson: `{"Name": "chris"}`,
	inputUri:  "/go?package=somepkg&struct=TestStruct&instance=myVar",
	expect: `
// this file was generated
// do not modify

package somepkg

type TestStruct struct {
    Name string
}

var myVar TestStruct = TestStruct{
    Name: "chris",
}
`,
}

var validJsonTestCases []validJsonTestCase = []validJsonTestCase{
	validJsonWithDefaultsTest,
	validJsonWithQueryParamsTest,
}

func createTestRequest(t *testing.T, method, url, body string) *http.Request {
	req, err := http.NewRequest(method, url, strings.NewReader(body))
	require.NoError(t, err)

	return req
}

func TestApi(t *testing.T) {
	api := NewApi()

	t.Run("handleGetGo", func(t *testing.T) {
		var responseWriterBuf strings.Builder
		var responseWriterStatusCode int
		w := mockResponseWriter{buf: &responseWriterBuf, statusCode: &responseWriterStatusCode}

		t.Run("invalid json input returns request error", func(t *testing.T) {
			responseWriterBuf.Reset()

			req := createTestRequest(t, "GET", "/go", "invalid json")
			api.handleGetGo(w, req)
			require.Equal(t, http.StatusBadRequest, responseWriterStatusCode)
		})

		t.Run("unreadable body returns internal error", func(t *testing.T) {
			responseWriterBuf.Reset()

			reqWithUnreadableBody, err := http.NewRequest("GET", "/go", mockReadCloser{err: fmt.Errorf("read failed")})
			require.NoError(t, err)
			api.handleGetGo(w, reqWithUnreadableBody)
			require.Equal(t, http.StatusInternalServerError, responseWriterStatusCode)
		})

		t.Run("valid json input", func(t *testing.T) {
			for _, testCase := range validJsonTestCases {
				t.Run(testCase.testName, func(t *testing.T) {
					responseWriterBuf.Reset()

					req := createTestRequest(t, "GET", testCase.inputUri, testCase.inputJson)
					api.handleGetGo(w, req)
					require.Equal(t, http.StatusOK, responseWriterStatusCode)
					assert.Equal(t, strings.TrimSpace(testCase.expect), strings.TrimSpace(responseWriterBuf.String()))
				})
			}
		})
	})
}
