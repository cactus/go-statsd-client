package statsd

import (
	"net"
	"math/rand"
	"fmt"
	"errors"
)

type StatsdClient struct {
	addr *net.UDPAddr
	conn *net.UDPConn
	prefx string
}

func (s *StatsdClient) Inc(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|c", value)
	return s.submit(stat, dap, rate)
}

func (s *StatsdClient) Dec(stat string, value int64, rate float32) error {
	return s.Inc(stat, -value, rate)
}

func (s *StatsdClient) Guage(stat string, value int64, rate float32) error {
	dap := fmt.Sprintf("%d|g", value)
	return s.submit(stat, dap, rate)
}

func (s *StatsdClient) Timing(stat string, delta int64, rate float32) error {
	dap := fmt.Sprintf("%d|ms", delta)
	return s.submit(stat, dap, rate)
}

func (s *StatsdClient) submit(stat string, value string, rate float32) error {
	if rate < 1 {
		if rand.Float32() < rate {
			value = fmt.Sprintf("%s|@%4.3", value, rate)
		} else {
			return nil
		}
	}

	if s.prefx != "" {
		stat = fmt.Sprintf("%s.%s", s.prefx, stat)
	}

	data := fmt.Sprintf("%s:%s", stat, value)

	i, err := s.conn.Write([]byte(data))
	if err != nil {
		return err
	}
	if i == 0 {
		return errors.New("Wrote no bytes")
	}
	return nil
}

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
		addr: udpaddr,
		conn: conn,
		prefx: prefix}

	return client, nil
}
