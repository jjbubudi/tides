package internal

import (
	"log"
	"strings"

	"github.com/nats-io/stan.go"
)

// ConnectToNats connects to Nats Streaming server
func ConnectToNats(clusterID string, clientID string, URL string) (stan.Conn, error) {
	return stan.Connect(
		clusterID,
		clientID,
		stan.NatsURL(URL),
		stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
			log.Fatal(err)
		}),
	)
}

// AutoEnv converts an argument name into environment variable style
func AutoEnv(name string) string {
	return strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
}
