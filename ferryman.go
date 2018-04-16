package main

import (
	"fmt"
	"net/http"
	"time"
	//"github.com/kavu/go_reuseport"
	reuseport "github.com/kavu/go_reuseport"
	"github.com/rharmse/ferryman/lib"
)

func main() {

	// Load from configuration
	conf, error := ferryman.GetConf("ferryman_test.json")
	if error != nil {
		panic(error)
	}

	//Initialize Pools
	pools := ferryman.InitPools(conf)
	fmt.Printf("\n%v", pools["TEST"].String())

	router := ferryman.New(pools["TEST"])
	//router.AddRoute("/", "/test/home", ferryman.StatusALLMap, ferryman.Default, nil, nil, nil)
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

	listener, err := reuseport.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	err = server.Serve(listener)

	/*listener, err := net.Listen("tcp", ":8080")
	defer listener.Close()
	err = server.Serve(listener)*/

	if err != nil {
		panic(err)
	}
}
