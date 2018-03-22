package ferryman

import (
	"encoding/json"
	"fmt"
	"os"
)

type RuleConfig struct {
	fromURI []string
	toURI   string
	regex   bool
}

type MemberConfig struct {
	hostname   string
	port       uint16
	relCtxRoot string
}

type HostConfig struct {
	hostname  string
	httpPort  uint16
	httpsPort uint16
}

type PoolConfig struct {
	poolName       string
	domain         string
	ctxRoot        string
	serveOn        HostConfig
	members        []MemberConfig
	rewriteRules   []RuleConfig
	tempRedirRules []RuleConfig
	permRedirRules []RuleConfig
	dropRules      []RuleConfig
}

type Config struct {
	ferrymanConf []PoolConfig
}

func LoadConf() (*Config, error) {
	file, _ := os.Open("../ferryman.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Config{}
	err := decoder.Decode(&Config)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Println(configuration)
}
