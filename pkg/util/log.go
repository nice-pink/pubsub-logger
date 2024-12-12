package util

import (
	"os"

	"github.com/nice-pink/goutil/pkg/log"
)

func NewRemoteLog(host string, port int, loglevel string, service string) *log.RLog {
	rlog := log.NewRLog(host, port, loglevel, "2006-01-02T01:04:05.100Z", true)
	rlog.UpdateKeys("", "level", "@timestamp")

	common := map[string]interface{}{}
	common["service"] = service
	common["pod_name"] = os.Getenv("POD_NAME")
	common["cluster"] = os.Getenv("CLUSTER_NAME")
	rlog.UpdateCommonData(common)

	if host != "" && port > 0 {
		rlog.Info("Configured logstash. Host:", host, "Port:", port, "Level:", loglevel)
	}

	return rlog
}
