// author: wsfuyibing <websearch@163.com>
// date: 2021-08-04

package app

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	Config *Configuration
)

type (
	Configuration struct {
		Host    string `yaml:"host" json:"host"`
		Name    string `yaml:"name" json:"name"`
		Port    int    `yaml:"port" json:"port"`
		Version string `yaml:"version" json:"version"`

		// Consul
		// enabled status.
		//
		// True is enabled otherwise disabled.
		Consul bool `yaml:"consul" json:"consul"`

		// Consul address.
		//
		// Register service to this address. Read and deregister
		// from this address.
		//
		// Example: 192.168.1.100:8500
		ConsulAddr string `yaml:"consul-addr" json:"consul-addr"`

		// Consul protocol.
		//
		// Accept http, https. Default is http.
		//
		// Example: http
		// Example: https
		ConsulScheme string `yaml:"consul-scheme" json:"consul-scheme"`

		// This application service address.
		//
		// Generate it as service address for consul service
		// registry.
		//
		// Example: 172.16.0.100
		// Example: myapp.example.com
		ConsulServiceAddr string `yaml:"consul-service-addr" json:"consul-service-addr"`

		// This application service port.
		//
		// Generate it as service port for consul service
		// registry.
		//
		// Example: 8080
		ConsulServicePort string `yaml:"consul-service-port" json:"consul-service-port"`

		// Executed address.
		// Example: 172.16.0.100:8080
		Addr string `yaml:"-" json:"-"`

		// Running pid.
		// Example: 3721
		Pid int `yaml:"-" json:"-"`

		// Executed software name.
		// Example: GMD/1.2.3
		Software string `yaml:"-" json:"-"`

		// Config initialized time.
		StartTime time.Time `yaml:"-" json:"-"`
	}
)

// LoadJson
// read json file and assign into fields.
func (o *Configuration) LoadJson(name string) error {
	buf, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, o)
}

// LoadYaml
// read yaml file and assign into fields.
func (o *Configuration) LoadYaml(name string) error {
	buf, err := os.ReadFile(name)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(buf, o)
}

// Update
// execution fields based on assigned.
func (o *Configuration) Update() {
	o.Addr = fmt.Sprintf("%s:%d", o.Host, o.Port)
	o.Software = fmt.Sprintf("%s/%s", o.Name, o.Version)
}

func (o *Configuration) init() *Configuration {
	o.loader()

	o.Pid = os.Getpid()
	o.StartTime = time.Now()
	o.Update()
	return o
}

func (o *Configuration) loader() *Configuration {
	for _, file := range []string{"tmp/app.yaml", "config/app.yaml", "../tmp/app.yaml", "../config/app.yaml"} {
		if o.LoadYaml(file) == nil {
			break
		}
	}
	return o
}
