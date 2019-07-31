package publisher

import (
	"log"
	"net/http"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jjbubudi/tides/internal/observatory"
	"github.com/jjbubudi/tides/internal/tides"
	stan "github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// NewPublisherCommand creates a new publisher command
func NewPublisherCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "publisher",
		Short: "Start the publisher",
		Long:  "Start the publisher",
		Run: func(cmd *cobra.Command, args []string) {
			natsURL := viper.GetString("nats-url")
			natsClusterID := viper.GetString("nats-cluster-id")
			natsClientID := viper.GetString("nats-client-id")
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
			runPublisher(nats)
		},
	}
}

func runPublisher(nats stan.Conn) {
	hongKong, _ := time.LoadLocation("Asia/Hong_Kong")
	observatory := observatory.New(&http.Client{})
	tidesService := tides.New(observatory.TidalDataAsOf, observatory.TidalPredictionsAsOf, nats.Publish)

	gocron.ChangeLoc(hongKong)

	gocron.Every(2).Minutes().Do(func() {
		err := tidesService.PublishRealTimeTidalData(time.Now().In(hongKong))
		if err != nil {
			log.Println(err)
		}
	})

	gocron.Every(1).Day().At("01:00").Do(func() {
		err := tidesService.PublishTidalPredictions(time.Now().In(hongKong))
		if err != nil {
			log.Println(err)
		}
	})

	<-gocron.Start()
}
