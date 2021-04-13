package packet

import (
	"bufio"
	"io"
)

// need to make this guy implement io.Reader. Just like bufio.Reader does. That way, we can treat them the same!!
type Reader struct {
	rdr       *bufio.Reader
	bytesRead int
}

// NewReader returns a new Reader whose buffer has the default size.
func NewReader(rdr io.Reader) *Reader {
	return &Reader{bufio.NewReader(rdr), 0}
}

// ReadByte reads and returns a single byte.
// If no byte is available, returns an error.
func (rdr *Reader) ReadByte() (byte, error) {
	b, err := rdr.rdr.ReadByte()
	if err != nil {
		return 0, err
	}
	rdr.bytesRead++
	return b, nil
}

// ResetBytesRead sets bytesRead to zero. Does NOT unread any bytes.
func (rdr *Reader) ResetBytesRead() {
	rdr.bytesRead = 0
}
