package main

import (
	"flag"
	"github.com/netsec-ethz/scion-apps/pkg/shttp"
	"github.com/sirupsen/logrus"
	"gitlab.com/clemens97/scionbttracker/registry"
	"gitlab.com/clemens97/scionbttracker/registry/inmem"
	"gitlab.com/clemens97/scionbttracker/registry/redis"
	"gitlab.com/clemens97/scionbttracker/server"
	"log"
)

var (
	flAddr        = flag.String("addr", ":9090", "address of the tracker")
	flDebug       = flag.Bool("debug", false, "enable debug mode for logging")
	flInterval    = flag.Int("interval", 120, "interval for when Peers should poll for new peers")
	flMinInterval = flag.Int("min-interval", 30, "min poll interval for new peers")
	flRedisAddr   = flag.String("redis-addr", "", "address to a redis server for persistent peer data")
	flRedisPass   = flag.String("redis-pass", "", "password to use to connect to the redis server")
)

func main() {
	flag.Parse()
	var (
		logger = logrus.New()
		r      registry.Registry
	)

	if *flDebug {
		logger.Level = logrus.DebugLevel
	}

	if *flRedisAddr != "" {
		r = redis.New(*flRedisAddr, *flRedisPass)
	} else {
		r = inmem.New()
	}

	s := server.New(*flInterval, *flMinInterval, r, logger)
	if err := shttp.ListenAndServe(*flAddr, s, nil); err != nil {
		log.Fatal(err)
	}
}
