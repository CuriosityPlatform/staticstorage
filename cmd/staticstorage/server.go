package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"staticstorage/pkg"
	appservice "staticstorage/pkg/service"

	"github.com/gorilla/mux"
)

func server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Run http server",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "no-warm-up-cache",
				Usage: "Dont warm up cache before start server",
				Value: false,
			},
		},
		Action: runServer,
	}
}

func runServer(ctx *cli.Context) error {
	config, err := getConfig(ctx)
	if err != nil {
		return err
	}

	service := pkg.Service()

	if !ctx.Bool("no-warm-up-cache") {
		err = service.WarmUpStorage(ctx.Context, config)
		if err != nil {
			return err
		}
	}

	handlers, err := service.CreateHandlers(ctx.Context, config)
	if err != nil {
		return err
	}

	r := mux.NewRouter()
	for _, handler := range handlers {
		r.Handle(handler.Path, createHTTPHandler(handler))
	}

	// Kubernetes probes route
	r.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%s", config.Port),
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		err = srv.ListenAndServe()
	}()

	<-ctx.Context.Done()
	return srv.Shutdown(context.Background())
}

func createHTTPHandler(handler appservice.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, request *http.Request) {
		file, err := os.ReadFile(handler.AssetPath)
		if err != nil {
			_, _ = w.Write([]byte(errors.Wrapf(err, "failed to read asset file").Error()))
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/octet-stream")
		_, _ = w.Write(file)
	})
}
