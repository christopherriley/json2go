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

	mux.HandleFunc("GET /go", handleGetGo)

	return &api
}

func (api Api) GetServeMux() *http.ServeMux {
	return api.mux
}

func handlerError(w http.ResponseWriter, s string, err error) {
	var httpErrorWriter HttpErrorWriter

	if errors.As(err, &httpErrorWriter) {
		httpErrorWriter.write(w)
	} else {
		log.Printf("*** %s: %s\n", s, err)
		ie := NewInternalError(s)
		ie.write(w)
	}
}

func handleGetGo(w http.ResponseWriter, req *http.Request) {
	goReq, err := NewGoRequest(w, req)
	if err != nil {
		handlerError(w, "internal error processing request", err)
		return
	}

	log.Println("GET go: pkgName: ", goReq.Package, ", structName: ", goReq.Struct, ", instanceName: ", goReq.Instance)

	generatedCode, err := goReq.generateCode()
	if err != nil {
		handlerError(w, "internal error generating code", err)
		return
	}

	w.Write([]byte(generatedCode))
}
