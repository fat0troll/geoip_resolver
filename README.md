# GeoIP resolver

GeoIP resolver — simple server, created for checking IP address geolocation by remote services. Currently supports ``freegeoip.net`` and any service, based on [https://github.com/nekudo/shiny_geoip](shiny_geoip).

## Installation

    go get -u github.com/fat0troll/geoip_resolver/cmd/geoip_resolver

Example configuration file located [https://github.com/fat0troll/geoip_resolver/blob/master/cmd/geoip_resolver/config.yml.dist](here).

## Running

    geoip_resolver --config=/path/to/config.yml

## Testing

1. Clone this repository
2. Change directory to ``cmd/geoip_resolver/``
3. Run ``go test``
