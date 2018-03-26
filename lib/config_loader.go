package ferryman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	Port       uint16 `json:"port"`
	RelCtxRoot string `json:"relCtxRoot"`
}

//Represents a upstream connection profile for pool members
type UpstreamConConfig struct {
	MaxCons       uint16 `json:"maxCon"`
	MaxIdleCons   uint32 `json:"maxIdleCon"`
	ConTimeout    uint16 `json:"timeout"`
	KeepAliveTime uint16 `json:"keepaliveTime"`
}

//Represents the host configuration to expose
type HostConfig struct {
	Hostname     string `json:"hostname"`
	HttpPort     uint16 `json:"httpPort"`
	HttpsPort    uint16 `json:"httpsPort"`
	ReadTimeout  uint16 `json:"readTimeout"`
	WriteTimeout uint16 `json:"writeTimeout"`
	IdleTimeout  uint16 `json:"idleTimeout"`
}

//Represents the pool configuration including rules, upstram members etc
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
	Pools []PoolConfig `json:"ferrymanConf"`
}

//This loads the configuration from users profile directory.
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
    
    config, error := LoadConf(userHome + "/ferryman.json")
    
    return config, error
}

//Loads the configuration from JSON config file, returns struct.
func LoadConf(pathToConf string) (*Config, error) {
    config := &Config{}
	if file, error := os.Stat(pathToConf); file != nil && error == nil {
		raw, error := ioutil.ReadFile(pathToConf)

		if error == nil {
			error = json.Unmarshal(raw, config)
			fmt.Println("Unmarshalled")
		}
		return &config, error
	} else {
		return &Config{}, error
	}
}

//Writes the configuration to JSON config file, after potential change
func StoreConf(pathToNewConf string, pathToCurrentConf string, conf *Config) error {
	raw, marshalError := json.Marshal(conf)
	fInf, fileError := os.Stat(pathToCurrentConf)

	if marshalError != nil && fileError != nil {
		return ioutil.WriteFile(pathToNewConf, raw, fInf.Mode())
	}
	if marshalError != nil {
		return marshalError
	} else {
		return fileError
	}
}
