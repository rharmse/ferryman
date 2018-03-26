package ferryman

import (
	"net"
	"net/http"
	"strings"
	"time"
)

//Track the request, response and relevant http response codes.
type NodeDataTracker struct {
	reqCnt     uint64 //total req cnt
	respCnt    uint64 //total resp cnt
	failCnt    uint64 //total 4xx and 5xx responses
	http1xxCnt uint64 //convert to map maybe
	http2xxCnt uint64
	http3xxCnt uint64
	http4xxCnt uint64
	http5xxCnt uint64
}

//This represents a upstream HTTP Serving node part of a Resource Pool
type PoolNode struct {
	hostname   string
	ip         net.IP
	port       uint64
	relCtxRoot string
	ctxRoot    string
	nodeURI    string
	scheme     string
	httpClient *http.Client
}

//Represents a container of upstream HTTP Servers
//serving client requests, will be utilized in a round robin
//fashion or least busy server
type Pool struct {
	name    string
	members map[string]*PoolNode
}

func LoadPoolNodes(poolName string) error {
	return nil
}

func (pool *Pool) RemovePoolNode(hostname string) error {
	return nil
}

func (node *PoolNode) GetNodeURI() (string, error) {
	if node.nodeURI == "" {
		b := &strings.Builder{}

		_, _ = b.WriteString(node.scheme)
		_, _ = b.WriteString("://")
		_, _ = b.WriteString(node.ip.String())
		_, _ = b.WriteString(node.ctxRoot)
		_, _ = b.WriteString(node.relCtxRoot)
		node.nodeURI = b.String()
	}
	return node.nodeURI, nil
}

//Setup transport for upstream member
func (node *PoolNode) setupClient(conf PoolConfig) {
	node.httpClient = &http.Client{
		Timeout: time.Duration(conf.UpstrConProf.ConTimeout) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        conf.UpstrConProf.MaxIdleCons,
			MaxIdleConnsPerHost: conf.UpstrConProf.MaxIdleCons,
			IdleConnTimeout:     time.Duration(conf.UpstrConProf.KeepAliveTime) * time.Second,
		}}
}

//This initializes a Pool, sets up the Pool and Member transports
//and clients, it does not start serving.
func BootstrapPools(config *Config) (map[string]*Pool, error) {
	/*pools := make(map[string]*Pool, len(config.Pools))

	for _, poolConf := range config.Pools {
		pool := &Pool{}

	}
	*/
	return nil, nil
}
