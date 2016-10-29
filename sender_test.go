package main

import (
	"github.com/contactless/wbgo"
	"github.com/contactless/wbgo/testutils"
	"reflect"
	"testing"
	"time"
)

type FakeInfluxClient struct {
	t      *testing.T
	points []Point
}

func NewFakeInfluxClient(t *testing.T) *FakeInfluxClient {
	return &FakeInfluxClient{t: t}
}

func (client *FakeInfluxClient) WritePoint(point Point) error {
	client.points = append(client.points, point)
	return nil
}

func (client *FakeInfluxClient) Verify(points []Point) {
	if !reflect.DeepEqual(points, client.points) {
		client.t.Errorf("Invalid points. Expected:\n%#v\nActual:\n%#v", points, client.points)
	}
}

func TestSender(t *testing.T) {
	broker := testutils.NewFakeMQTTBroker(t, nil)
	outerClient := broker.MakeClient("outer")
	outerClient.Start()
	wbinfluxMqttClient := broker.MakeClient("wbinflux")
	wbinfluxMqttClient.Start()
	fakeInflux := NewFakeInfluxClient(t)
	startTime := time.Now()
	currentTime := startTime
	sender := NewSender(&sampleConfig, fakeInflux, func() time.Time { return currentTime })
	sender.Subscribe(wbinfluxMqttClient)

	messages := []wbgo.MQTTMessage{
		{
			Topic:   "/some/ignored/topic",
			Payload: "123",
		},
		{
			Topic:   "/wbrules/log/info",
			Payload: "Some info message",
		},
		{
			Topic:   "/devices/wb-w1/controls/28-0000058e1692",
			Payload: "-1.5",
		},
	}
	for _, message := range messages {
		outerClient.Publish(message)
		currentTime = currentTime.Add(5 * time.Second)
	}

	fakeInflux.Verify([]Point{
		{
			Measurement: "log",
			Tags:        map[string]string{"level": "info"},
			Time:        startTime.Add(5 * time.Second),
			Fields:      map[string]interface{}{"value": "Some info message"},
			Precision:   "n",
		},
		{
			Measurement: "temps",
			Tags:        map[string]string{"name": "Outside"},
			Time:        startTime.Add(10 * time.Second),
			Fields:      map[string]interface{}{"value": -1.5},
			Precision:   "n",
		},
	})
}
