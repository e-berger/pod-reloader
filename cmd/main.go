package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/e-berger/pod-reloader/internal/config"
	"github.com/e-berger/pod-reloader/internal/process"
	"github.com/e-berger/pod-reloader/internal/registry"
)

const (
	FREQUENCY_CHECK_SECONDS = 30
)

func main() {
	config.SetLogger()
	registryConfig, err := config.GetRegistryConfig()
	if err != nil {
		slog.Error("Error getting registry config", "error", err)
		panic(err)
	}

	p, err := process.NewProcess(registryConfig)

	switch registryConfig.GetType() {
	case registry.ECR:
		cfg, err := config.CreateIAMConfig()
		if err != nil {
			slog.Error("Error creating IAM config", "error", err)
			panic(err)
		}
		p.Registry.(*registry.RegistryEcr).SetAwsConfig(*cfg)

	case registry.DOCKER:
		auths, err := config.GetRegistryAuth(p)
		if err != nil {
			slog.Error("Error getting docker registry auth", "error", err)
			panic(err)
		}
		if auths != nil {
			p.Registry.(*registry.RegistryDocker).SetAuths(auths)
		}
	}

	frequency := FREQUENCY_CHECK_SECONDS
	frequencyEnv := os.Getenv("FREQUENCY_CHECK_SECONDS")
	if frequencyEnv != "" {
		frequencyInt64, _ := strconv.ParseInt(frequencyEnv, 10, 64)
		frequency = int(frequencyInt64)
	}

	ticker := time.NewTicker(time.Second * time.Duration(frequency))
	done := make(chan bool)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				err = p.Tick()
				if err != nil {
					slog.Error("Error during main loop", "error", err)
				}
			}
		}
	}()
	<-ctx.Done()
	stop()
	done <- true
}
