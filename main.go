package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/n3k0fi5t/wallet/app/api"
	"github.com/sirupsen/logrus"
)

var (
	wait    = flag.Duration("GRACEFULL_TIMEOUT", 15*time.Second, "the duration for which the server gracefully wait for existing connections to finish")
	apiPort = os.Getenv("API_PORT")
)

func main() {
	flag.Parse()

	rt := api.BuildRouter()

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%s", apiPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      rt,
	}

	// Run our server in a goroutine so that main won't be blocked.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			logrus.WithField("err", err).Fatal("ListenAndServe failed")
			panic(err)
		}
	}()

	// Graceful Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	<-quit

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), *wait)
	defer cancel()

	srv.Shutdown(ctx)
	logrus.Info("Service shutdown")

	os.Exit(0)
}
