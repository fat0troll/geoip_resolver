package requester

import (
	"encoding/json"
	"github.com/Tomasen/realip"
	"github.com/fat0troll/geoip_resolver/lib/appcontext"
	"github.com/fat0troll/geoip_resolver/lib/requester/requesterinterface"
	"net"
	"net/http"
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

	// Registering handler for response
	response := make(map[string]string)
	responseBody := make([]byte, 0)
	c.HTTPServerMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		clientIP := ""
		ipInQuery, ipInQueryFound := r.URL.Query()["ip"]
		if ipInQueryFound && len(ipInQuery) > 0 {
			realIP := net.ParseIP(ipInQuery[0])
			if realIP != nil {
				clientIP = ipInQuery[0]
			}
		} else {
			clientIP = realip.FromRequest(r)
		}

		w.Header().Add("Content-Type", "application/json; charset=utf-8")
		if clientIP != "" {
			response = c.Requester.ProcessRequest(clientIP)
		} else {
			w.WriteHeader(404)
			response["status"] = "error"
			response["description"] = "Invalid IP address query"
		}

		responseBody, _ = json.Marshal(response)
		w.Write(responseBody)
	})
}
