package main

import (
	"errors"
	"fmt"
	"github.com/influxdata/influxdb/client"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	DEFAULT_INFLUX_HOST = "localhost"
	DEFAULT_INFLUX_PORT = 8086
)

type Point struct {
	Measurement string
	Tags        map[string]string
	Time        time.Time
	Fields      map[string]interface{}
	Precision   string
}

type InfluxClient interface {
	WritePoint(point Point) error
}

type RealInfluxClient struct {
	database    string
	innerClient *client.Client
}

// See: https://github.com/influxdata/influxdb/blob/master/client/example_test.go

func NewRealInfluxClient() (*RealInfluxClient, error) {
	database := os.Getenv("INFLUX_DATABASE")
	if database == "" {
		return nil, errors.New("must specify INFLUX_DATABASE")
	}

	host := os.Getenv("INFLUX_HOST")
	if host == "" {
		host = DEFAULT_INFLUX_HOST
	}

	var err error
	port := DEFAULT_INFLUX_PORT
	portStr := os.Getenv("INFLUX_PORT")
	if portStr != "" {
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, err
		}
	}

	url, err := url.Parse(fmt.Sprintf("http://%s:%d", host, port))
	if err != nil {
		return nil, err
	}

	conf := client.Config{
		URL:      *url,
		Username: os.Getenv("INFLUX_USER"),
		Password: os.Getenv("INFLUX_PWD"),
	}
	innerClient, err := client.NewClient(conf)
	if err != nil {
		return nil, err
	}

	return &RealInfluxClient{
		innerClient: innerClient,
		database:    database,
	}, nil
}

func (c *RealInfluxClient) WritePoint(point Point) error {
	_, err := c.innerClient.Write(client.BatchPoints{
		Points: []client.Point{
			{
				Measurement: point.Measurement,
				Tags:        point.Tags,
				Time:        point.Time,
				Fields:      point.Fields,
				Precision:   point.Precision,
			},
		},
		Database:        c.database,
		RetentionPolicy: "",
	})
	return err
}
