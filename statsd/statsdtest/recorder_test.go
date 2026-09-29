package statsdtest

import (
	"fmt"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/cactus/go-statsd-client/v6/statsd"
)

func TestRecordingSenderIsSender(t *testing.T) {
	// This ensures that if the Sender interface changes in the future we'll get
	// compile time failures should the RecordingSender not be updated to meet
	// the new definition. This keeps changes from inadvertently breaking tests
	// of folks that use go-statsd-client.
	var _ statsd.Sender = NewRecordingSender()
}

func TestRecordingSender(t *testing.T) {
	start := time.Now()
	rs := new(RecordingSender)
	statter, err := statsd.NewClientWithSender(rs, "test", 0)
	if err != nil {
		t.Errorf("failed to construct client")
		return
	}

	statter.Inc("stat", 4444, 1.0)
	statter.Dec("stat", 5555, 1.0)
	statter.Set("set-stat", "some string", 1.0)

	d := time.Since(start)
	statter.TimingDuration("timing", d, 1.0)

	sent := rs.GetSent()
	if len(sent) != 4 {
		// just dive out because everything else relies on ordering
		t.Fatalf("Did not capture all stats sent; got: %s", sent)
	}

	ms := float64(d) / float64(time.Millisecond)
	// somewhat fragile in that it assumes float rendering within client *shrug*
	msStr := string(strconv.AppendFloat([]byte(""), ms, 'f', -1, 64))

	expected := Stats{
		{"test.stat", "4444", "c", "", []byte("test.stat:4444|c"), true},
		{"test.stat", "-5555", "c", "", []byte("test.stat:-5555|c"), true},
		{"test.set-stat", "some string", "s", "", []byte("test.set-stat:some string|s"), true},
		{"test.timing", msStr, "ms", "", []byte(fmt.Sprintf("test.timing:%s|ms", msStr)), true},
	}

	if !reflect.DeepEqual(sent, expected) {
		t.Errorf("got: %s, want: %s", sent, expected)
	}
}
