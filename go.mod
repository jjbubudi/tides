module github.com/jjbubudi/tides

go 1.12

require (
	github.com/gin-gonic/gin v1.4.0
	github.com/golang/protobuf v1.3.1
	github.com/jasonlvhit/gocron v0.0.0-20190402024347-5bcdd9fcfa9b
	github.com/jjbubudi/protos-go v0.0.0-20190513145009-4dcc2b0121af
	github.com/nats-io/go-nats-streaming v0.4.2
	github.com/spf13/cobra v0.0.4-0.20190321000552-67fc4837d267
	github.com/spf13/viper v1.3.2
	github.com/stretchr/testify v1.3.0
)

replace (
	github.com/nats-io/go-nats-streaming v0.4.2 => github.com/jjbubudi/go-nats-streaming v0.4.3-0.20190420024036-3a359ddc011c
	github.com/ugorji/go v1.1.4 => github.com/ugorji/go/codec v0.0.0-20190204201341-e444a5086c43
)
