package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jasonlvhit/gocron"
	"github.com/jjbubudi/tides/internal/tides"
	"github.com/jjbubudi/tides/pkg/observatory"
	stan "github.com/nats-io/go-nats-streaming"
)

var hongKong, _ = time.LoadLocation("Asia/Hong_Kong")

func main() {
	natsCluster := os.Getenv("NATS_CLUSTER")
	natsURL := os.Getenv("NATS_URL")
	natsClientID := os.Getenv("NATS_CLIENT_ID")

	nats, err := stan.Connect(natsCluster, natsClientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf(`Error connecting to Nats! Cluster: "%s", URL: "%s", ClientID: "%s"`,
			natsCluster, natsURL, natsClientID)
	}
	defer nats.Close()

	observatory := observatory.New(&http.Client{})
	tidesService := tides.New(observatory.TidalDataAsOf, observatory.TidalPredictionsAsOf, nats.Publish)

	gocron.ChangeLoc(hongKong)

	gocron.Every(5).Minutes().Do(func() {
		err := tidesService.PublishRealTimeTidalData(time.Now())
		if err != nil {
			log.Println(err)
		}
	})

	gocron.Every(1).Day().At("01:00").Do(func() {
		err := tidesService.PublishTidalPredictions(time.Now())
		if err != nil {
			log.Println(err)
		}
	})

	<-gocron.Start()
}
