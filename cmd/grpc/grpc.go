package grpc

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/jjbubudi/protos-go/tides"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

// NewServerCommand creates a new server command
func NewServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Start the server",
		Long:  "Start the server",
		Run: func(cmd *cobra.Command, args []string) {
			natsURL := viper.GetString("nats-url")
			natsClusterID := viper.GetString("nats-cluster-id")
			natsClientID := viper.GetString("nats-client-id")
			bindAddress := viper.GetString("bind-address")
			nats, err := stan.Connect(natsClusterID, natsClientID, stan.NatsURL(natsURL))
			if err != nil {
				log.Fatal(err)
			}
			defer nats.Close()
			runServer(nats, bindAddress)
		},
	}
	cmd.Flags().String("bind-address", "0.0.0.0:50051", "Address to bind the gRPC server to")
	viper.BindPFlags(cmd.Flags())
	return cmd
}

func runServer(nats stan.Conn, bindAddress string) {
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := &server{}
	server.Subscribe(nats)

	s := grpc.NewServer()
	tides.RegisterTidesServiceServer(s, server)
	if err := s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type server struct {
	realtimeTides  tides.TideRecorded
	predictedTides tides.TidePredicted
}

func (s *server) Subscribe(nats stan.Conn) {
	nats.Subscribe("tides", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.realtimeTides)
	}, stan.StartWithLastReceived())
	nats.Subscribe("tides_predictions", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.predictedTides)
	}, stan.StartWithLastReceived())
}

func (s *server) RealtimeTides(ctx context.Context, req *tides.RealtimeTidesRequest) (*tides.RealtimeTidesResponse, error) {
	return &tides.RealtimeTidesResponse{
		Time:   s.realtimeTides.Time,
		Meters: s.realtimeTides.Meters,
	}, nil
}

func (s *server) PredictedTides(ctx context.Context, req *tides.PredictedTidesRequest) (*tides.PredictedTidesResponse, error) {
	var predictions []*tides.PredictedTidesResponse_Prediction
	for _, t := range s.predictedTides.Predictions {
		predictions = append(predictions, &tides.PredictedTidesResponse_Prediction{
			Time:   t.Time,
			Meters: t.Meters,
		})
	}
	return &tides.PredictedTidesResponse{
		Predictions: predictions,
	}, nil
}
