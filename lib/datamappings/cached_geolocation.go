package datamappings

import (
	"time"
)

// CachedGeolocation handles geolocation cached data for storing in DataCache
type CachedGeolocation struct {
	IPAddress   string
	CountryCode string
	CountryName string
	CreateTime  time.Time
}
