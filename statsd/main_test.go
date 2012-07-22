package statsd

import (
	"testing"
	"bytes"
	"net"
	//"fmt"
)

type fakeUDPConn struct {
	*net.UDPConn
	buf bytes.Buffer
}

func (f *fakeUDPConn) GetBuff() string {
	res := f.buf.String()
	f.buf.Reset()
	return res
}

func (f *fakeUDPConn) Write(b []byte) (int, error) {
	i, err := f.buf.Write(b)
	return i, err
}


func TestGuage(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u, prefix: "test"}

	err := f.Guage("guage", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "test.guage:1|g"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestIncRatio(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u, prefix: "test"}

	err := f.Inc("count", 1, 0.999999)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "test.count:1|c|@0.999999"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestInc(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u, prefix: "test"}

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "test.count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestDec(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u, prefix: "test"}

	err := f.Dec("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "test.count:-1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestTiming(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u, prefix: "test"}

	err := f.Timing("timing", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "test.timing:1|ms"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

func TestEmptyPrefix(t *testing.T) {
	u := &fakeUDPConn{}
	f := &Client{conn: u}

	err := f.Inc("count", 1, 1.0)
	if err != nil {
		t.Fatal(err)
	}

	b := u.GetBuff()
	expected := "count:1|c"
	if b != expected {
		t.Fatalf("got '%s' expected '%s'", b, expected)
	}
}

