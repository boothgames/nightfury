package metrics

import (
	"github.com/influxdata/influxdb-client-go"
	"github.com/influxdata/influxdb-client-go/api/write"
	"net"
	"net/http"
)

var client = influxdb2.NewClient("http://localhost:8086", "")
var org = "boothgames"
var bucket = "nightfury"

// Initialize initialise the metrics server
func Initialize(metricsHost, metricsBucket, authToken string) error {
	client = influxdb2.NewClient(metricsHost, authToken)
	bucket = metricsBucket
	return nil
}

// Write write data to metrics server
func Write(point *write.Point) {
	if client != nil {
		writeAPI := client.WriteAPI(org, bucket)
		writeAPI.WritePoint(point)
	}
}

// IPAddressFromRequest extracts the user IP address from req, if present.
func IPAddressFromRequest(req *http.Request) string {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return "NA"
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return "NA"
	}
	return userIP.String()
}

func RealIPAddressFromRequest(req *http.Request) string {
	return req.Header.Get("X-Real-Ip")
}
