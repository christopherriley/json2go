package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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
	pkgName := "main"
	structName := "Anonymous"
	instanceName := "Instance"

	inputBytes, err := io.ReadAll(req.Body)
	if err != nil {
		log.Println("** Error: cannot get read request body: ", err)
		http.Error(w, fmt.Sprintf("cannot read request body: %s", err), http.StatusInternalServerError)
		return
	}

	req.Body.Close()

	input := strings.TrimSpace(string(inputBytes))
	if len(input) == 0 {
		log.Println("** Error: must provide json in request body: ", err)
		http.Error(w, "must provide json in request body", http.StatusBadRequest)
		return
	}

	if pv := req.URL.Query().Get("package"); len(pv) > 0 {
		pkgName = pv
	}
	if pv := req.URL.Query().Get("struct"); len(pv) > 0 {
		structName = pv
	}
	if pv := req.URL.Query().Get("instance"); len(pv) > 0 {
		instanceName = pv
	}

	generatedFileComment := "this file was generated"

	log.Println("GET go: pkgName: ", pkgName, ", structName: ", structName, ", instanceName: ", instanceName)

	generatedCode, err := generate.Generate(generatedFileComment, input, pkgName, structName, instanceName)
	if err != nil {
		log.Println("** Error: code could not be generated: ", err)
		http.Error(w, fmt.Sprintf("code could not be generated: %s", err), http.StatusBadRequest)
		return
	}

	w.Write([]byte(generatedCode))
}
