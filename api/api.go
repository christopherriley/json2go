package api

import (
	"errors"
	"log"
	"net/http"
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

func (Api) handlerError(w http.ResponseWriter, s string, err error) {
	var httpErrorWriter HttpErrorWriter

	if errors.As(err, &httpErrorWriter) {
		httpErrorWriter.write(w)
	} else {
		log.Printf("*** internal error %s: %s\n", s, err)
		ie := NewInternalError(s)
		ie.write(w)
	}
}

func (api Api) handleGetGo(w http.ResponseWriter, req *http.Request) {
	goReq, err := NewGoRequest(w, req)
	if err != nil {
		api.handlerError(w, "processing request", err)
		return
	}

	log.Println("GET go: pkgName: ", goReq.Package, ", structName: ", goReq.Struct, ", instanceName: ", goReq.Instance)

	generatedCode, err := goReq.generateCode()
	if err != nil {
		api.handlerError(w, "generating code", err)
		return
	}

	w.Write([]byte(generatedCode))
}
