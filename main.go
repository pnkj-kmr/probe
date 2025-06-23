package main

import (
	"context"
	_ "embed"
	"log"
	"log/slog"
	"os"
	"os/signal"
	probe "probe/internal"
	"syscall"
)

//go:embed version.txt
var appVersion string

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token. Example: "Bearer {token}"

// @title Probe
// @version 1.0.0
// @description agent
// @termsOfService https://github.com/pnkj-kmr/probe
// @contact.name Pankaj Kumar
// @contact.url https://github.com/pnkj-kmr
// @contact.email pnkj3092@gmail.com
// @license.name MIT License
// @license.url https://opensource.org/licenses/MIT

func main() {
	log.Println("app version", appVersion)
	signalChan := shutdownSignal()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// probe, err := probe.NewProbeWithContext(ctx, probe.WithICMP())
	probe, err := probe.NewProbeWithContext(ctx)
	log.Println("probe declration -->", probe)
	if err != nil {
		slog.Error("", "err", err)
		log.Fatalln(err)
	}

	go probe.Start()
	<-signalChan
	probe.Stop()
}

func shutdownSignal() <-chan os.Signal {
	// Create a channel to listen for OS signals (like SIGINT, SIGTERM)
	signalChan := make(chan os.Signal, 1)
	// Notify the channel when a termination signal is received (Ctrl+C or SIGTERM)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	return signalChan
}
