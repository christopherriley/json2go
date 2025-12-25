package api

import (
	"errors"
	"log"
	"net/http"

	"github.com/christopherriley/json2go/generate"
)

type Api struct {
	mux *http.ServeMux
}

func NewApi() *Api {
	mux := http.NewServeMux()

	api := Api{
		mux: mux,
	}

	mux.HandleFunc("GET /go", api.handleGetGo)

	return &api
}

func (api Api) GetServeMux() *http.ServeMux {
	return api.mux
}

func (Api) handleGetGo(w http.ResponseWriter, req *http.Request) {
	goReq, err := NewGoRequest(req)
	if err != nil {
		var httpError HttpError

		if errors.As(err, &httpError) {
			httpError.write(w)
		} else {
			log.Println("*** internal error processing request: ", err)
			ie := InternalError{"internal error handling request"}
			ie.write(w)
		}

		return
	}

	generatedFileComment := "this file was generated"

	log.Println("GET go: pkgName: ", goReq.Package, ", structName: ", goReq.Struct, ", instanceName: ", goReq.Instance)

	generatedCode, err := generate.Generate(generatedFileComment, goReq.Input, goReq.Package, goReq.Struct, goReq.Instance)
	if err != nil {
		re := RequestError{
			Err:   "code could not be generated",
			Cause: err.Error(),
		}

		re.write(w)
		return
	}

	w.Write([]byte(generatedCode))
}
