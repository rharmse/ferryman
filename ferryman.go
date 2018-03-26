package main

import (
	"fmt"
    "runtime"
    "os"
    "errors"
	"github.com/rharmse/ferryman/lib"
)

func GetConf() (*Config, error) {
    usrHome := ""
    hasHome := false
    
    switch opsys := runtime.GOOS; opsys {
        case "windows":
            usrHome, hasHome := os.LookupEnv("USERPROFILE")
        case "linux":
            usrHome, hasHome := os.LookupEnv("HOME")
        default:
            fmt.Printf("OS is => %s\n", opsys)
            return nil, errors.New("Unsupported OS.")
    }
    
    if !hasHome || "" == usrHome
        return nil, errors.New("User profile home environment variable not set or present.")
    
    conf, error := ferryman.LoadConf(userHome + "/ferryman.json")
    
    return &conf, error
}

func main() {
    
    

	conf, error := ferryman.LoadConf("/home/rharmse/go/src/github.com/rharmse/ferryman/ferryman.json")
	if error == nil {
		fmt.Printf("%v", conf)
	} else {
		panic(error)
	}

	/**
	  error = ferryman.WriteConf("c:/Users/harmseru/go/src/github.com/rharmse/ferryman/ferryman2.json", "c:/Users/harmseru/go/src/github.com/rharmse/ferryman/ferryman.json", conf)
		if error == nil {
			fmt.Printf("%v", conf)
		} else {
			panic(error)
		}
	*/
	ferryman.BootstrapPool(&conf)
	// 1. Determine if there is configuration, and load from it
	// 2. Load Rules
	// 3. Setup Content Rewrites
	// 4. Setup server & Clients
}
