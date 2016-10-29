package main

import (
	"reflect"
	"testing"
)

var sampleConfigStr = `
measurements:
  - name: log
    topics:
    - topic: /wbrules/log/debug
      tags:
        level: debug
    - topic: /wbrules/log/info
      tags:
        level: info
    - topic: /wbrules/log/warning
      tags:
        level: warning
    - topic: /wbrules/log/error
      tags:
        level: error
  - name: temps
    topics:
    - name: Outside
      devctl: wb-w1/28-0000058e1692
    - name: Dining Room
      devctl: wb-w1/28-000007558653
`

var sampleConfig = Config{
	Measurements: []Measurement{
		{
			Name: "log",
			Topics: []Topic{
				{
					Topic: "/wbrules/log/debug",
					Tags: map[string]string{
						"level": "debug",
					},
				},
				{
					Topic: "/wbrules/log/info",
					Tags: map[string]string{
						"level": "info",
					},
				},
				{
					Topic: "/wbrules/log/warning",
					Tags: map[string]string{
						"level": "warning",
					},
				},
				{
					Topic: "/wbrules/log/error",
					Tags: map[string]string{
						"level": "error",
					},
				},
			},
		},
		{
			Name: "temps",
			Topics: []Topic{
				{
					Topic: "/devices/wb-w1/controls/28-0000058e1692",
					Tags: map[string]string{
						"name": "Outside",
					},
				},
				{
					Topic: "/devices/wb-w1/controls/28-000007558653",
					Tags: map[string]string{
						"name": "Dining Room",
					},
				},
			},
		},
	},
}

func verifyConfig(t *testing.T, source string, expectedConfig Config) {
	actualConfig, err := ParseConfig([]byte(source))
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if !reflect.DeepEqual(*actualConfig, expectedConfig) {
		t.Fatalf("Config mismatch. Expected:\n%#v\nActual:\n%#v", *actualConfig, expectedConfig)
	}
}

func TestParseConfigNoPrefix(t *testing.T) {
	verifyConfig(t, sampleConfigStr, sampleConfig)
}
