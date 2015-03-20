package statsd

import (
	"bytes"
	"log"
	"reflect"
	"testing"
	"time"
)

func TestBufferedClient(t *testing.T) {
	l, err := newUDPListener("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	for _, tt := range statsdPacketTests {
		c, err := NewBufferedClient(l.LocalAddr().String(), tt.Prefix, 10*time.Millisecond, 100)
		if err != nil {
			c.Close()
			t.Fatal(err)
		}
		method := reflect.ValueOf(c).MethodByName(tt.Method)
		e := method.Call([]reflect.Value{
			reflect.ValueOf(tt.Stat),
			reflect.ValueOf(tt.Value),
			reflect.ValueOf(tt.Rate)})[0]
		errInter := e.Interface()
		if errInter != nil {
			c.Close()
			t.Fatal(errInter.(error))
		}

		data := make([]byte, 128)
		_, _, err = l.ReadFrom(data)
		if err != nil {
			c.Close()
			t.Fatal(err)
		}

		data = bytes.TrimRight(data, "\x00\n")
		if bytes.Equal(data, []byte(tt.Expected)) != true {
			t.Fatalf("%s got '%s' expected '%s'", tt.Method, data, tt.Expected)
		}
		c.Close()
	}
}

var statsdPacketTests2 = []struct {
	Prefix   string
	Method   string
	Stat     string
	Value    int64
	Rate     float32
	Expected string
}{
	{"test", "Gauge", "gauge", 1, 1.0, "test.gauge:1|g"},
	{"test", "Inc", "count", 1, 0.999999, "test.count:1|c|@0.999999"},
	{"test", "Inc", "count", 1, 1.0, "test.count:1|c"},
	{"test", "Dec", "count", 1, 1.0, "test.count:-1|c"},
	{"test", "Timing", "timing", 1, 1.0, "test.timing:1|ms"},
	{"test", "Inc", "count", 1, 1.0, "test.count:1|c"},
	{"test", "GaugeDelta", "gauge", 1, 1.0, "test.gauge:+1|g"},
	{"test", "GaugeDelta", "gauge", -1, 1.0, "test.gauge:-1|g"},
}

func TestBufferedClient2(t *testing.T) {
	l, err := newUDPListener("127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()

	c, err := NewBufferedClient(l.LocalAddr().String(), "test", 10*time.Millisecond, 1024)
	if err != nil {
		t.Fatal(err)
	}
	defer c.Close()

	for _, tt := range statsdPacketTests2 {
		method := reflect.ValueOf(c).MethodByName(tt.Method)
		e := method.Call([]reflect.Value{
			reflect.ValueOf(tt.Stat),
			reflect.ValueOf(tt.Value),
			reflect.ValueOf(tt.Rate)})[0]
		errInter := e.Interface()
		if errInter != nil {
			t.Fatal(errInter.(error))
		}
	}

	var expected string
	for _, tt := range statsdPacketTests2 {
		expected = expected + tt.Expected + "\n"
	}

	data := make([]byte, 1024)
	_, _, err = l.ReadFrom(data)
	if err != nil {
		t.Fatal(err)
	}

	data = bytes.TrimRight(data, "\x00")
	if bytes.Equal(data, []byte(expected)) != true {
		t.Fatalf("got '%s' expected '%s'", data, expected)
	}
}

func ExampleClient_buffered() {
	// first create a client
	client, err := NewBufferedClient("127.0.0.1:8125", "test-client", 10*time.Millisecond, 0)
	// handle any errors
	if err != nil {
		log.Fatal(err)
	}
	// make sure to clean up
	defer client.Close()

	// Send a stat
	err = client.Inc("stat1", 42, 1.0)
	// handle any errors
	if err != nil {
		log.Printf("Error sending metric: %+v", err)
	}
}
