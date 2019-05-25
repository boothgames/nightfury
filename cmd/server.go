package cmd

import (
	"context"
	"fmt"
	"github.com/boothgames/nightfury/api"
	"github.com/boothgames/nightfury/cmd/cli"
	"github.com/boothgames/nightfury/pkg/db"
	"github.com/boothgames/nightfury/pkg/nightfury"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start nightfury server and broadcast its information using mdns",
	Run: func(cmd *cobra.Command, args []string) {
		router := gin.Default()
		err := db.Initialize(dbPath)
		cli.DieIf(err)

		cli.DieIf(err)
		cleanup := func() {
			clients := nightfury.Clients{}
			repository := db.DefaultRepository()
			clientsFromRepo, err := nightfury.NewClientsFromRepo(repository)
			cli.DieIf(err)

			mapstructure.Decode(clientsFromRepo, &clients)
			cli.Warn("deleting all clients from db")
			clients.Delete(repository)

			cli.Warn("shutting down db")
			err = db.Close()
			cli.DieIf(err)
		}

		api.Bind(router)
		srv := &http.Server{
			Addr:    fmt.Sprintf("%s:%d", bindAddress, bindPort),
			Handler: router,
		}

		go func() {
			// service connections
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				cli.DieIf(err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		cli.Warn("\ngracefully shutting down server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			cli.Success("server shutdown:", err)
		}
		select {
		case <-ctx.Done():
			cleanup()
		}
		cli.Success("done")
	},
}

var bindAddress string
var bindPort int
var dbPath string

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&bindAddress, "bind-address", "", "0.0.0.0", "specify the advertise address to use")
	serverCmd.Flags().IntVarP(&bindPort, "bind-port", "p", 5624, "specify the advertise port to use")
	serverCmd.Flags().StringVarP(&dbPath, "db-path", "", "nightfury.db", "specify the database path where db will be stored")
}
