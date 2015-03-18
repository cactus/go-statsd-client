package statsd

import (
	"time"
)

type BufferedSender struct {
	flushIntervalBytes  int
	flushIntervalMillis int
	sender              *SimpleSender
	buffer              []byte
	reqs                chan []byte
	shutdown            chan bool
}

func (s *BufferedSender) Send(data []byte) (int, error) {
	s.reqs <- data
	return len(data), nil
}

func (s *BufferedSender) Close() error {
	s.shutdown <- true
	err := s.sender.Close()
	return err
}

func (s *BufferedSender) Start() {
	ticker := time.NewTicker(time.Duration(s.flushIntervalMillis) * time.Millisecond)

	for {
		select {
		case <-ticker.C:
			s.flush()
		case req := <-s.reqs:
			// StatsD supports receiving multiple metrics in a single packet by separating them with a newline.
			s.buffer = append(s.buffer, append(req, []byte("\n")...)...)
			if len(s.buffer) >= s.flushIntervalBytes {
				s.flush()
			}
		case <-s.shutdown:
			break
		}
	}
}

func (s *BufferedSender) flush() (int, error) {
	n, err := s.sender.Send(s.buffer)
	s.buffer = []byte{} // clear the buffer
	return n, err
}

func NewBufferedSender(addr string, flushIntervalBytes, flushIntervalMillis int) (*BufferedSender, error) {
	simpleSender, err := NewSimpleSender(addr)
	if err != nil {
		return nil, err
	}

	sender := &BufferedSender{
		flushIntervalBytes:  flushIntervalBytes,
		flushIntervalMillis: flushIntervalMillis,
		sender:              simpleSender,
		reqs:                make(chan []byte),
		shutdown:            make(chan bool),
	}

	go sender.Start()
	return sender, nil
}

func NewBufferedClient(addr, prefix string, flushIntervalBytes, flushIntervalMillis int) (*Client, error) {
	if flushIntervalBytes <= 0 {
		flushIntervalBytes = 1432 // https://github.com/etsy/statsd/blob/master/docs/metric_types.md#multi-metric-packets
	}
	if flushIntervalMillis <= 0 {
		flushIntervalMillis = 1000
	}
	sender, err := NewBufferedSender(addr, flushIntervalBytes, flushIntervalMillis)
	if err != nil {
		return nil, err
	}

	client := &Client{
		prefix: prefix,
		sender: sender,
	}

	return client, nil
}
