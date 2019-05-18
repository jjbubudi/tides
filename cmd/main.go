package main

import (
	"log"
	"strings"

	"github.com/jjbubudi/tides/cmd/grpc"
	"github.com/jjbubudi/tides/cmd/publisher"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	rootCommand := &cobra.Command{
		Use: "tides",
	}
	rootCommand.PersistentFlags().String("nats-url", "nats://127.0.0.1:4222", "Nats URL to connect to")
	rootCommand.PersistentFlags().String("nats-cluster-id", "test-cluster", "Nats cluster ID")
	rootCommand.PersistentFlags().String("nats-client-id", "test-client", "Nats client ID")
	viper.BindPFlags(rootCommand.PersistentFlags())
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))

	rootCommand.AddCommand(publisher.NewPublisherCommand())
	rootCommand.AddCommand(grpc.NewServerCommand())

	if err := rootCommand.Execute(); err != nil {
		log.Fatal(err)
	}
}
