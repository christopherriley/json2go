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
		var httpErrorWriter HttpErrorWriter

		if errors.As(err, &httpErrorWriter) {
			httpErrorWriter.write(w)
		} else {
			log.Println("*** internal error processing request: ", err)
			ie := NewInternalError("internal error handling request")
			ie.write(w)
		}

		return
	}

	generatedFileComment := "this file was generated"

	log.Println("GET go: pkgName: ", goReq.Package, ", structName: ", goReq.Struct, ", instanceName: ", goReq.Instance)

	generatedCode, err := generate.Generate(generatedFileComment, goReq.Input, goReq.Package, goReq.Struct, goReq.Instance)
	if err != nil {
		re := NewRequestError("code could not be generated", err)
		re.write(w)
		return
	}

	w.Write([]byte(generatedCode))
}
