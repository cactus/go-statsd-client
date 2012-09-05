package statsd

import (
	"bufio"
	"bytes"
	"testing"
	//"fmt"
)

func NewTestClient(prefix string) (*Client, *bytes.Buffer) {
	b := &bytes.Buffer{}
	buf := bufio.NewReadWriter(bufio.NewReader(b), bufio.NewWriter(b))
	f := &Client{buf: buf, prefix: prefix}
	return f, b
}

func TestGauge(t *testing.T) {
	f, buf := NewTestClient("test")

	err := f.Gauge("gauge", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "test.gauge:1|g"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestIncRatio(t *testing.T) {
	f, buf := NewTestClient("test")

	err := f.Inc("count", 1, 0.999999)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "test.count:1|c|@0.999999"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestInc(t *testing.T) {
	f, buf := NewTestClient("test")

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "test.count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestDec(t *testing.T) {
	f, buf := NewTestClient("test")

	err := f.Dec("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "test.count:-1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestTiming(t *testing.T) {
	f, buf := NewTestClient("test")

	err := f.Timing("timing", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "test.timing:1|ms"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestEmptyPrefix(t *testing.T) {
	f, buf := NewTestClient("")

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := buf.String()
	buf.Reset()
	expected := "count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func ExampleClient() {
	// first create a client
	client, err := statsd.Dial("127.0.0.1:8125", "test-client")
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
