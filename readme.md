# golang protobuf stuff

this is a simple example of how to use protobuf with golang. in this case, we're using protobuf to serialize a message and send it over mqtt.

the programs assume you have an mqtt broker on localhost which does not require authentication or encryption. if you don't, you should check out the mosquitto broker.

## building

```bash
$ make
Cleaning up...
rm -rf bin internal/protocol
Building the protocol...
mkdir -p internal/protocol
protoc --go_out=paths=source_relative:internal/protocol protocol.proto
Building main program...
mkdir -p bin
go mod tidy
go build -o bin/gpbs ./
```

## running

the program is split into two parts: a producer and a consumer. the producer generates a protobuf message and sends it over mqtt. the consumer receives the message and prints details about it.

start the consumer:

```bash
$ ./bin/gpbs consumer
messages received per second: 0
messages received per second: 0
messages received per second: 0
messages received per second: 12478
messages received per second: 19036
messages received per second: 17824
messages received per second: 17794
...
```

start the producer:

```bash
$ ./bin/gpbs producer
messages sent per second: 18127
messages sent per second: 18694
messages sent per second: 18011
messages sent per second: 17614
...
```
