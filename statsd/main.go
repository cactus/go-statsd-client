package statsd

import (
	"net"
	"math/rand"
	"fmt"
	"errors"
)

type StatsdClient struct {
	// connection of type interface net.Conn
	conn net.Conn
	// prefix for statsd name
	prefix string
}


// Increments a statsd count type
// stat is a string name for the metric
// value is the integer value
// rate is the sample rate (0.0 to 1.0)
func (s *StatsdClient) Inc(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|c", value)
	return s.submit(stat, dap, rate)
}

// Decrements a statsd count type
// stat is a string name for the metric
// value is the integer value
// rate is the sample rate (0.0 to 1.0)
func (s *StatsdClient) Dec(stat string, value int64, rate float32) error {
	return s.Inc(stat, -value, rate)
}

// Submits/Updates a statsd guage type
// stat is a string name for the metric
// value is the integer value
// rate is the sample rate (0.0 to 1.0)
func (s *StatsdClient) Guage(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|g", value)
	return s.submit(stat, dap, rate)
}

// Submits a statsd timing type
// stat is a string name for the metric
// value is the integer value
// rate is the sample rate (0.0 to 1.0)
func (s *StatsdClient) Timing(stat string, delta int64, rate float32) error {
	dap := fmt.Sprintf("%d|ms", delta)
	return s.submit(stat, dap, rate)
}

// Sets/Updates the statsd client prefix
func (s *StatsdClient) SetPrefix(prefix string) {
	s.prefix = prefix
}

// submit formats the statsd event data, handles sampling, and prepares it,
// and sends it to the server.
func (s *StatsdClient) submit(stat string, value string, rate float32) error {
	if rate < 1 {
		if rand.Float32() < rate {
			value = fmt.Sprintf("%s|@%f", value, rate)
		} else {
			return nil
		}
	}

	if s.prefix != "" {
		stat = fmt.Sprintf("%s.%s", s.prefix, stat)
	}

	data := fmt.Sprintf("%s:%s", stat, value)

	i, err := s.send([]byte(data))
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("Wrote no bytes")
	}
	return nil
}

// sends the data to the server endpoint over the net.Conn
func (s *StatsdClient) send(data []byte) (int, error) {
	i, err := s.conn.Write([]byte(data))
	return i, err
}

// Returns a pointer to a new StatsdClient
// addr is a string of the format "hostname:port", and must be parsable by
// net.ResolveUDPAddr.
// prefix is the statsd client prefix. Can be "" if no prefix is desired.
func New(addr string, prefix string) (*StatsdClient, error) {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP(udpaddr.Network(), nil, udpaddr)
	if err != nil {
		return nil, err
	}

	client := &StatsdClient{
		conn: conn,
		prefix: prefix}

	return client, nil
}
