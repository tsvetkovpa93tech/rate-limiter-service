package pool

import (
	"bytes"
	"sync"
)

var (
	// BufferPool is a pool for bytes.Buffer objects
	BufferPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}

	// JSONEncoderPool is a pool for JSON encoders
	JSONEncoderPool = sync.Pool{
		New: func() interface{} {
			return &bytes.Buffer{}
		},
	}
)

// GetBuffer retrieves a buffer from the pool
func GetBuffer() *bytes.Buffer {
	return BufferPool.Get().(*bytes.Buffer)
}

// PutBuffer returns a buffer to the pool
func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	BufferPool.Put(buf)
}

// GetJSONBuffer retrieves a JSON buffer from the pool
func GetJSONBuffer() *bytes.Buffer {
	return JSONEncoderPool.Get().(*bytes.Buffer)
}

// PutJSONBuffer returns a JSON buffer to the pool
func PutJSONBuffer(buf *bytes.Buffer) {
	buf.Reset()
	JSONEncoderPool.Put(buf)
}

