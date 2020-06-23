package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	dbDSN   string
	address int
)

func main() {
	a := cli.NewApp()
	a.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "db-dsn",
			EnvVars:     []string{"DB_DSN"},
			Destination: &dbDSN,
			Required:    true,
		},
		&cli.IntFlag{
			Name:        "address",
			EnvVars:     []string{"SERVER_ADDRESS"},
			Required:    true,
			Destination: &address,
		},
	}
	a.Action = run

	if err := a.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", address),
		Handler:      http.DefaultServeMux,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	go gracefulServerShutdown(c.Context, srv)
	return srv.ListenAndServe()
}

func gracefulServerShutdown(ctx context.Context, srv *http.Server) {
	<-ctx.Done()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	srv.SetKeepAlivesEnabled(false)
	if err := srv.Shutdown(ctx); err != nil {
		panic(err)
	}
}
