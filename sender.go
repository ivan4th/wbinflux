package main

import (
	"github.com/contactless/wbgo"
	"strconv"
	"strings"
	"time"
)

type Clock func() time.Time

type Sender struct {
	config       *Config
	influxClient InfluxClient
	clock        Clock
}

func NewSender(config *Config, influxClient InfluxClient, clock Clock) *Sender {
	if clock == nil {
		clock = time.Now
	}
	return &Sender{
		config:       config,
		influxClient: influxClient,
		clock:        clock,
	}
}

func (sender *Sender) Subscribe(mqttClient wbgo.MQTTClient) {
	// TBD: detect type
	for _, m := range sender.config.Measurements {
		for _, t := range m.Topics {
			measurementName := m.Name
			isNumeric := strings.Contains(t.Topic, "/devices/")
			tags := t.Tags
			mqttClient.Subscribe(func(message wbgo.MQTTMessage) {
				p := Point{
					Measurement: measurementName,
					Tags:        map[string]string{},
					Time:        sender.clock(),
					Fields:      map[string]interface{}{},
					Precision:   "n",
				}
				for name, value := range tags {
					p.Tags[name] = value
				}
				if isNumeric {
					if val, err := strconv.ParseFloat(message.Payload, 64); err == nil {
						p.Fields["value"] = val
					} else {
						wbgo.Error.Printf("Invalid value received for %v: %v", message.Topic, message.Payload)
					}
				} else {
					p.Fields["value"] = message.Payload
				}
				if err := sender.influxClient.WritePoint(p); err != nil {
					wbgo.Error.Printf("Failed to post value %v for %v: %v", message.Topic, message.Payload, err)
				}
			}, t.Topic)
		}
	}
}
