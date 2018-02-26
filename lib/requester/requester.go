package requester

import (
	"errors"
	"net/http"
	"strconv"
	"time"
)

var (
	client = http.Client{Timeout: 10 * time.Second}
)

func (rq *Requester) initServicesStorage() {
	c.Log.Debug("Initializing services storage...")
	rq.services = make(map[string]*ServiceConfiguration)
}

func (rq *Requester) fillServices() {
	c.Log.Debug("Filling services...")
	for i := range c.Cfg.GeoIPServices {
		service := c.Cfg.GeoIPServices[i]
		switch service.ParseMode {
		case "shiny_geoip", "freegeoip":
			serviceConfiguration := ServiceConfiguration{}
			serviceConfiguration.ServiceType = service.ParseMode
			serviceConfiguration.URL = service.URL
			serviceConfiguration.Limit = service.Limit
			serviceConfiguration.LastMinuteRequests = 0

			rq.services[strconv.Itoa(len(rq.services))] = &serviceConfiguration
		default:
			c.Log.Info("Can't accept service with parse mode=" + service.ParseMode + ": unknown parse mode.")
		}

	}
	c.Log.Debug("Found " + strconv.Itoa(len(rq.services)) + " services in configuration...")
	if len(rq.services) == 0 {
		c.Log.Fatal("Can't start: no valid service to get GeoIP found!")
	}
}

func (rq *Requester) flushLimits() {
	for _ = range rq.limitsTicker.C {
		rq.servicesMutex.Lock()
		c.Log.Debug("Flushing limits...")
		for i := range rq.services {
			rq.services[i].LastMinuteRequests = 0
		}
		rq.servicesMutex.Unlock()
	}
}

func (rq *Requester) findFreeService() (*ServiceConfiguration, error) {
	// There we will find service to proceed request. If none, throwing error
	rq.servicesMutex.Lock()
	for i := range rq.services {
		if rq.services[i].Limit > rq.services[i].LastMinuteRequests {
			rq.services[i].LastMinuteRequests++
			rq.servicesMutex.Unlock()
			return rq.services[i], nil
		}
	}

	// If we're here, nothing found
	rq.servicesMutex.Unlock()
	return nil, errors.New("There is no GeoIP service available right now, try again later")
}

// ProcessRequest returns country code, name and state (cached/requested) for given IP
func (rq *Requester) ProcessRequest(IPAddress string) map[string]string {
	result := make(map[string]string)

	IPAddressCachedObject, err := c.Cache.GetCachedDataForIP(IPAddress)
	if err != nil {
		// Cache miss, but it's okay
		c.Log.Info(err.Error())
		c.Log.Debug("Gathering data from remote services...")
		service, err := rq.findFreeService()
		if err != nil {
			c.Log.Error(err.Error())
			result["status"] = "error"
			result["description"] = "There is no valid GeoIP services available right now"
			return result
		}

		c.Log.Debug("Requests, fullfilled by this service in last minute: " + strconv.Itoa(service.LastMinuteRequests))

		switch service.ServiceType {
		case "freegeoip":
			requestURL := service.URL + IPAddress
			gatheredData := rq.getIPFromFreeGeoIP(requestURL)
			return gatheredData
		}
	}

	// Filling data by cached object
	result["status"] = "success"
	result["description"] = "Found IP address data in cache"
	result["ip"] = IPAddressCachedObject.IPAddress
	result["country_code"] = IPAddressCachedObject.CountryCode
	result["country_name"] = IPAddressCachedObject.CountryName
	result["source"] = "cached"

	return result
}
