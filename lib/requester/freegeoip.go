package requester

import (
	"encoding/json"
	"github.com/fat0troll/geoip_resolver/lib/datamappings"
	"time"
)

// FreeGeoIPResponse holds response data for freegeoip.net
type FreeGeoIPResponse struct {
	IP          string  `json:"ip"`
	CountryCode string  `json:"country_code"`
	CountryName string  `json:"country_name"`
	RegionCode  string  `json:"region_code,omitempty"`
	RegionName  string  `json:"region_name,omitempty"`
	City        string  `json:"city,omitempty"`
	ZIP         string  `json:"zip_code,omitempty"`
	TimeZone    string  `json:"time_zone,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	MetroCode   int     `json:"metro_code,omitempty"`
}

func (rq *Requester) getIPFromFreeGeoIP(requestURL string) map[string]string {
	result := make(map[string]string)

	c.Log.Debug("Sending request: " + requestURL)

	req, err := client.Get(requestURL)
	if err != nil {
		c.Log.Error(err.Error())
		result["status"] = "error"
		result["description"] = err.Error()

		return result
	}
	defer req.Body.Close()

	response := FreeGeoIPResponse{}
	err = json.NewDecoder(req.Body).Decode(&response)
	if err != nil {
		c.Log.Error(err.Error())
		result["status"] = "error"
		result["description"] = err.Error()

		return result
	}

	if response.CountryCode != "" {
		result["status"] = "success"
		result["description"] = "Found IP address data in freegeoip service!"
		result["ip"] = response.IP
		result["country_code"] = response.CountryCode
		result["country_name"] = response.CountryName

		cachedData := datamappings.CachedGeolocation{}
		cachedData.IPAddress = response.IP
		cachedData.CountryCode = response.CountryCode
		cachedData.CountryName = response.CountryName
		t := time.Now().UTC()
		cachedData.CreateTime = t
		c.Cache.AddIPAddressToCache(&cachedData)
	} else {
		result["status"] = "error"
		result["description"] = "Service freegeoip doesn't know country of this IP address"
		result["ip"] = response.IP
	}

	return result
}
