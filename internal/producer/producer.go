package producer

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
	"github.com/uminac/go-pb-stuff/internal/protocol"
	"google.golang.org/protobuf/proto"
)

// Producer is the interface used for all types of producers
type Producer interface {
	Run() error
}

//
// MQTT
//

func NewMQTTProducer() Producer {
	return &mqttProducer{}
}

type mqttProducer struct {
	// mqttClient is an interface
	mqttClient mqtt.Client

	// mqttOptions is a struct
	mqttOptions *mqtt.ClientOptions
}

func (p *mqttProducer) Run() error {
	// get a new set of options
	p.mqttOptions = mqtt.NewClientOptions()

	// the broker will always be localhost for this project
	p.mqttOptions.AddBroker("tcp://localhost:1883")

	// set the client ID for this producer (easier to review logs)
	//p.mqttOptions.SetClientID("go-pb-stuff-producer")

	// automatically reconnect if the connection is lost
	p.mqttOptions.AutoReconnect = true

	// create a MQTT client
	p.mqttClient = mqtt.NewClient(p.mqttOptions)

	// connect to the broker (token.Wait() blocks until connection is established and ready to use)
	if token := p.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		// bail if there's an error
		return token.Error()
	}

	// track our sequence number
	var sequenceNumber uint64 = 0

	// create a ticker which ticks every second
	ticker := time.NewTicker(time.Second)

	// track the count of received messages
	var sentCountSecond uint64 = 0

	// kick off a goroutine to print the count per second
	go func() {
		for range ticker.C {
			fmt.Println("messages sent per second:", sentCountSecond)
			sentCountSecond = 0 // reset the counter
		}
	}()

	// start our main loop
	for {
		// increment the sequence number even if we don't publish a message
		sequenceNumber++

		if p.mqttClient.IsConnected() {
			// create a Thing
			thing := protocol.Thing{
				Name: "Thing",
				Type: protocol.Thing_TYPEA,
			}

			// create an Action
			action := protocol.Action{
				Uuid:           uuid.New().String(),
				Time:           time.Now().UnixNano(),
				SequenceNumber: sequenceNumber,
				Thing:          &thing,
			}

			// marshal the Action into a byte array
			if data, err := proto.Marshal(&action); err == nil {
				// publish a message
				if token := p.mqttClient.Publish("protocol", 2, false, data); token.Wait() && token.Error() != nil {
					// bail if there's an error
					return token.Error()
				}

				// increment the counter
				sentCountSecond++
			} else {
				// bail if there's an error
				fmt.Printf("producer: failed to marshal message: [%s] (error: %s)", action.String(), err)
				return err
			}
		} else {
			// let's give the client a second to reconnect on its own
			time.Sleep(1 * time.Second)
		}
	}
}
