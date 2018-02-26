package requester

import (
	"encoding/json"
	"github.com/fat0troll/geoip_resolver/lib/datamappings"
	"time"
)

// ShinyGeoIPResponse holds response data for shiny_geoip based services
type ShinyGeoIPResponse struct {
	City     ShinyGeoIPCity     `json:"city,omitempty"`
	Country  ShinyGeoIPCountry  `json:"country,omitempty"`
	Location ShinyGeoIPLocation `json:"location,omitempty"`
	IP       string             `json:"ip,omitempty"`
	Type     string             `json:"type,omitempty"`
	Message  string             `json:"msg,omitempty"`
}

// ShinyGeoIPCountry holds response's `country` fields
type ShinyGeoIPCountry struct {
	CountryName string `json:"name"`
	CountryCode string `json:"code"`
}

// ShinyGeoIPLocation holds response's `location` fields
type ShinyGeoIPLocation struct {
	AccuracyRadius int     `json:"accuracy_radius"`
	Latitude       float64 `json:"latitude"`
	Longitude      float64 `json:"longitude"`
}

// ShinyGeoIPCity is a workaround type for shiny_geoip responses, because they're sending `false` on empty city data
type ShinyGeoIPCity struct {
	City string
}

// UnmarshalJSON for ShinyGeoIPCity implements the workaround for shiny_geoip response
func (sgc *ShinyGeoIPCity) UnmarshalJSON(data []byte) error {
	dataAsString := string(data)
	sgc.City = dataAsString
	return nil
}

func (rq *Requester) getIPFromShinyGeoIP(requestURL string) map[string]string {
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

	response := ShinyGeoIPResponse{}
	err = json.NewDecoder(req.Body).Decode(&response)
	if err != nil {
		c.Log.Error(err.Error())
		result["status"] = "error"
		result["description"] = err.Error()

		return result
	}

	if response.IP != "" {
		if response.Country.CountryCode != "" {
			if response.Country.CountryCode != "" {
				result["status"] = "success"
				result["description"] = "Found IP address data in freegeoip service!"
				result["ip"] = response.IP
				result["country_code"] = response.Country.CountryCode
				result["country_name"] = response.Country.CountryName
				result["source"] = "shiny_geoip"

				cachedData := datamappings.CachedGeolocation{}
				cachedData.IPAddress = response.IP
				cachedData.CountryCode = response.Country.CountryCode
				cachedData.CountryName = response.Country.CountryName
				t := time.Now().UTC()
				cachedData.CreateTime = t
				c.Cache.AddIPAddressToCache(&cachedData)
			} else {
				result["status"] = "error"
				result["description"] = "Service freegeoip doesn't know country of this IP address"
				result["ip"] = response.IP
			}
		}
		result["status"] = "error"
		result["description"] = "Service freegeoip doesn't know this IP address"
	}

	return result
}
