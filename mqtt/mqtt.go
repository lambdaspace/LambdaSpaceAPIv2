package mqtt

import (
	"fmt"
	"log"
	"strconv"

	"github.com/lambdaspace/LambdaSpaceAPIv2/config"
	"github.com/yosssi/gmq/mqtt/client"
)

var (
	statusChan *chan int
)

func check(err error) {
	if err != nil {
		log.Println(err)
	}
}

func updateHackersCount(_, message []byte) {
	// Cast message to int
	result, err := strconv.Atoi(string(message))
	check(err)
	// Send message to channel
	*statusChan <- result
}

func Main(upstreamChan chan int) {

	statusChan = &upstreamChan

	// Get configuration
	cfg := config.Load()

	// Create an MQTT Client.
	cli := client.New(&client.Options{
		ErrorHandler: check,
	})

	// Terminate the Client.
	defer cli.Terminate()

	// Connect to the MQTT Server.
	err := cli.Connect(&client.ConnectOptions{
		// Network is the network on which the Client connects to.
		Network: "tcp",
		// Address is the address which the Client connects to.
		Address: fmt.Sprintf("%s:%s", cfg.Mqtt.Broker.Host, cfg.Mqtt.Broker.Port),
		// CleanSession is the Clean Session of the CONNECT Packet.
		CleanSession: true,
	})

	check(err)

	// Subscribe to topics.
	err = cli.Subscribe(&client.SubscribeOptions{
		SubReqs: []*client.SubReq{
			&client.SubReq{
				// TopicFilter is the Topic Filter of the Subscription.
				TopicFilter: []byte(cfg.Mqtt.Topic),
				// Handler is the handler which handles the Application Message
				// sent from the Server.
				Handler: updateHackersCount,
			},
		},
	})

	check(err)

}
