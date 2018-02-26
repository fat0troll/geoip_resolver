package cache

import (
	"github.com/fat0troll/geoip_resolver/lib/appcontext"
	"github.com/fat0troll/geoip_resolver/lib/cache/cacheinterface"
	"github.com/fat0troll/geoip_resolver/lib/datamappings"
	"sync"
	"time"
)

var (
	c *appcontext.Context
)

// Cache is an object which handles all the application cached objects
type Cache struct {
	// Geolocation information cache
	// The key is IP address
	geolocations            map[string]*datamappings.CachedGeolocation
	geolocationsMutex       sync.Mutex
	geolocationsFlushTicker *time.Ticker
}

// New passes app context to cache package
func New(ac *appcontext.Context) {
	c = ac
	cd := &Cache{}
	c.RegisterCacheInterface(cacheinterface.CacheInterface(cd))
}

// Init is a initialization function for cache
func (cc *Cache) Init() {
	c.Log.Info("Initializing cache...")

	cc.initGeolocationCacheStorage()
	cc.geolocationsFlushTicker = time.NewTicker(1 * time.Minute)

	go cc.flushOldData()
}
