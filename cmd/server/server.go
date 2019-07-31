package server

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/jjbubudi/protos-go/tides"
	stan "github.com/nats-io/stan.go"
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
			nats, err := stan.Connect(
				natsClusterID,
				natsClientID,
				stan.NatsURL(natsURL),
				stan.SetConnectionLostHandler(func(conn stan.Conn, err error) {
					log.Fatal(err)
				}),
			)
			if err != nil {
				log.Fatal(err)
			}
			defer nats.Close()
			runServer(nats, bindAddress)
		},
	}
	cmd.Flags().String("bind-address", "0.0.0.0:50051", "Address to bind the server to")
	viper.BindPFlags(cmd.Flags())
	return cmd
}

func runServer(nats stan.Conn, bindAddress string) {
	store := &store{}
	store.Subscribe(nats)

	lis, err := net.Listen("tcp", bindAddress)
	if err != nil {
		log.Fatal(err)
	}

	g := grpc.NewServer()
	tides.RegisterTidesServiceServer(g, &server{
		store: store,
	})
	g.Serve(lis)
}

type store struct {
	realtimeTides  tides.TideRecorded
	predictedTides tides.TidePredicted
}

func (s *store) GetRealtimeTides() tides.TideRecorded {
	return s.realtimeTides
}

func (s *store) GetPredictedTides() tides.TidePredicted {
	return s.predictedTides
}

func (s *store) Subscribe(nats stan.Conn) {
	nats.Subscribe("tides.realtime", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.realtimeTides)
	}, stan.StartWithLastReceived())
	nats.Subscribe("tides.predictions", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.predictedTides)
	}, stan.StartWithLastReceived())
}

type server struct {
	store *store
}

func (s *server) RealtimeTides(ctx context.Context, req *tides.RealtimeTidesRequest) (*tides.RealtimeTidesResponse, error) {
	t := s.store.GetRealtimeTides()
	return &tides.RealtimeTidesResponse{
		Time:   t.GetTime(),
		Meters: t.GetMeters(),
	}, nil
}

func (s *server) PredictedTides(ctx context.Context, req *tides.PredictedTidesRequest) (*tides.PredictedTidesResponse, error) {
	t := s.store.GetPredictedTides()

	var data []*tides.PredictedTidesResponse_Prediction
	for _, prediction := range t.GetPredictions() {
		data = append(data, &tides.PredictedTidesResponse_Prediction{
			Time:   prediction.GetTime(),
			Meters: prediction.GetMeters(),
		})
	}

	return &tides.PredictedTidesResponse{
		Time:        t.GetTime(),
		Predictions: data,
	}, nil
}
