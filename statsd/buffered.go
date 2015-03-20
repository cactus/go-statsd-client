package statsd

import (
	"bytes"
	"time"
)

// BufferedSender provides a buffered statsd udp, sending multiple
// metrics, where possible.
type BufferedSender struct {
	flushBytes    int
	flushInterval time.Duration
	sender        Sender
	buffer        *bytes.Buffer
	reqs          chan []byte
	shutdown      chan bool
}

// Send bytes
func (s *BufferedSender) Send(data []byte) (int, error) {
	s.reqs <- data
	return len(data), nil
}

// Close Buffered Sender
func (s *BufferedSender) Close() error {
	s.shutdown <- true
	err := s.sender.Close()
	return err
}

// Start Buffered Sender
// Begins ticker and read loop
func (s *BufferedSender) Start() {
	ticker := time.NewTicker(s.flushInterval)

	for {
		select {
		case <-ticker.C:
			if s.buffer.Len() > 0 {
				s.flush()
			}
		case req := <-s.reqs:
			// StatsD supports receiving multiple metrics in a single packet by
			// separating them with a newline.
			newLine := append(req, '\n')
			if s.buffer.Len()+len(newLine) > s.flushBytes {
				s.flush()
			}
			s.buffer.Write(newLine)
			if s.buffer.Len() >= s.flushBytes {
				s.flush()
			}
		case <-s.shutdown:
			break
		}
	}
}

// flush the buffer/send to remove endpoint.
func (s *BufferedSender) flush() (int, error) {
	n, err := s.sender.Send(s.buffer.Bytes())
	s.buffer.Reset() // clear the buffer
	return n, err
}

// Returns a new BufferedSender
//
// addr is a string of the format "hostname:port", and must be parsable by
// net.ResolveUDPAddr.
//
// flushInterval is a time.Duration, and specifies the maximum interval for
// packet sending. Note that if you send lots of metrics, you will send more
// often. This is just a maximal threshold.
//
// flushBytes specifies the maximum udp packet size you wish to send. If adding
// a metric would result in a larger packet than flushBytes, the packet will
// first be send, then the new data will be added to the next packet.
func NewBufferedSender(addr string, flushInterval time.Duration, flushBytes int) (Sender, error) {
	simpleSender, err := NewSimpleSender(addr)
	if err != nil {
		return nil, err
	}

	sender := &BufferedSender{
		flushBytes:    flushBytes,
		flushInterval: flushInterval,
		sender:        simpleSender,
		buffer:        bytes.NewBuffer(make([]byte, 0, flushBytes)),
		reqs:          make(chan []byte),
		shutdown:      make(chan bool),
	}

	go sender.Start()
	return sender, nil
}

// Return a new BufferedClient
//
// addr is a string of the format "hostname:port", and must be parsable by
// net.ResolveUDPAddr.
//
// prefix is the statsd client prefix. Can be "" if no prefix is desired.
//
// flushInterval is a time.Duration, and specifies the maximum interval for
// packet sending. Note that if you send lots of metrics, you will send more
// often. This is just a maximal threshold.
//
// flushBytes specifies the maximum udp packet size you wish to send. If adding
// a metric would result in a larger packet than flushBytes, the packet will
// first be send, then the new data will be added to the next packet.
//
// If flushBytes is 0, defaults to 1432 bytes, which is considered safe
// for local traffic. If sending over the public internet, 512 bytes is
// the recommended value.
func NewBufferedClient(addr, prefix string, flushInterval time.Duration, flushBytes int) (Statter, error) {
	if flushBytes <= 0 {
		// https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
		flushBytes = 1432
	}
	if flushInterval <= time.Duration(0) {
		flushInterval = 300 * time.Millisecond
	}
	sender, err := NewBufferedSender(addr, flushInterval, flushBytes)
	if err != nil {
		return nil, err
	}

	client := &Client{
		prefix: prefix,
		sender: sender,
	}

	return client, nil
}
