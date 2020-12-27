package main

import (
	"flag"
	"log"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.com/clemens97/scionbttracker/registry"
	"gitlab.com/clemens97/scionbttracker/registry/inmem"
	"gitlab.com/clemens97/scionbttracker/registry/redis"
	"gitlab.com/clemens97/scionbttracker/server"
)

var (
	flAddr        = flag.String("addr", ":9090", "address of the tracker")
	flDebug       = flag.Bool("debug", false, "enable debug mode for logging")
	flInterval    = flag.Int("interval", 120, "interval for when Peers should poll for new peers")
	flMinInterval = flag.Int("min-interval", 30, "min poll interval for new peers")
	flRedisAddr   = flag.String("redis-addr", "", "address to a redis server for persistent peer data")
	flRedisPass   = flag.String("redis-pass", "", "password to use to connect to the redis server")

	mux sync.Mutex
)

func main() {
	flag.Parse()
	var (
		logger   = logrus.New()
		registry registry.Registry
	)

	if *flDebug {
		logger.Level = logrus.DebugLevel
	}

	if *flRedisAddr != "" {
		registry = redis.New(*flRedisAddr, *flRedisPass)
	} else {
		registry = inmem.New()
	}

	s := server.New(*flInterval, *flMinInterval, registry, logger)
	if err := http.ListenAndServe(*flAddr, s); err != nil {
		log.Fatal(err)
	}
}
