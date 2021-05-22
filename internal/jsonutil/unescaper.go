package jsonutil

import (
	"bytes"
	"io"
)

var (
	leftAngle  = []byte("003c")
	rightAngle = []byte("003e")
	ampersand  = []byte("0026")
)

type HTMLUnescaper struct {
	buf    []byte
	offset int
}

func NewHTMLUnescaper(buf []byte) *HTMLUnescaper {
	return &HTMLUnescaper{
		buf: buf,
	}
}

func (u *HTMLUnescaper) Read(b []byte) (int, error) {
	n := 0

	for i := 0; i < len(b); i++ {
		if u.offset >= len(u.buf) {
			return n, io.EOF
		}

		if len(u.buf)-u.offset < 6 {
			b[i] = u.buf[u.offset]
			n++
			u.offset++
			continue
		}

		sextet := u.buf[u.offset : u.offset+6]
		if !bytes.HasPrefix(sextet, []byte(`\u`)) {
			b[i] = u.buf[u.offset]
			n++
			u.offset++
			continue
		}

		switch 0 {
		case bytes.Compare(sextet[2:], leftAngle):
			b[i] = '<'
			n++
			u.offset += 6
		case bytes.Compare(sextet[2:], rightAngle):
			b[i] = '>'
			n++
			u.offset += 6
		case bytes.Compare(sextet[2:], ampersand):
			b[i] = '&'
			n++
			u.offset += 6
		default:
			b[i] = u.buf[u.offset]
			n++
			u.offset++
		}
	}

	return n, nil
}
