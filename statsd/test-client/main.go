package main

import (
	"flag"
	"time"
	"log"
	"github.com/cactus/go-statsd-client/statsd"
)

func main() {

	address := flag.String("address", "127.0.0.1:8125", "Address:port of statsd")
	prefix := flag.String("prefix", "test-client", "Statsd prefix")
	name := flag.String("name", "counter", "stat name")
	rate := flag.Float64("rate", 1.0, "Sample rate")
	statType := flag.String("type", "count", "Stat type to send. Can be timing, count, gauge")
	statValue := flag.Int64("value", 1, "Value to send")
	volume := flag.Int("volume", 1000, "Number of stats to send")
	duration := flag.Duration("duration", 10*time.Second, "How long to spread the volume across. Each second of duration volume/seconds events will be sent.")
	flag.Parse()

	client, err := statsd.Dial(*address, *prefix)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	var stat func(stat string, value int64, rate float32) error
	switch *statType {
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

	pertick := *volume / int(duration.Seconds()) / 10
	// add some extra tiem, because the first tick takes a while
	ender := time.After(*duration + 100 * time.Millisecond)
	c := time.Tick(time.Second/10)
	count := 0
	for {
		select {
		case <-c:
			for x := 0; x < pertick; x++ {
				err := stat(*name, *statValue, float32(*rate))
				if err != nil {
					log.Printf("Got Error: %+v", err)
					break
				}
				count += 1
			}
		case <-ender:
			log.Printf("%d events called", count)
			return
		}
	}
}
