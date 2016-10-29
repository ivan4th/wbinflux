package main

import (
	"errors"
	"fmt"
	"github.com/go-yaml/yaml"
	"strings"
)

type Config struct {
	Measurements []Measurement
}

type Measurement struct {
	Name   string
	Topics []Topic
}

type Topic struct {
	Topic string
	Tags  map[string]string
}

type RawConfig struct {
	Measurements []RawMeasurement
}

func (rc *RawConfig) Convert() (*Config, error) {
	config := &Config{}
	for _, rm := range rc.Measurements {
		m, err := rm.Convert()
		if err != nil {
			return nil, err
		}
		config.Measurements = append(config.Measurements, m)
	}
	return config, nil
}

type RawMeasurement struct {
	Name   string
	Topics []RawTopic
}

func (rm RawMeasurement) Convert() (Measurement, error) {
	m := Measurement{Name: rm.Name}
	for _, rt := range rm.Topics {
		t, err := rt.Convert()
		if err != nil {
			return m, err
		}
		m.Topics = append(m.Topics, t)
	}
	return m, nil
}

type RawTopic struct {
	Name          string
	Topic         string
	DeviceControl string `yaml:"devctl"`
	Tags          map[string]string
}

func (rt RawTopic) Convert() (Topic, error) {
	t := Topic{}

	t.Tags = make(map[string]string)
	for name, value := range rt.Tags {
		t.Tags[name] = value
	}
	if rt.Name != "" {
		t.Tags["name"] = rt.Name
	}

	switch {
	case rt.DeviceControl != "" && rt.Topic != "":
		return t, fmt.Errorf("cannot specify both devctl and topic: %#v", rt)
	case rt.DeviceControl != "":
		parts := strings.SplitN(rt.DeviceControl, "/", 2)
		if len(parts) != 2 {
			return t, errors.New("invalid deviceControl")
		}
		t.Topic = fmt.Sprintf("/devices/%s/controls/%s", parts[0], parts[1])
	case rt.Topic != "":
		t.Topic = rt.Topic
	default:
		return t, fmt.Errorf("must specify either devctl or topic: %#v", rt)
	}

	return t, nil
}

func ParseConfig(in []byte) (*Config, error) {
	var rawConfig RawConfig
	err := yaml.Unmarshal(in, &rawConfig)
	if err != nil {
		return nil, err
	}

	return rawConfig.Convert()
}
