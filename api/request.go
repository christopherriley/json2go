package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type goRequest struct {
	Input    string
	Package  string
	Struct   string
	Instance string
}

func NewGoRequest(req *http.Request) (goRequest, error) {
	goReq := goRequest{
		Package:  "main",
		Struct:   "Anonymous",
		Instance: "Instance",
	}

	inputBytes, err := io.ReadAll(req.Body)
	if err != nil {
		return goRequest{}, fmt.Errorf("failed to read request body: %s", err)
	}

	req.Body.Close()

	goReq.Input = strings.TrimSpace(string(inputBytes))
	if len(goReq.Input) == 0 {
		return goRequest{}, NewRequestError("must provide source json as request body", nil)
	}

	if qp := req.URL.Query().Get("package"); len(qp) > 0 {
		goReq.Package = qp
	}
	if qp := req.URL.Query().Get("struct"); len(qp) > 0 {
		goReq.Struct = qp
	}
	if qp := req.URL.Query().Get("instance"); len(qp) > 0 {
		goReq.Instance = qp
	}

	return goReq, nil
}
