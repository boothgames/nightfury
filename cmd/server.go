package cmd

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jskswamy/herman/cmd/cli"
	"github.com/spf13/cobra"
	"gitlab.com/jskswamy/nightfury/api"
	"gitlab.com/jskswamy/nightfury/pkg/db"
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
		cli.Warn("\nshutdown server ...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
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
