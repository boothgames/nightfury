package cmd

import (
	"fmt"
	"github.com/boothgames/nightfury/api"
	"github.com/boothgames/nightfury/api/metrics"
	"github.com/boothgames/nightfury/cmd/cli"
	"github.com/boothgames/nightfury/log"
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start nightfury server",
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLogLevel(logLevel)
		go startServer()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		cli.Warn("\ngracefully shutting down server ...")
		cleanup()
		cli.Success("done")
	},
}

var (
	bindAddress            string
	bindPort               int
	dbPath                 string
	logLevel               string
	metricsServer          string
	metricsBucket          string
	metricsServerAuthToken string
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&bindAddress, "bind-address", "", "0.0.0.0", "specify the advertise address to use")
	serverCmd.Flags().IntVarP(&bindPort, "bind-port", "p", 5624, "specify the advertise port to use")
	serverCmd.Flags().StringVarP(&logLevel, "log-level", "l", "error", "specify the log level (panic, fatal, error, warn, info, debug, trace)")
	serverCmd.Flags().StringVarP(&dbPath, "db-path", "", "nightfury.db", "specify the database path where db will be stored")
	serverCmd.Flags().StringVarP(&metricsServer, "metrics-server", "m", "http://localhost:8086", "specify the URI for influx db to store application metrics")
	serverCmd.Flags().StringVarP(&metricsBucket, "metrics-bucket", "b", "nightfury", "bucket name to be used for storing the metrics information")
	serverCmd.Flags().StringVarP(&metricsServerAuthToken, "metrics-server-auth-token", "", "", "metrics server auth token")
}

func releaseMode() string {
	switch strings.ToLower(logLevel) {
	case "panic", "fatal", "error":
		return gin.ReleaseMode
	default:
		return gin.DebugMode
	}
}

func startServer() {
	address := fmt.Sprintf("%s:%d", bindAddress, bindPort)
	cli.Info(fmt.Sprintf("starting nightfury at %s", address))
	cli.Info(fmt.Sprintf("using metrics server at %s with bucket %s", metricsServer, metricsBucket))

	gin.SetMode(releaseMode())
	router := gin.Default()

	err := db.Initialize(dbPath)
	cli.DieIf(err)

	api.Bind(router)
	err = metrics.Initialize(metricsServer, metricsBucket, metricsServerAuthToken)
	cli.DieIf(err)

	srv := &http.Server{Addr: address, Handler: router}
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		cli.DieIf(err)
	}
}

func cleanup() {
	clients := nightfury.Clients{}
	repository := db.DefaultRepository()
	clientsFromRepo, err := nightfury.NewClientsFromRepo(repository)
	cli.DieIf(err)

	err = mapstructure.Decode(clientsFromRepo, &clients)
	cli.DieIf(err)

	cli.Warn("deleting all clients from db")

	err = clients.Delete(repository)
	cli.DieIf(err)

	cli.Warn("shutting down db")
	err = db.Close()
	cli.DieIf(err)
}
