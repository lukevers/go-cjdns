// Package admin provides easy methods to the cjdns admin interface
package admin

import (
	"encoding/json"
	"github.com/3M3RY/go-bencode"
	//"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os/user"
	"regexp"
	"sync"
	"time"
)

type CjdnsAdminConfig struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Config   string `json:"config,omitempty"`
}

// Conn is an object for interacting with the CJDNS administration port
type Client struct {
	Admin               *Admin
	AdminLog            *AdminLog
	AuthorizedPasswords *AuthorizedPasswords
	Core                *Core
	EthInterface        *EthInterface
	InterfaceController *InterfaceController
	IPTunnel            *IPTunnel
	NodeStore           *NodeStore
	RouterModule        *RouterModule
	SearchRunner        *SearchRunner
	Security            *Security
	SwitchPinger        *SwitchPinger
	UDPInterface        *UDPInterface
	password            string
	addr                *net.UDPAddr
	enc                 *bencode.Encoder
	conn                *net.UDPConn
	mu                  sync.Mutex
	queries             chan *request
	responses           map[string]chan *packet
	logStreams          map[string]chan<- *LogMessage
}

func Connect(config *CjdnsAdminConfig) (admin *Client, err error) {
	if config == nil {
		config = new(CjdnsAdminConfig)
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		rawFile, err := ioutil.ReadFile(u.HomeDir + "/.cjdnsadmin")
		if err != nil {
			return nil, err
		}

		raw, err := stripComments(rawFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(config.Addr),
		Port: config.Port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	admin = &Client{
		password:  config.Password,
		addr:      addr,
		conn:      conn,
		queries:   make(chan *request),
		responses: make(map[string]chan *packet),
	}

	admin.AdminLog = &AdminLog{admin}
	admin.Admin = &Admin{admin}
	admin.AuthorizedPasswords = &AuthorizedPasswords{admin}
	admin.Core = &Core{admin}
	admin.EthInterface = &EthInterface{admin}
	admin.InterfaceController = &InterfaceController{admin}
	admin.IPTunnel = &IPTunnel{admin}
	admin.NodeStore = &NodeStore{admin}
	admin.RouterModule = &RouterModule{admin}
	admin.SearchRunner = &SearchRunner{admin}
	admin.Security = &Security{admin}
	admin.SwitchPinger = &SwitchPinger{admin}
	admin.UDPInterface = &UDPInterface{admin}

	go admin.readFromConn()
	go admin.writeToConn()
	return admin, err
}

const (
	readerChanSize       = 10
	socketReaderChanSize = 100
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
func stripComments(b []byte) ([]byte, error) {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		return nil, err
	}
	out := regComment.ReplaceAllLiteral(b, nil)
	return out, nil
}
