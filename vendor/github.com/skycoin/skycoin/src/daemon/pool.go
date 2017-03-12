package daemon

import (
	"time"

	//"github.com/skycoin/skycoin/src/daemon/gnet"
	"github.com/skycoin/skycoin/src/daemon/gnet"
)

type PoolConfig struct {
	// Timeout when trying to connect to new peers through the pool
	DialTimeout time.Duration
	// How often to process message buffers and generate events
	MessageHandlingRate time.Duration
	// How long to wait before sending another ping
	PingRate time.Duration
	// How long a connection can idle before considered stale
	IdleLimit time.Duration
	// How often to check for needed pings
	IdleCheckRate time.Duration
	// How often to check for stale connections
	ClearStaleRate time.Duration
	// Buffer size for gnet.ConnectionPool's network Read events
	EventChannelSize int
	// These should be assigned by the controlling daemon
	address string
	port    int
}

func NewPoolConfig() PoolConfig {
	//defIdleLimit := time.Minute
	return PoolConfig{
		port:                6677,
		address:             "",
		DialTimeout:         time.Second * 30,
		MessageHandlingRate: time.Millisecond * 50,
		PingRate:            5 * time.Second,
		IdleLimit:           60 * time.Second,
		IdleCheckRate:       1 * time.Second,
		ClearStaleRate:      1 * time.Second,
		EventChannelSize:    4096,
	}
}

type Pool struct {
	Config PoolConfig
	Pool   *gnet.ConnectionPool
}

func NewPool(c PoolConfig) *Pool {
	return &Pool{
		Config: c,
		Pool:   nil,
	}
}

// Begins listening on port for connections and periodically scanning for
// messages on read_interval
func (self *Pool) Init(d *Daemon) {
	logger.Info("InitPool on port %d", self.Config.port)
	cfg := gnet.NewConfig()
	cfg.DialTimeout = self.Config.DialTimeout
	cfg.Port = uint16(self.Config.port)
	cfg.Address = self.Config.address
	cfg.ConnectCallback = d.onGnetConnect
	cfg.DisconnectCallback = d.onGnetDisconnect
	// cfg.EventChannelSize = cfg.EventChannelSize
	pool := gnet.NewConnectionPool(cfg, d)
	self.Pool = pool
}

// Closes all connections and stops listening
func (self *Pool) Shutdown() {
	if self.Pool != nil {
		self.Pool.Shutdown()
		logger.Info("Shutdown pool")
	}
}

// Starts listening on the configured Port
// no goroutine
func (self *Pool) Start() {
	self.Pool.Run()
}

// Accepts connections, run in goroutine
// func (self *Pool) AcceptConnections() {
// 	self.Pool.AcceptConnections()
// }

// Send a ping if our last message sent was over pingRate ago
func (pool *Pool) sendPings() {
	pool.Pool.SendPings(pool.Config.PingRate, &PingMessage{})
}

// Removes connections that have not sent a message in too long
func (self *Pool) clearStaleConnections() {
	self.Pool.ClearStaleConnections(self.Config.IdleLimit, DisconnectIdle)
}
