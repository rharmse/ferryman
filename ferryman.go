package main

import (
	"fmt"
	"net"
	"net/http"
	"time"
	//"github.com/kavu/go_reuseport"
	"github.com/rharmse/ferryman/lib"
)

func main() {

	// Load from configuration
	conf, error := ferryman.GetConf("ferryman_online.json")
	if error != nil {
		panic(error)
	}

	//Initialize Pools
	pools := ferryman.InitPools(conf)
	fmt.Printf("\n%v", pools["VCOZA"].String())

	router := ferryman.New(pools["VCOZA"])
	//router
	fmt.Println("Router up")

	// 2. Load Rules - Normal handlers, with pre/post middleware
	// 3. Setup Content Rewrites -  post request middleware
	// 4. Setup server & Clients

	//Need prerequest middleware eg, header rewrites
	//need post request middleware eg header rewrites, body content replacement

	server := &http.Server{
		IdleTimeout:  5000 * time.Millisecond,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      router}

	/*
		    listener, err := reuseport.Listen("tcp", "localhost:8080")
		    defer listener.Close()
			err = server.Serve(listener)
	*/

	listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	err = server.Serve(listener)

	if err != nil {
		fmt.Printf("%v", err)
	}
}
