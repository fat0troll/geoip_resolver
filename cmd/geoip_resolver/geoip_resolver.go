package main

import (
	"github.com/fat0troll/geoip_resolver/lib/appcontext"
	"github.com/fat0troll/geoip_resolver/lib/cache"
	"github.com/fat0troll/geoip_resolver/lib/requester"
)

func main() {
	c := appcontext.New()
	c.Init()
	c.InitializeStartupFlags()
	c.StartupFlags.Parse()

	configPath, err := c.StartupFlags.GetStringValue("config")
	if err != nil {
		c.Log.Errorln(err)
		c.Log.Fatal("Can't get config file parameter from command line. Exiting.")
	}
	c.InitializeConfig(configPath)

	cache.New(c)
	requester.New(c)

	c.StartHTTPListener()
}
