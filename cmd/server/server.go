package server

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jjbubudi/protos-go/tides"
	stan "github.com/nats-io/go-nats-streaming"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	cmd.Flags().String("bind-address", "0.0.0.0:8080", "Address to bind the server to")
	viper.BindPFlags(cmd.Flags())
	return cmd
}

func runServer(nats stan.Conn, bindAddress string) {
	store := &store{}
	store.Subscribe(nats)

	var hongKong, _ = time.LoadLocation("Asia/Hong_Kong")

	router := gin.New()
	router.Use(gin.Recovery())

	v1 := router.Group("v1")
	v1.GET("/tides/realtime", func(c *gin.Context) {
		r := store.realtimeTides
		c.JSON(http.StatusOK, tidalData{
			Time:   timestampToString(r.GetTime(), hongKong),
			Meters: r.GetMeters(),
		})
	})
	v1.GET("/tides/predicted", func(c *gin.Context) {
		r := store.predictedTides
		var data []tidalData
		for _, prediction := range r.Predictions {
			data = append(data, tidalData{
				Time:   timestampToString(prediction.GetTime(), hongKong),
				Meters: prediction.GetMeters(),
			})
		}
		c.JSON(http.StatusOK, predictions{
			Time:        timestampToString(r.GetTime(), hongKong),
			Predictions: data,
		})
	})

	gin.SetMode(gin.ReleaseMode)
	router.Run(bindAddress)
}

func timestampToString(timestamp *timestamp.Timestamp, location *time.Location) string {
	return time.Unix(timestamp.GetSeconds(), int64(timestamp.GetNanos())).In(location).Format(time.RFC3339Nano)
}

type tidalData struct {
	Time   string  `json:"time"`
	Meters float64 `json:"meters"`
}

type predictions struct {
	Time        string      `json:"time"`
	Predictions []tidalData `json:"predictions"`
}

type store struct {
	realtimeTides  tides.TideRecorded
	predictedTides tides.TidePredicted
}

func (s *store) Subscribe(nats stan.Conn) {
	nats.Subscribe("tides", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.realtimeTides)
	}, stan.StartWithLastReceived())
	nats.Subscribe("tides_predictions", func(m *stan.Msg) {
		proto.Unmarshal(m.Data, &s.predictedTides)
	}, stan.StartWithLastReceived())
}
