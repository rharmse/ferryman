package ferryman

import (
	"bytes"
	"net/http"
)

//Router to implement adaptive radix tree https://pdfs.semanticscholar.org/6abf/5107efc723c655956f027b4a67565b048799.pdf
//Reference http://daslab.seas.harvard.edu/classes/cs265/files/presentations/CS265_presentation_Sinyagin.pdf

type Router struct {
}

//Build a new router load rules should be assigned here from
//config
func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {

    w.Header().Add("Served By", "Ferryman")
	w.WriteHeader(200)
	
	buf := new(bytes.Buffer)
	buf.ReadFrom(req.Body)
	b := buf.Bytes()
	req.Body.Close()
	w.Write(b)
}