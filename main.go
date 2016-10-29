package main

import (
	"flag"
	"github.com/contactless/wbgo"
	"io/ioutil"
	"time"
)

const DRIVER_CLIENT_ID = "wbinflux"

func main() {
	brokerAddress := flag.String("broker", "tcp://localhost:1883", "MQTT broker url")
	configPath := flag.String("config", "/etc/wbinflux.conf", "wbinflux config path")
	debug := flag.Bool("debug", false, "Enable debugging")
	flag.Parse()

	if *debug {
		wbgo.SetDebuggingEnabled(true)
	}
	confBytes, err := ioutil.ReadFile(*configPath)
	if err != nil {
		wbgo.Error.Fatalf("can't load wbinflux config: %v", err)
	}
	config, err := ParseConfig(confBytes)
	if err != nil {
		wbgo.Error.Fatalf("can't parse wbinflux config: %v", err)
	}

	influxClient, err := NewRealInfluxClient()
	if err != nil {
		wbgo.Error.Fatalf("can't create influxdb client: %v", err)
	}
	mqttClient := wbgo.NewPahoMQTTClient(*brokerAddress, DRIVER_CLIENT_ID, false)
	mqttClient.Start()
	sender := NewSender(config, influxClient, nil)
	sender.Subscribe(mqttClient)

	for {
		time.Sleep(1 * time.Second)
	}
}
