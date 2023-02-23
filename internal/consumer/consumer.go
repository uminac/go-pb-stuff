package consumer

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/uminac/go-pb-stuff/internal/protocol"
	"google.golang.org/protobuf/proto"
)

const BUFFERED_CHANNEL_SIZE = 10000

// Consumer is the interface used for all types of consumers
type Consumer interface {
	Run() error
}

//
// MQTT
//

func NewMQTTConsumer() Consumer {
	return &mqttConsumer{}
}

type mqttConsumer struct {
	// mqttClient is an interface
	mqttClient mqtt.Client

	// mqttOptions is a struct
	mqttOptions *mqtt.ClientOptions
}

func (p *mqttConsumer) Run() error {
	// get a new set of options
	p.mqttOptions = mqtt.NewClientOptions()

	// the broker will always be localhost for this project
	p.mqttOptions.AddBroker("tcp://localhost:1883")

	// set the client ID for this producer (easier to review logs)
	//p.mqttOptions.SetClientID("go-pb-stuff-consumer")

	// automatically reconnect if the connection is lost
	p.mqttOptions.AutoReconnect = true

	// create a MQTT client
	p.mqttClient = mqtt.NewClient(p.mqttOptions)

	// connect to the broker (token.Wait() blocks until connection is established and ready to use)
	if token := p.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		// bail if there's an error
		return token.Error()
	}

	// create a channel to receive messages (buffered to avoid blocking the producer, even though it's probably not necessary in this case)
	recvChannel := make(chan mqtt.Message, BUFFERED_CHANNEL_SIZE)

	// subscribe to the topic
	if token := p.mqttClient.Subscribe("protocol", 0, func(client mqtt.Client, msg mqtt.Message) {
		recvChannel <- msg
	}); token.Wait() && token.Error() != nil {
		// bail if there's an error
		return token.Error()
	}

	// track our sequence number
	var expectedSequenceNumber uint64 = 1

	// create a ticker which ticks every second
	ticker := time.NewTicker(time.Second)

	// track the count of received messages
	var receivedCountSecond uint64 = 0

	// kick off a goroutine to print the count per second
	go func() {
		for range ticker.C {
			fmt.Println("messages received per second:", receivedCountSecond)
			receivedCountSecond = 0 // reset the counter
		}
	}()

	// start our main loop
	for {
		// wait for a message
		msg := <-recvChannel

		// increment our received count
		receivedCountSecond++

		// unmarshal the message payload into a protocol.Action
		action := &protocol.Action{}

		// if we can unmarshal the message, process it
		if err := proto.Unmarshal(msg.Payload(), action); err == nil {
			// if the sequence number is not what we expect, produce an error
			if action.SequenceNumber != expectedSequenceNumber {
				fmt.Println("ERROR: expected sequence number", expectedSequenceNumber, "but got", action.SequenceNumber)
			}

			// reset our sequence number
			expectedSequenceNumber = action.SequenceNumber + 1
		} else {
			// if we can't unmarshal the message, produce an error
			fmt.Println("ERROR: could not unmarshal message:", err)
		}
	}
}
