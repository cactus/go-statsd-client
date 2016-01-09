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
	return (bp.Pool.Get()).(*bytes.Buffer)
}

func (bp *bufferPool) Put(b *bytes.Buffer) {
	b.Truncate(0)
	bp.Pool.Put(b)
}
