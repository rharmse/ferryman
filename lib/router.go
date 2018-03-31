package ferryman

import (
	"io"
	"net/http"
    "github.com/kellydunn/go-art"
)

//Router to implement adaptive radix tree https://pdfs.semanticscholar.org/6abf/5107efc723c655956f027b4a67565b048799.pdf
//Reference http://daslab.seas.harvard.edu/classes/cs265/files/presentations/CS265_presentation_Sinyagin.pdf

type Router struct {
	pool *Pool
    routes *ArtTree
}

//Build a new router load rules should be assigned here from
//config
func New(pool *Pool) *Router {
	return &Router{pool: pool}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Served-By", "Ferryman")
	/*member := r.pool.members["hostname"]
	if member != nil {
		resp, err := member.httpClient.Do()
		if err != nil {
			//defer resp.Body.Close()
			io.Copy(w, resp.Body)
		} else {
			fmt.Printf("%v\n", err)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte{})
	}
	return
	fmt.Println("Served here")*/
	w.WriteHeader(http.StatusOK)
	if req.Body != nil {
		defer req.Body.Close()
		io.Copy(w, req.Body)
	}
	w.Write([]byte{})
	return
}
