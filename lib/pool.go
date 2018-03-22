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
	port       uint16
	relCtxRoot string
	ctxRoot    string
	nodeURI    string
	scheme     string
}

//Represents a container of upstream HTTP Servers
//serving client requests, will be utilized in a round robin
//fashion or least busy server
type Pool struct {
	name       string
	members    map[string]*PoolNode
	httpClient *http.Client
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

func (pool *Pool) SetupClient() error {
	pool.httpClient = &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        len(pool.members) * 110,
			MaxIdleConnsPerHost: 100}}
	return nil
}
