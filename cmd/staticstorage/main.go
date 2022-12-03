package main

import (
	"bytes"
	"context"
	stdlog "log"
	"os"
	"os/signal"
	"path"
	"path/filepath"
	"syscall"

	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"

	"staticstorage/pkg/config"
)

const (
	appID = "staticstorage"
)

var (
	version = "v1"
)

func main() {
	ctx := context.Background()

	ctx = subscribeForKillSignals(ctx)

	err := runApp(ctx, os.Args)
	if err != nil {
		stdlog.Fatal(err)
	}
}

func runApp(ctx context.Context, args []string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	currentDir, err := getCurrentDir()
	if err != nil {
		return err
	}

	serverCmd := server()

	app := &cli.App{
		Name:        appID,
		Version:     version,
		Description: "Serve external assets without keeping them locally or in repository",
		Usage:       "Simple service to serve external assets",
		Action:      serverCmd.Action,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "Path to config file",
				Value:   path.Join(currentDir, "config.json"),
			},
		},
		Commands: []*cli.Command{
			serverCmd,
			warmUpCache(),
		},
	}

	return app.RunContext(ctx, args)
}

func subscribeForKillSignals(ctx context.Context) context.Context {
	ctx, cancel := context.WithCancel(ctx)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		defer cancel()
		select {
		case <-ctx.Done():
			signal.Stop(ch)
		case <-ch:
		}
	}()

	return ctx
}

func getCurrentDir() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get current dir")
	}
	return filepath.Dir(ex), nil
}

func getConfig(ctx *cli.Context) (config.Config, error) {
	configPath := ctx.String("config")
	if configPath == "" {
		return config.Config{}, errors.New("no config path passed")
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		return config.Config{}, errors.Wrapf(err, "failed to read config at path %s", configPath)
	}

	return config.Parser{}.ParseConfig(bytes.NewReader(configBytes))
}
