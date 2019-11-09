package server

import (
	"context"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/jjbubudi/protos-go/tides"
	"github.com/jjbubudi/tides/internal"
	stan "github.com/nats-io/stan.go"
	"github.com/urfave/cli"
	"google.golang.org/grpc"
)

const (
	pNatsURL       = "nats-url"
	pNatsClusterID = "nats-cluster-id"
	pNatsClientID  = "nats-client-id"
	pBindAddress   = "bind-address"
)

// NewServerCommand creates a new server command
func NewServerCommand() cli.Command {
	return cli.Command{
		Name:  "server",
		Usage: "Server serving tidal data queries",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "start",
				Usage: "Start the server",
				Flags: []cli.Flag{
					cli.StringFlag{
						Name:     pNatsURL,
						EnvVar:   internal.AutoEnv(pNatsURL),
						Required: true,
					},
					cli.StringFlag{
						Name:     pNatsClusterID,
						EnvVar:   internal.AutoEnv(pNatsClusterID),
						Required: true,
					},
					cli.StringFlag{
						Name:     pNatsClientID,
						EnvVar:   internal.AutoEnv(pNatsClientID),
						Required: true,
					},
					cli.StringFlag{
						Name:     pBindAddress,
						EnvVar:   internal.AutoEnv(pBindAddress),
						Value:    "0.0.0.0:50051",
						Required: false,
					},
				},
				Action: func(c *cli.Context) error {
					nats, err := internal.ConnectToNats(
						c.String(pNatsClusterID),
						c.String(pNatsClientID),
						c.String(pNatsURL),
					)
					if err != nil {
						return err
					}
					defer nats.Close()
					return startServer(nats, c.String(pBindAddress))
				},
			},
		},
	}
}

func startServer(nats stan.Conn, bindAddress string) error {
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

	return g.Serve(lis)
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
