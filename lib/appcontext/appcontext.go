package appcontext

import (
	"github.com/fat0troll/geoip_resolver/lib/cache/cacheinterface"
	"github.com/fat0troll/geoip_resolver/lib/config"
	"github.com/fat0troll/geoip_resolver/lib/requester/requesterinterface"
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
	c.Log.Info("HTTP server started at http://" + c.Cfg.HTTPServer.Host + ":" + c.Cfg.HTTPServer.Port)
	err := http.ListenAndServe(c.Cfg.HTTPServer.Host+":"+c.Cfg.HTTPServer.Port, c.HTTPServerMux)
	c.Log.Fatalln(err)
}
