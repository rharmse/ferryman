package main

import (
    "fmt"
    "github.com/rharmse/ferryman/lib"
)

func main() {

    conf, error := ferryman.LoadConf()
    if error == nil {
        fmt.Println("%v", conf)
    } else {
        panic(error)
    }
    
	// 1. Determine if there is configuration, and load from it
	// 2. Load Rules
	// 3. Setup Content Rewrites
	// 4. Setup server & Clients
}
