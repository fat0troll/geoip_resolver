package cacheinterface

import (
	"github.com/fat0troll/geoip_resolver/lib/datamappings"
)

// CacheInterface is an interface which represents functions of Cache package
type CacheInterface interface {
	Init()

	AddIPAddressToCache(data *datamappings.CachedGeolocation)
	GetAllCachedIPs() []datamappings.CachedGeolocation
	GetCachedDataForIP(IPAddress string) (*datamappings.CachedGeolocation, error)
}
