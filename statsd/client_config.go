// Copyright (c) 2012-2016 Eli Janssen
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package statsd

import (
	"fmt"
	"time"
)

type ClientConfig struct {
	// addr is a string of the format "hostname:port", and must be something
	// validly parsable by net.ResolveUDPAddr.
	Address string

	// prefix is the statsd client prefix. Can be "" if no prefix is desired.
	Prefix string

	// ResInterval is the interval over which the addr is re-resolved.
	// Do note that this /does/ add overhead!
	// If you need higher performance, leave unset (or set to 0),
	// in which case the address will not be re-resolved.
	//
	// Note that if Address is an {ip}:{port} and not a {hostname}:{port}, then
	// ResInterval will be ignored.
	ResInterval time.Duration

	// UseBuffered determines whether a buffered sender is used or not.
	// If a buffered sender is /not/ used, FlushInterval and FlushBytes values are
	// ignored. Default is false.
	UseBuffered bool

	// FlushInterval is a time.Duration, and specifies the maximum interval for
	// packet sending. Note that if you send lots of metrics, you will send more
	// often. This is just a maximal threshold.
	// If FlushInterval is 0, defaults to 300ms.
	FlushInterval time.Duration

	// If flushBytes is 0, defaults to 1432 bytes, which is considered safe
	// for local traffic. If sending over the public internet, 512 bytes is
	// the recommended value.
	FlushBytes int
}

// NewClientWithConfig returns a new BufferedClient
//
// config is a ClientConfig, which holds various configuration values.
func NewClientWithConfig(config *ClientConfig) (Statter, error) {
	var sender Sender
	var err error

	// guard against nil config
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// Use a re-resolving simple sender iff:
	// *  The time duration greater than 0
	// *  The Address is not an ip (eg. {ip}:{port}).
	// Otherwise, re-resolution is not required.
	if config.ResInterval > 0 && !mustBeIP(config.Address) {
		sender, err = NewResolvingSimpleSender(config.Address, config.ResInterval)
	} else {
		sender, err = NewSimpleSender(config.Address)
	}
	if err != nil {
		return nil, err
	}

	if config.UseBuffered {
		return newBufferedC(sender, config)
	} else {
		return NewClientWithSender(sender, config.Prefix)
	}
}

func newBufferedC(baseSender Sender, config *ClientConfig) (Statter, error) {

	flushBytes := config.FlushBytes
	if flushBytes <= 0 {
		// ref:
		// github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
		flushBytes = 1432
	}

	flushInterval := config.FlushInterval
	if flushInterval <= time.Duration(0) {
		flushInterval = 300 * time.Millisecond
	}

	bufsender, err := NewBufferedSenderWithSender(baseSender, flushInterval, flushBytes)
	if err != nil {
		return nil, err
	}

	return NewClientWithSender(bufsender, config.Prefix)
}

// NewClientWithSender returns a pointer to a new Client and an error.
//
// sender is an instance of a statsd.Sender interface and may not be nil
//
// prefix is the stastd client prefix. Can be "" if no prefix is desired.
func NewClientWithSender(sender Sender, prefix string) (Statter, error) {
	if sender == nil {
		return nil, fmt.Errorf("Client sender may not be nil")
	}

	return &Client{prefix: prefix, sender: sender}, nil
}
