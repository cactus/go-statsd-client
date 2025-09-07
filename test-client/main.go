// Copyright (c) 2012-2016 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/alecthomas/kong"
	"github.com/cactus/go-statsd-client/v6/statsd"
)

func main() {
	// command line flags
	var opts struct {
		HostPort  string        `name:"host" default:"127.0.0.1:8125" help:"host:port of statsd server"`
		Prefix    string        `name:"prefix" default:"test-client" help:"Statsd prefix"`
		StatType  string        `name:"type" enum:"timing,count,gauge" default:"count" help:"stat type to send. Can be one of: timing, count, gauge"`
		StatValue int64         `name:"value" default:"1" help:"Value to send"`
		Name      string        `name:"name" short:"n" default:"counter" help:"stat name"`
		Rate      float32       `name:"rate" short:"r" default:"1.0" help:"sample rate"`
		Volume    int           `name:"count" short:"c" default:"1000" help:"Number of stats to send. Volume."`
		Nil       bool          `name:"nil" help:"Use nil client"`
		Buffered  bool          `name:"buffered" help:"Use a buffered client"`
		Duration  time.Duration `name:"duration" short:"d" default:"10s" help:"How long to spread the volume across. For each second of duration, volume/seconds events will be sent."`
	}

	// parse said flags
	kong.Parse(&opts,
		kong.Name("test-client"),
		kong.UsageOnError(),
	)
	fmt.Printf("%+v\n", opts)

	if opts.Nil && opts.Buffered {
		fmt.Printf("Specifying both nil and buffered together is invalid\n")
		os.Exit(1)
	}

	if opts.Name == "" || statsd.CheckName(opts.Name) != nil {
		fmt.Printf("Stat name contains invalid characters\n")
		os.Exit(1)
	}

	if statsd.CheckName(opts.Prefix) != nil {
		fmt.Printf("Stat prefix contains invalid characters\n")
		os.Exit(1)
	}

	config := &statsd.ClientConfig{
		Address:     opts.HostPort,
		Prefix:      opts.Prefix,
		ResInterval: time.Duration(0),
	}

	var client statsd.Statter = (*statsd.Client)(nil)
	var err error
	if !opts.Nil {
		if opts.Buffered {
			config.UseBuffered = true
			config.FlushInterval = opts.Duration / time.Duration(4)
			config.FlushBytes = 0
		}

		client, err = statsd.NewClientWithConfig(config)
		if err != nil {
			log.Fatal(err)
		}
		defer client.Close()
	}

	var stat func(stat string, value int64, rate float32) error
	switch opts.StatType {
	case "count":
		stat = func(stat string, value int64, rate float32) error {
			return client.Inc(stat, value, rate)
		}
	case "gauge":
		stat = func(stat string, value int64, rate float32) error {
			return client.Gauge(stat, value, rate)
		}
	case "timing":
		stat = func(stat string, value int64, rate float32) error {
			return client.Timing(stat, value, rate)
		}
	default:
		log.Fatal("Unsupported state type")
	}

	log.Printf("sending to %s\n", opts.HostPort)

	pertick := opts.Volume / int(opts.Duration.Seconds()) / 10
	// add some extra time, because the first tick takes a while
	ender := time.After(opts.Duration + 100*time.Millisecond)
	c := time.Tick(time.Second / 10)
	count := 0
	for {
		select {
		case <-c:
			for x := 0; x < pertick; x++ {
				err := stat(opts.Name, opts.StatValue, opts.Rate)
				if err != nil {
					log.Printf("Got Error: %+v\n", err)
					break
				}
				count++
			}
		case <-ender:
			log.Printf("%d events called\n", count)
			return
		}
	}
}
