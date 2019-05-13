module github.com/jjbubudi/tides

go 1.12

require (
	github.com/golang/protobuf v1.3.1
	github.com/jasonlvhit/gocron v0.0.0-20190402024347-5bcdd9fcfa9b
	github.com/jjbubudi/protos-go v0.0.0-20190513145009-4dcc2b0121af
	github.com/nats-io/go-nats-streaming v0.4.2
	github.com/stretchr/testify v1.3.0
)

replace github.com/nats-io/go-nats-streaming v0.4.2 => github.com/jjbubudi/go-nats-streaming v0.4.3-0.20190420024036-3a359ddc011c
