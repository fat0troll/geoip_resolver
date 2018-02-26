package cache

import (
	"errors"
	"github.com/fat0troll/geoip_resolver/lib/datamappings"
	"time"
)

func (cc *Cache) initGeolocationCacheStorage() {
	c.Log.Debug("Cache: Initializing empty geolocations storage")
	cc.geolocations = make(map[string]*datamappings.CachedGeolocation)
}

func (cc *Cache) flushOldData() {
	for _ = range cc.geolocationsFlushTicker.C {
		c.Log.Debug("Flushing old IP addresses from cache...")
		oldestPermittedTime := time.Now().Add(-1 * time.Duration(c.Cfg.Cache.ValidPeriod) * time.Minute)
		cc.geolocationsMutex.Lock()
		for i := range cc.geolocations {
			if cc.geolocations[i].CreateTime.Before(oldestPermittedTime) {
				delete(cc.geolocations, i)
			}
		}
		cc.geolocationsMutex.Unlock()
	}
}

// AddIPAddressToCache adds IP address data to cache
func (cc *Cache) AddIPAddressToCache(data *datamappings.CachedGeolocation) {
	cc.geolocationsMutex.Lock()
	cc.geolocations[data.IPAddress] = data
	cc.geolocationsMutex.Unlock()
}

// GetAllCachedIPs returns array of all cached data
func (cc *Cache) GetAllCachedIPs() []datamappings.CachedGeolocation {
	items := make([]datamappings.CachedGeolocation, 0)
	cc.geolocationsMutex.Lock()
	for i := range cc.geolocations {
		items = append(items, *cc.geolocations[i])
	}
	cc.geolocationsMutex.Unlock()
	return items
}

// GetCachedDataForIP returns cached data for IP address if any
func (cc *Cache) GetCachedDataForIP(IPAddress string) (*datamappings.CachedGeolocation, error) {
	cc.geolocationsMutex.Lock()
	for i := range cc.geolocations {
		if cc.geolocations[i].IPAddress == IPAddress {
			cc.geolocationsMutex.Unlock()
			return cc.geolocations[i], nil
		}
	}

	cc.geolocationsMutex.Unlock()
	return nil, errors.New("There is no cached data for IP: " + IPAddress)
}
