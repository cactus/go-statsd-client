package statsd

import (
	"bytes"
	"testing"
	"time"
)

type mockSender struct {
	closeCallCount int
}

func (m *mockSender) Send(data []byte) (int, error) {
	return 0, nil
}

func (m *mockSender) Close() error {
	m.closeCallCount++
	return nil
}

func TestClose(t *testing.T) {
	mockSender := &mockSender{}
	sender := &BufferedSender{
		flushBytes:    512,
		flushInterval: 1 * time.Second,
		sender:        mockSender,
		buffer:        bytes.NewBuffer(make([]byte, 0, 512)),
		shutdown:      make(chan chan error),
	}

	sender.Close()
	if mockSender.closeCallCount != 0 {
		t.Fatalf("expected close to have been called zero times, but got %d", mockSender.closeCallCount)
	}

	sender.Start()
	if sender.running != 1 {
		t.Fatal("sender failed to start")
	}

	sender.Close()
	if mockSender.closeCallCount != 1 {
		t.Fatalf("expected close to have been called once, but got %d", mockSender.closeCallCount)
	}
}

func TestCloseConcurrent(t *testing.T) {
	mockSender := &mockSender{}
	sender := &BufferedSender{
		flushBytes:    512,
		flushInterval: 1 * time.Second,
		sender:        mockSender,
		buffer:        bytes.NewBuffer(make([]byte, 0, 512)),
		shutdown:      make(chan chan error),
	}
	sender.Start()

	const N = 10
	c := make(chan struct{}, N)
	for i := 0; i < N; i++ {
		go func() {
			sender.Close()
			c <- struct{}{}
		}()
	}

	for i := 0; i < N; i++ {
		<-c
	}

	if mockSender.closeCallCount != 1 {
		t.Errorf("expected close to have been called once, but got %d", mockSender.closeCallCount)
	}
}
