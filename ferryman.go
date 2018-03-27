package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/rharmse/ferryman/lib"
)

func main() {

	conf, error := ferryman.GetConf("ferryman_online.json")
	if error != nil {
		panic(error)
	}

	pools := ferryman.InitPools(conf)
	fmt.Printf("%v", pools["VCOZA"].String())

	//ferryman.BootstrapPool(&conf)
	// 1. Determine if there is configuration, and load from it
	// 2. Load Rules
	// 3. Setup Content Rewrites
	// 4. Setup server & Clients
	/*router := httprouter.New()
	router.GET("/asd", Index)

	server := &http.Server{
		IdleTimeout:  1 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router}

	listener, error := net.Listen("tcp", ":8080")
	error = server.Serve(netutil.LimitListener(listener, 400))
	*/
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
