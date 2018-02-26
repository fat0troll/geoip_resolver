package requester

import (
	"github.com/fat0troll/geoip_resolver/lib/appcontext"
	"github.com/fat0troll/geoip_resolver/lib/requester/requesterinterface"
	"sync"
	"time"
)

var (
	c *appcontext.Context
)

// ServiceConfiguration holds current state of registered service
type ServiceConfiguration struct {
	ServiceType        string
	URL                string
	Limit              int
	LastMinuteRequests int
}

// Requester is a sturct which holds all requests-related stuff
type Requester struct {
	services      map[string]*ServiceConfiguration
	servicesMutex sync.Mutex
	limitsTicker  *time.Ticker
}

// New passes app context to cache package
func New(ac *appcontext.Context) {
	c = ac
	rq := &Requester{}
	c.RegisterRequesterInterface(requesterinterface.RequesterInterface(rq))
}

// Init is a initialization function for requester
func (rq *Requester) Init() {
	c.Log.Info("Initializing requester...")

	rq.initServicesStorage()
	rq.fillServices()

	rq.limitsTicker = time.NewTicker(1 * time.Minute)
	go rq.flushLimits()
}
