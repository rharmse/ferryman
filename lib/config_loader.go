package ferryman

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
)

//Represents a rule configuration
type RuleConfig struct {
	FromURI []string `json:"fromURI"`
	ToURI   string   `json:"toURI"`
	Regex   bool     `json:"regex"`
}

//Represents a upstream member configuration.
type MemberConfig struct {
	Hostname   string `json:"hostname"`
	Port       int    `json:"port"`
	RelCtxRoot string `json:"relCtxRoot"`
}

//Represents a upstream connection profile for pool members
type UpstreamConConfig struct {
	MaxCons       int `json:"maxCon"`
	MaxIdleCons   int `json:"maxIdleCon"`
	ConTimeout    int `json:"timeout"`
	KeepAliveTime int `json:"keepaliveTime"`
}

//Represents the host configuration to expose
type HostConfig struct {
	Hostname     string `json:"hostname"`
	HttpPort     int    `json:"httpPort"`
	HttpsPort    int    `json:"httpsPort"`
	ReadTimeout  int    `json:"readTimeout"`
	WriteTimeout int    `json:"writeTimeout"`
	IdleTimeout  int    `json:"idleTimeout"`
}

//Represents the pool configuration including rules, upstream members etc
type PoolConfig struct {
	PoolName       string            `json:"poolName"`
	Domain         string            `json:"domain"`
	CtxRoot        string            `json:"ctxRoot"`
	ServeOn        HostConfig        `json:"serveOn"`
	UpstrConProf   UpstreamConConfig `json:"upstreamCnctConf"`
	Members        []MemberConfig    `json:"members"`
	RewriteRules   []RuleConfig      `json:"rewriteRules"`
	TempRedirRules []RuleConfig      `json:"tempRedirRules"`
	PermRedirRules []RuleConfig      `json:"permRedirRules"`
	DropRules      []RuleConfig      `json:"dropRules"`
}

//Represents multiple pool configurations
type Config struct {
	Pools    []PoolConfig `json:"ferrymanConf"`
	ConfFile string       `json:"-"`
}

//This loads the configuration from users profile directory.
func GetConf() (*Config, error) {
	usrHome := ""
	hasHome := false

	switch opsys := runtime.GOOS; opsys {
	case "windows":
		usrHome, hasHome = os.LookupEnv("USERPROFILE")
	case "linux":
		usrHome, hasHome = os.LookupEnv("HOME")
	default:
		fmt.Printf("OS is => %s\n", opsys)
		return nil, errors.New("Unsupported OS.")
	}

	if !hasHome || "" == usrHome {
		return nil, errors.New("User profile home environment variable not set or present.")
	}

	config, error := LoadConf(usrHome + "/ferryman.json")
	return config, error
}

//Loads the configuration from JSON config file, returns struct.
func LoadConf(pathToConf string) (*Config, error) {
	config := &Config{}
	if file, error := os.Stat(pathToConf); file != nil && error == nil {
		raw, error := ioutil.ReadFile(pathToConf)

		if error == nil {
			error = json.Unmarshal(raw, config)
			config.ConfFile = pathToConf
			fmt.Println("Unmarshalled")
		}
		return config, error
	} else {
		return config, error
	}
}

//Writes the configuration to JSON config file, after potential change
func StoreConf(pathToNewConf string, conf *Config) error {

	fInf, fileError := os.Stat(conf.ConfFile)

	if fileError != nil {
		return fileError
	}

	raw, marshalError := json.Marshal(conf)

	if marshalError == nil {
		if pathToNewConf != "" {
			return ioutil.WriteFile(pathToNewConf, raw, fInf.Mode())
		} else {
			return ioutil.WriteFile(conf.ConfFile, raw, fInf.Mode())
		}
	} else {
		return marshalError
	}
}
