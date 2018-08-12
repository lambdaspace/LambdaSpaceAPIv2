package config

import (
	"log"
	"sync"

	"github.com/joeshaw/envdecode"
)

// Configuration stores application configuration options
type Configuration struct {
	Mqtt struct {
		Broker struct {
			Host string `env:"BROKER_HOST,default=localhost"`
			Port string `env:"BROKER_PORT,default=1883"`
		}
		Topic string `env:"MQTT_TOPIC_HACKERS,default=lambdaspace/spacestatus/hackers"`
	}
}

var (
	configInstance = Configuration{}
	configOnce     sync.Once
)

// Load loads app configuration from env variables using defaults
func Load() Configuration {
	configOnce.Do(func() {
		err := envdecode.Decode(&configInstance)
		if err != nil {
			log.Println(err)
		}
	},
	)
	return configInstance
}
