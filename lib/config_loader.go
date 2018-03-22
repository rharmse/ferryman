package ferryman

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type RuleConfig struct {
	FromURI []string `json:"fromURI"`
	ToURI   string   `json:"toURI"`
	Regex   bool     `json:"regex"`
}

type MemberConfig struct {
	Hostname   string `json:"hostname"`
	Port       uint16 `json:"port"`
	RelCtxRoot string `json:"relCtxRoot"`
}

type HostConfig struct {
	Hostname  string `json:"hostname"`
	HttpPort  uint16 `json:"httpPort"`
	HttpsPort uint16 `json:"httpsPort"`
}

type PoolConfig struct {
	PoolName       string         `json:"poolName"`
	Domain         string         `json:"domain"`
	CtxRoot        string         `json:"ctxRoot"`
	ServeOn        HostConfig     `json:serveOn`
	Members        []MemberConfig `json:members`
	RewriteRules   []RuleConfig   `json:rewriteRules`
	TempRedirRules []RuleConfig   `json:tempRedirRules`
	PermRedirRules []RuleConfig   `json:permRedirRules`
	DropRules      []RuleConfig   `json:dropRules`
}

type Config struct {
	Pools []PoolConfig `json:"ferrymanConf"`
}

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

func WriteConf(pathToNewConf string, pathToCurrentConf string, conf Config) error {
	raw, _ := json.Marshal(&conf)
	fInf, _ := os.Stat(pathToCurrentConf)
	return ioutil.WriteFile(pathToNewConf, raw, fInf.Mode())
}
