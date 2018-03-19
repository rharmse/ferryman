package ferryman

import (
	"net"
	"net/http"
	"time"
)

type PoolNode struct {
	hostname string
	ip       net.IP
	port     uint16
}

type Pool struct {
	name       string
	members    map[string]*PoolNofde
	httpClient *http.Client
}

func (pool *Pool) AddPoolNode(hostname string) error {

}

func (pool *Pool) RemovePoolNode(hostname string) error {
	return nil
}

func (pool *Pool) SetupClient() error {
	pool.httpClient = &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        len(pool.members) * 110,
			MaxIdleConnsPerHost: 100}}
	return nil
}
