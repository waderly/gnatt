package main

import (
	"flag"
	"os"
	"os/signal"

	G "github.com/alsm/gnatt/gateway/gate"
)

func main() {
	var gateway G.Gateway
	stopsig := registerSignals()
	gatewayconf := setup()

	G.InitLogger(os.Stdout, os.Stderr) // todo: configurable

	if gatewayconf.IsAggregating() {
		G.INFO.Println("GNATT Gateway starting in aggregating mode")
		gateway = initAggregating(gatewayconf, stopsig)
	} else {
		G.INFO.Println("GNATT Gateway starting in transparent mode")
		gateway = initTransparent(gatewayconf, stopsig)
	}

	gateway.Start()
}

func setup() *G.GatewayConfig {
	var configFile string
	var port int

	flag.StringVar(&configFile, "c", "", "Configuration File")
	flag.IntVar(&port, "port", 0, "MQTT-G UDP Listening Port")
	flag.Parse()

	if configFile != "" {
		if gc, err := G.ParseConfigFile(configFile); err != nil {
			G.ERROR.Fatal(err)
		} else {
			return gc
		}
	}

	G.ERROR.Fatal("-configuration <file> must be specified")
	return nil
}

func initAggregating(gc *G.GatewayConfig, stopsig chan os.Signal) *G.AggGate {
	ag := G.NewAggGate(gc, stopsig)
	return ag
}

func initTransparent(gc *G.GatewayConfig, stopsig chan os.Signal) *G.TransGate {
	tg := G.NewTransGate(gc, stopsig)
	return tg
}

func registerSignals() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	return c
}
