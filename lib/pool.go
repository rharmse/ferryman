package ferryman

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Track the request, response and relevant http response codes.
type NodeDataTracker struct {
	reqCnt     int //total req cnt
	respCnt    int //total resp cnt
	failCnt    int //total 4xx and 5xx responses
	http1xxCnt int //convert to map maybe
	http2xxCnt int
	http3xxCnt int
	http4xxCnt int
	http5xxCnt int
}

// This represents a upstream HTTP Serving node part of a Resource Pool
type PoolMember struct {
	hostname   string
	ip         net.IP
	port       int
	relCtxRoot string
	ctxRoot    string
	nodeURI    string
	scheme     string
	httpClient *http.Client
	requestCnt uint
}

// Represents an HTTP Session on an upstream server.
type Session struct {
	id            string
	pinned        *PoolMember
	fallback      *PoolMember
	lastRequested time.Time
}

// Represents a container of upstream HTTP Servers
// serving client requests, will be utilized in a round robin
// fashion or least busy server, supports session affinity and
// content rewrites
type Pool struct {
	sticky         bool
	sessionKeyId   string
	sessionIdType  string
	name           string
	members        map[string]*PoolMember
	sessions       map[string]*Session
	contentRewrite map[string]map[string]string
}

//Build the base URI to utilize when interacting with this upstream Server
func (node *PoolMember) buildNodeURI() {
	if node.nodeURI == "" {
		b := &strings.Builder{}
		b.WriteString(node.scheme)
		b.WriteString("://")
		b.WriteString(node.ip.String())
		b.WriteString(":")
		b.WriteString(strconv.Itoa(node.port))
		b.WriteString(node.ctxRoot)
		b.WriteString(node.relCtxRoot)
		node.nodeURI = b.String()
	}
}

//Setup transport for upstream member
func (node *PoolMember) setupClient(conf UpstreamConConfig) {
	node.httpClient = &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: time.Duration(conf.ConTimeout) * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives:   false,
			DisableCompression:  false,
			MaxIdleConns:        conf.MaxIdleCons,
			MaxIdleConnsPerHost: conf.MaxIdleCons,
			IdleConnTimeout:     time.Duration(conf.KeepAliveTime) * time.Second,
		}}
}

// Returns the pool member with the least amount of requests served
func (pool *Pool) getLeastBusy() (member *PoolMember) {
	var mlast *PoolMember
	for _, m := range pool.members {
		if mlast == nil || m.requestCnt <= mlast.requestCnt {
			member = m
			mlast = m
		}
	}
	return member
}

// Returns the current session active for this session id which contains
// a member target or a pool member with the least amount of requests
// served if no session exists
func (pool *Pool) getUpstreamTarget(r *http.Request) (member *PoolMember, session *Session) {
	sessionId := ""

	for _, cookie := range r.Cookies() {
		if cookie.Name == pool.sessionKeyId {
			sessionId = cookie.Value
			break
		}
	}

	if sessionId != "" && pool.sessions[sessionId] != nil {
		return nil, pool.sessions[sessionId]
	} else {
		return pool.getLeastBusy(), nil
	}
}

// Creates a session
func (pool *Pool) createSession(resp *http.Response, member *PoolMember) {
	for _, cookie := range resp.Cookies() {
		if cookie.Name == pool.sessionKeyId {
			session := &Session{
				pinned:        member,
				id:            cookie.Value,
				lastRequested: time.Now(),
			}
			pool.sessions[session.id] = session
			fmt.Printf("\nSession started: %v=%v pinned to host %v\n", cookie.Name, cookie.Value, member.hostname)
			break
		}
	}
}

// Maintains last request time on this session
func (pool *Pool) maintainSession(session *Session) {
	session.lastRequested = time.Now()
}

//Resolve IP associated to member hostname, only grabs first resolved IPV4.
func (node *PoolMember) resolveIP() {
	ips, error := net.LookupIP(node.hostname)
	if error == nil {
		for _, ip := range ips {
			fmt.Printf("ip:%v, %v, %v\n", ip, len(ip.To4()), net.IPv4len)
			if len(ip.To4()) == net.IPv4len {
				node.ip = ip
				break
			}
		}
	} else {
		panic(error)
	}
}

//Creates a pool member and initializes the transport and base uri
func addHTTPPoolMember(contextRoot string, memberConfig MemberConfig, uStreamConf UpstreamConConfig) *PoolMember {
	poolMember := &PoolMember{}
	poolMember.scheme = "http"
	poolMember.hostname = memberConfig.Hostname
	poolMember.port = memberConfig.Port
	poolMember.ctxRoot = contextRoot
	poolMember.relCtxRoot = memberConfig.RelCtxRoot
	poolMember.resolveIP()
	poolMember.setupClient(uStreamConf)
	poolMember.buildNodeURI()
	return poolMember
}

// Creates a pool with its specific members
func addPool(config PoolConfig) *Pool {
	pool := &Pool{}
	if config.Session.Affinity {
		pool = &Pool{
			sessions: make(map[string]*Session, 20000),
		}
	}

	poolMembers := make(map[string]*PoolMember, len(config.Members))
	for _, poolMemConf := range config.Members {
		poolMember := addHTTPPoolMember(config.CtxRoot, poolMemConf, config.UpstrConProf)
		poolMembers[poolMember.hostname] = poolMember
	}

	pool.members = poolMembers
	pool.name = config.PoolName
	pool.sticky = config.Session.Affinity
	pool.sessionKeyId = config.Session.Identifier
	pool.sessionIdType = config.Session.Type

	if config.ContentRewrite != nil && len(config.ContentRewrite) > 0 {
		pool.contentRewrite = make(map[string]map[string]string, len(config.ContentRewrite))

		for _, rwConf := range config.ContentRewrite {
			if pool.contentRewrite[rwConf.ContentType] == nil {
				pool.contentRewrite[rwConf.ContentType] = make(map[string]string, 5)
			}
			pool.contentRewrite[rwConf.ContentType][rwConf.Find] = rwConf.Replace
		}
	}

	return pool
}

//This initializes a Pool, sets up the Pool and Member transports
//and clients, it does not start serving.
func InitPools(config *Config) map[string]*Pool {
	pools := make(map[string]*Pool, len(config.Pools))
	for _, poolConf := range config.Pools {
		pool := addPool(poolConf)
		pools[pool.name] = pool
	}
	return pools
}

func (pool *Pool) String() string {
	b := &strings.Builder{}
	b.WriteString("{\"name\":\"")
	b.WriteString(pool.name)
	b.WriteString("\",")
	b.WriteString("\"members\":[")
	index := len(pool.members)
	for _, member := range pool.members {
		b.WriteString("{\"hostname\":\"")
		b.WriteString(member.hostname)
		b.WriteString("\",")
		b.WriteString("\"ip\":\"")
		b.WriteString(member.ip.String())
		b.WriteString("\",")
		b.WriteString("\"port\":")
		b.WriteString(strconv.Itoa(member.port))
		b.WriteString(",")
		b.WriteString("\"ctxRoot\":\"")
		b.WriteString(member.ctxRoot)
		b.WriteString("\",")
		b.WriteString("\"relCtxRoot\":\"")
		b.WriteString(member.relCtxRoot)
		b.WriteString("\",")
		b.WriteString("\"scheme\":\"")
		b.WriteString(member.scheme)
		b.WriteString("\",")
		b.WriteString("\"nodeURI\":\"")
		b.WriteString(member.nodeURI)
		b.WriteString("\"}")
		if index > 1 {
			b.WriteString(",")
			index--
		}
	}
	b.WriteString("]}")
	return b.String()
}
