package main

import (
	"fmt"
	"net"
	"net/http"
	"time"

	ferryman "github.com/rharmse/ferryman/lib"
	"golang.org/x/net/netutil"
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
		IdleTimeout:  100 * time.Millisecond,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router}

	listener, error := net.Listen("tcp", ":8080")
	error = server.Serve(netutil.LimitListener(listener, 400))

	if error != nil {
		fmt.Printf("%v", error)
	}
}
