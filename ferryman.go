package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/kavu/go_reuseport"
	ferryman "github.com/rharmse/ferryman/lib"
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
	// 2. Load Rules - Normal handlers, with pre/post middleware
	// 3. Setup Content Rewrites -  post request middleware
	// 4. Setup server & Clients

	//Need prerequest middleware eg, header rewrites
	//need post request middleware eg header rewrites, body content replacement
	router := ferryman.NewRouter(pools["VCOZA"])
	fmt.Println("Router up")
	server := &http.Server{
		IdleTimeout:  5000 * time.Millisecond,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router}

	listener, err := reuseport.Listen("tcp", "localhost:8080")
	err = server.Serve(listener)

	//listener, error := net.Listen("tcp", ":8080")
	//error = server.Serve(netutil.LimitListener(listener, 400))

	if err != nil {
		fmt.Printf("%v", err)
	}
}
