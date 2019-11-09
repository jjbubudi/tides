package fetch

import (
	"net/http"
	"time"

	"github.com/jjbubudi/tides/internal"
	"github.com/jjbubudi/tides/internal/observatory"
	"github.com/jjbubudi/tides/internal/tides"
	stan "github.com/nats-io/stan.go"
	"github.com/urfave/cli"
)

const (
	pNatsURL       = "nats-url"
	pNatsClusterID = "nats-cluster-id"
	pNatsClientID  = "nats-client-id"
)

// NewFetchCommand creates a new fetch command
func NewFetchCommand() cli.Command {
	return cli.Command{
		Name:  "fetch",
		Usage: "Fetches tidal data",
		Subcommands: []cli.Command{
			cli.Command{
				Name:  "realtime",
				Usage: "Realtime tidal data",
				Flags: publisherFlags(),
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
					return publish(nats, true)
				},
			},
			cli.Command{
				Name:  "predicted",
				Usage: "Predicted tidal data for today",
				Flags: publisherFlags(),
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
					return publish(nats, false)
				},
			},
		},
	}
}

func publisherFlags() []cli.Flag {
	return []cli.Flag{
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
	}
}

func publish(nats stan.Conn, realtime bool) error {
	hongKong, _ := time.LoadLocation("Asia/Hong_Kong")
	observatory := observatory.New(&http.Client{})
	tidesService := tides.New(observatory.TidalDataAsOf, observatory.TidalPredictionsAsOf, nats.Publish)
	if realtime {
		return tidesService.PublishRealTimeTidalData(time.Now().In(hongKong))
	}
	return tidesService.PublishTidalPredictions(time.Now().In(hongKong))
}
