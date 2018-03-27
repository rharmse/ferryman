package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/rharmse/ferryman/lib"
	"golang.org/x/net/netutil"
)

func main() {

	conf, error := ferryman.GetConf()
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	error = ferryman.StoreConf("/home/rharmse/ferryman2.json", conf)
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	error = ferryman.StoreConf("", conf)
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	//ferryman.BootstrapPool(&conf)
	// 1. Determine if there is configuration, and load from it
	// 2. Load Rules - Normal handlers, with pre/post middleware
	// 3. Setup Content Rewrites -  post request middleware
	// 4. Setup server & Clients
    
    //Need prerequest middleware eg, header rewrites
    //need post request middleware eg header rewrites, body content replacement
	router := httprouter.New()
	router.GET("/asd", Index)

	server := &http.Server{
		IdleTimeout:  1 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router}

	listener, error := net.Listen("tcp", ":8080")
	error = server.Serve(netutil.LimitListener(listener, 400))

}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

}
