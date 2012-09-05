package statsd

import (
	"bufio"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"time"
)

type Client struct {
	// connection buffer
	buf *bufio.ReadWriter
	// underlying connection
	conn *net.Conn
	// prefix for statsd name
	prefix string
	// write mutex
	mutex sync.Mutex
}

// Close closes the connection and cleans up.
func (s *Client) Close() error {
	// flush any outstanding data
	s.buf.Flush()
	s.buf = nil
	err := (*s.conn).Close()
	return err
}

// Increments a statsd count type.
// stat is a string name for the metric.
// value is the integer value
// rate is the sample rate (0.0 to 1.0)
func (s *Client) Inc(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|c", value)
	return s.submit(stat, dap, rate)
}

// Decrements a statsd count type.
// stat is a string name for the metric.
// value is the integer value.
// rate is the sample rate (0.0 to 1.0).
func (s *Client) Dec(stat string, value int64, rate float32) error {
	return s.Inc(stat, -value, rate)
}

// Submits/Updates a statsd gauge type.
// stat is a string name for the metric.
// value is the integer value.
// rate is the sample rate (0.0 to 1.0).
func (s *Client) Gauge(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|g", value)
	return s.submit(stat, dap, rate)
}

// Submits a statsd timing type.
// stat is a string name for the metric.
// value is the integer value.
// rate is the sample rate (0.0 to 1.0).
func (s *Client) Timing(stat string, delta int64, rate float32) error {
	dap := fmt.Sprintf("%d|ms", delta)
	return s.submit(stat, dap, rate)
}

// Sets/Updates the statsd client prefix
func (s *Client) SetPrefix(prefix string) {
	s.prefix = prefix
}

// submit formats the statsd event data, handles sampling, and prepares it,
// and sends it to the server.
func (s *Client) submit(stat string, value string, rate float32) error {
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

	_, err := s.send([]byte(data))
	if err != nil {
		return err
	}
	return nil
}

// sends the data to the server endpoint over the net.Conn
func (s *Client) send(data []byte) (int, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	n, err := s.buf.Write([]byte(data))
	if err != nil {
		return 0, err
	}
	if n == 0 {
		return n, errors.New("Wrote no bytes")
	}
	err = s.buf.Flush()
	if err != nil {
		return n, err
	}
	return n, nil
}

func newClient(conn *net.Conn, prefix string) *Client {
	buf := bufio.NewReadWriter(bufio.NewReader(*conn), bufio.NewWriter(*conn))
	client := &Client{
		buf:    buf,
		conn:   conn,
		prefix: prefix}
	return client
}

// Returns a pointer to a new Client.
// addr is a string of the format "hostname:port", and must be parsable by
// net.ResolveUDPAddr.
// prefix is the statsd client prefix. Can be "" if no prefix is desired.
func Dial(addr string, prefix string) (*Client, error) {
	conn, err := net.Dial("udp", addr)
	if err != nil {
		return nil, err
	}

	client := newClient(&conn, prefix)
	return client, nil
}

// Returns a pointer to a new Client.
// addr is a string of the format "hostname:port", and must be parsable by
// net.ResolveUDPAddr.
// timeout is the connection timeout. Since statsd is UDP, there is no
// connection, but the timeout applies to name resolution (if relevant).
// prefix is the statsd client prefix. Can be "" if no prefix is desired.
func DialTimeout(addr string, timeout time.Duration, prefix string) (*Client, error) {
	conn, err := net.DialTimeout("udp", addr, timeout)
	if err != nil {
		return nil, err
	}

	client := newClient(&conn, prefix)
	return client, nil
}
