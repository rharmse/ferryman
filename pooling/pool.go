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

//This represents a HTTP Serving node part of a Resource Pool
type PoolNode struct {
	hostname   string
	ip         net.IP
	port       uint16
	relCtxRoot string
	ctxRoot    string
	nodeURI    string
	scheme     string
}

//Represents a container of HTTP Servers serving client requests
type Pool struct {
	name       string
	members    map[string]*PoolNode
	httpClient *http.Client
}

func (pool *Pool) AddPoolNode(hostname string) error {
	return nil
}

func (pool *Pool) RemovePoolNode(hostname string) error {
	return nil
}


func (node *PoolNode) GetNodeURI() (string, error) {
	if node.nodeURI == nil {
		builder := &strings.Builder{}

		cnt, _ = b.WriteString(node.sheme)
		cnt, _ = b.WriteString("://")
		cnt, _ = b.WriteString(ip)
		cnt, _ = b.WriteString(ctxRoot)
		cnt, _ = b.WriteString(relCtxRoot)
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
