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

//Represents the host configuration to expose
type HostConfig struct {
	Hostname  string `json:"hostname"`
	HttpPort  uint16 `json:"httpPort"`
	HttpsPort uint16 `json:"httpsPort"`
}

//Represents the pool configuration including rules, upstram members etc
type PoolConfig struct {
	PoolName       string         `json:"poolName"`
	Domain         string         `json:"domain"`
	CtxRoot        string         `json:"ctxRoot"`
	ServeOn        HostConfig     `json:"serveOn"`
	Members        []MemberConfig `json:"members"`
	RewriteRules   []RuleConfig   `json:"rewriteRules"`
	TempRedirRules []RuleConfig   `json:"tempRedirRules"`
	PermRedirRules []RuleConfig   `json:"permRedirRules"`
	DropRules      []RuleConfig   `json:"dropRules"`
}

//Represents multiple pool configurations
type Config struct {
	Pools []PoolConfig `json:"ferrymanConf"`
}

//Loads the configuration from JSON config file, returns struct.
func LoadConf(pathToConf string) (Config, error) {
	var config = Config{}
	if file, error := os.Stat(pathToConf); file != nil && error == nil {
		raw, error := ioutil.ReadFile(pathToConf)

		if error == nil {
			error = json.Unmarshal(raw, &config)
			fmt.Println("Unmarshalled")
		}
		return config, error
	} else {
		return Config{}, error
	}
}

//Writes the configuration to JSON config file, after potential change
func WriteConf(pathToNewConf string, pathToCurrentConf string, conf Config) error {
	raw, marshalError := json.Marshal(&conf)
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
