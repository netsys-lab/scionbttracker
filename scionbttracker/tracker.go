package main

import (
	"flag"
	"github.com/netsec-ethz/scion-apps/pkg/shttp"
	"github.com/netsys-lab/scionbttracker/registry"
	"github.com/netsys-lab/scionbttracker/registry/inmem"
	"github.com/netsys-lab/scionbttracker/registry/redis"
	"github.com/netsys-lab/scionbttracker/server"
	"github.com/sirupsen/logrus"
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
	logger.Debugf("attempting to listen on - %s", *flAddr)
	if err := shttp.ListenAndServe(*flAddr, s, nil); err != nil {
		log.Fatal(err)
	}
}
