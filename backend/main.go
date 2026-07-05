package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/danielgtaylor/huma/v2/humacli"
	"github.com/spf13/cobra"

	"github.com/meshbox-studio/meshbox-studio/internal/server"
	"github.com/meshbox-studio/meshbox-studio/internal/stats"
)

type Options struct {
	Port    int    `doc:"Port to listen on." short:"p" default:"8080"`
	Host    string `doc:"Host to listen on." default:"0.0.0.0"`
	DataDir string `doc:"Directory for persistent project data." default:"./data"`
}

func main() {
	cli := humacli.New(func(hooks humacli.Hooks, opts *Options) {
		opts.loadFromEnv()

		_, mux := server.NewAPI(os.Getenv("GO_ENV") == "development", opts.DataDir)

		addr := fmt.Sprintf("%s:%d", opts.Host, opts.Port)
		srv := &http.Server{Addr: addr, Handler: mux}

		hooks.OnStart(func() {
			fmt.Printf("Server is running on http://%s\n", addr)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				panic(err)
			}
		})

		hooks.OnStop(func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			srv.Shutdown(ctx)
		})
	})

	cli.Root().AddCommand(&cobra.Command{
		Use:   "openapi",
		Short: "Print the OpenAPI spec to stdout",
		Run: func(cmd *cobra.Command, args []string) {
			api, _ := server.NewAPI(false, "./data")
			b, err := api.OpenAPI().YAML()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(b))
		},
	})

	cmd := cli.Root()
	cmd.Use = "meshbox-studio"
	cmd.Version = stats.AppVersion()

	cli.Run()
}

func (o *Options) loadFromEnv() {
	if v := os.Getenv("PORT"); v != "" {
		fmt.Sscanf(v, "%d", &o.Port)
	}
	if v := os.Getenv("HOST"); v != "" {
		o.Host = v
	}
	if v := os.Getenv("DATA_DIR"); v != "" {
		o.DataDir = v
	}
}

func init() {
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)
}