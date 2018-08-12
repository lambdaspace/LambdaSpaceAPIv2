package config

import (
	"log"
	"sync"

	"github.com/joeshaw/envdecode"
)

type configuration struct {
	Mqtt struct {
		Broker struct {
			Host string `env:"BROKER_HOST,default=localhost"`
			Port string `env:"BROKER_PORT,default=1883"`
		}
		Topic string `env:"MQTT_TOPIC_HACKERS,default=lambdaspace/spacestatus/hackers"`
	}
}

var (
	configInstance = configuration{}
	configOnce     sync.Once
)

func Load() configuration {
	configOnce.Do(func() {
		err := envdecode.Decode(&configInstance)
		if err != nil {
			log.Println(err)
		}
	},
	)
	return configInstance
}
