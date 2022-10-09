package statsd

import "time"

// NullSender  is disabled sender (if stat need to be disabled)
type NullSender struct{}

func (NullSender) Inc(string, int64, float32, ...Tag) error                    { return nil }
func (NullSender) Dec(string, int64, float32, ...Tag) error                    { return nil }
func (NullSender) Gauge(string, int64, float32, ...Tag) error                  { return nil }
func (NullSender) GaugeDelta(string, int64, float32, ...Tag) error             { return nil }
func (NullSender) Timing(string, int64, float32, ...Tag) error                 { return nil }
func (NullSender) TimingDuration(string, time.Duration, float32, ...Tag) error { return nil }
func (NullSender) Set(string, string, float32, ...Tag) error                   { return nil }
func (NullSender) SetInt(string, int64, float32, ...Tag) error                 { return nil }
func (NullSender) Raw(string, string, float32, ...Tag) error                   { return nil }
func (NullSender) NewSubStatter(string) SubStatter                             { return NullSender{} }
func (NullSender) SetPrefix(string)                                            {}
func (NullSender) SetSamplerFunc(SamplerFunc)                                  {}
func (NullSender) Close() error                                                { return nil }
