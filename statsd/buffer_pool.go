package statsd

import (
	"bytes"
	"sync"
)

type bufferPool struct {
	*sync.Pool
}

func newBufferPool() *bufferPool {
	return &bufferPool{
		&sync.Pool{New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 0, 1500))
		}},
	}
}

func (bp *bufferPool) Get() *bytes.Buffer {
	return (bp.Pool.Get()).(*bytes.Buffer)
}

func (bp *bufferPool) Put(b *bytes.Buffer) {
	b.Truncate(0)
	bp.Pool.Put(b)
}
