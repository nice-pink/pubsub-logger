package main

import (
	"flag"
	"os"
	"strings"
	"time"

	"github.com/nice-pink/goutil/pkg/log"

	"github.com/nice-pink/pubsub-util/pkg/msg"
	"github.com/nice-pink/pubsub-util/pkg/util"
)

func main() {
	var serviceName = flag.String("serviceName", "pubsub-util", "Will appear in logs.")
	var project = flag.String("project", "", "Google project.")
	var subscription = flag.String("subscription", "", "Subscription.")
	var observeAttr = flag.String("observeAttr", "", "Observe specific attributes.")
	var logHost = flag.String("logHost", "", "Logstash host.")
	var logPort = flag.Int("logPort", 0, "Logstash port.")
	var logLevel = flag.String("logLevel", "debug", "Logstash level.")
	flag.Parse()

	attrs := strings.Split(*observeAttr, ",")
	log.Info(attrs)

	log.Info("*** Start")
	log.Info(os.Args)

	rLog := util.NewRemoteLog(*logHost, *logPort, *logLevel, *serviceName)
	ObservePubSub(*project, *subscription, *observeAttr, rLog)
}

func ObservePubSub(project, subscription, observeAttr string, l *log.RLog) error {
	client, err := msg.NewPubSubHandler(project)
	if err != nil {
		log.Err(err, "could not create pubsub client")
		return nil
	}

	for {
		data, attr, err := client.Subscribe(subscription, 30)
		if err != nil {
			time.Sleep(time.Duration(500) * time.Millisecond)
			continue
		}
		l.Debug(string(data), "::", attr)
	}
}
