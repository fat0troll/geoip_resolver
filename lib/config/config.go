package config

import (
	// stdlib
	"io/ioutil"
	"path/filepath"
	// 3rd-party
	"gopkg.in/yaml.v2"
	"source.pztrn.name/golibs/mogrus"
)

// CacheConfiguration handles DataCache configuration
type CacheConfiguration struct {
	ValidPeriod int `yaml:"valid_period"`
}

// GeoIPServiceConfiguration handles GeoIP services API configuration
type GeoIPServiceConfiguration struct {
	URL       string `yaml:"url"`
	ParseMode string `yaml:"parse_mode"`
	Limit     int    `yaml:"limit"`
}

// HTTPServerConfiguration handles HTTP server configuration in config file
type HTTPServerConfiguration struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

// Config is a struct which represents config file structure
type Config struct {
	Cache         CacheConfiguration          `yaml:"cache"`
	GeoIPServices []GeoIPServiceConfiguration `yaml:"geoip_services"`
	HTTPServer    HTTPServerConfiguration     `yaml:"server"`
}

// Init is a configuration initializer
func (c *Config) Init(log *mogrus.LoggerHandler, configPath string) {
	log.Info("Config file path: " + configPath)
	fname, _ := filepath.Abs(configPath)
	yamlFile, yerr := ioutil.ReadFile(fname)
	if yerr != nil {
		log.Fatal("Can't read config file")
	}

	yperr := yaml.Unmarshal(yamlFile, c)
	if yperr != nil {
		log.Fatal("Can't parse config file")
	}
}
