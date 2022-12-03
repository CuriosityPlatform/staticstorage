package main

import (
	"github.com/urfave/cli/v2"

	"staticstorage/pkg"
)

func warmUpCache() *cli.Command {
	return &cli.Command{
		Name:   "warm-cache",
		Usage:  "Warm up cache without running server",
		Action: executeWarmUpCache,
	}
}

func executeWarmUpCache(ctx *cli.Context) error {
	config, err := getConfig(ctx)
	if err != nil {
		return err
	}

	service := pkg.Service()

	return service.WarmUpStorage(ctx.Context, config)
}
