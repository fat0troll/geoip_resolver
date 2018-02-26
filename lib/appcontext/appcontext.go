package appcontext

import (
	"encoding/json"
	"github.com/Tomasen/realip"
	"github.com/fat0troll/geoip_resolver/lib/cache/cacheinterface"
	"github.com/fat0troll/geoip_resolver/lib/config"
	"github.com/fat0troll/geoip_resolver/lib/requester/requesterinterface"
	"net"
	"net/http"
	"os"
	"source.pztrn.name/golibs/flagger"
	"source.pztrn.name/golibs/mogrus"
)

// Context is an application context struct
type Context struct {
	Cache         cacheinterface.CacheInterface
	Cfg           *config.Config
	HTTPServerMux *http.ServeMux
	Requester     requesterinterface.RequesterInterface
	Log           *mogrus.LoggerHandler
	StartupFlags  *flagger.Flagger
}

// Init is an initialization function for context
func (c *Context) Init() {
	l := mogrus.New()
	l.Initialize()
	c.Log = l.CreateLogger("stdout")
	c.Log.CreateOutput("stdout", os.Stdout, true, "debug")

	c.Cfg = config.New()

	c.StartupFlags = flagger.New(c.Log)
	c.StartupFlags.Initialize()
	c.HTTPServerMux = http.NewServeMux()
}

// InitializeConfig fills config struct with data from given file
func (c *Context) InitializeConfig(configPath string) {
	c.Cfg.Init(c.Log, configPath)
}

// InitializeStartupFlags gives information about available startup flags
func (c *Context) InitializeStartupFlags() {
	configFlag := flagger.Flag{}
	configFlag.Name = "config"
	configFlag.Description = "Configuration file path"
	configFlag.Type = "string"
	configFlag.DefaultValue = "config.yaml"
	err := c.StartupFlags.AddFlag(&configFlag)
	if err != nil {
		c.Log.Errorln(err)
	}
}

// RegisterCacheInterface registers cache interface in application context
func (c *Context) RegisterCacheInterface(ci cacheinterface.CacheInterface) {
	c.Cache = ci
	c.Cache.Init()
}

// RegisterRequesterInterface registers requester interface in application context
func (c *Context) RegisterRequesterInterface(ri requesterinterface.RequesterInterface) {
	c.Requester = ri
	c.Requester.Init()
}

// StartHTTPListener starts HTTP server on given port
func (c *Context) StartHTTPListener() {
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
			response["ip"] = "null"
		}

		responseBody, _ = json.Marshal(response)
		w.Write(responseBody)
	})

	c.Log.Info("HTTP server started at http://" + c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port)
	err := http.ListenAndServe(c.Cfg.HTTPServer.Host+":"+c.Cfg.HTTPServer.Port, c.HTTPServerMux)
	c.Log.Fatalln(err)
}
