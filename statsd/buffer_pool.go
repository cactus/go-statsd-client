package statsd

import (
	"bytes"
	"sync"
)

type bufferPool struct {
	sync.Pool
}

func newBufferPool() *bufferPool {
	bp := &bufferPool{}

	bp.New = func() interface{} {
		return &bytes.Buffer{}
	}

	return bp
}

func (bp *bufferPool) Get() *bytes.Buffer {
	b := (bp.Pool.Get()).(*bytes.Buffer)
	b.Truncate(0)
	return b
}

func (bp *bufferPool) Put(b *bytes.Buffer) {
	bp.Pool.Put(b)
}
