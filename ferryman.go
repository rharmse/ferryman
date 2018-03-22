package main

import (
	"fmt"

	"github.com/rharmse/ferryman/lib"
)

func main() {

	conf, error := ferryman.LoadConf("/home/rharmse/go/src/github.com/rharmse/ferryman/ferryman.json")
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	error = ferryman.WriteConf("/home/rharmse/go/src/github.com/rharmse/ferryman/ferryman2.json", "/home/rharmse/go/src/github.com/rharmse/ferryman/ferryman.json", conf)
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	// 1. Determine if there is configuration, and load from it
	// 2. Load Rules
	// 3. Setup Content Rewrites
	// 4. Setup server & Clients
}
