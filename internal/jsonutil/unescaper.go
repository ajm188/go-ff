package jsonutil

import (
	"bytes"
	"io"
)

var escapeChars = map[byte]byte{
	0x003c: '<',
	0x003e: '>',
	0x0026: '&',
}

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
	i := 0

	for i < len(b) {
		if u.offset >= len(u.buf) {
			return n, io.EOF
		}

		if len(u.buf[u.offset:]) < 6 {
			x := copy(b[i:], u.buf[u.offset:])
			n += x
			u.offset += x
			break
		}

		j := bytes.Index(u.buf[u.offset:], []byte(`\u`))
		if j == -1 {
			// no escaped characters remaining in u.buf, just copy everything
			x := copy(b[i:], u.buf[u.offset:])
			n += x
			u.offset += x
			break
		}

		if j > 0 {
			// escaped character is not the front of u.buf[u.offset:], copy
			// everything up to it first
			x := copy(b[i:], u.buf[u.offset:u.offset+j])
			i += x
			n += x
			u.offset += x
		}

		if i >= len(b) {
			// We've filled the destination buffer
			break
		}

		if len(u.buf[u.offset:]) < 6 {
			// even though we saw the beginning of an escape sequence, there's
			// not even remaining characters to have a real escape sequence.
			// this is a duplication of the block at the beginning of the copy
			// loop
			x := copy(b[i:], u.buf[u.offset:])
			n += x
			u.offset += x
			break
		}

		sextet := u.buf[u.offset : u.offset+6]
		hexcode := sextet[2]<<3 | sextet[3]<<2 | sextet[4]<<1 | sextet[5]

		if char, ok := escapeChars[hexcode]; ok {
			b[i] = char
			n++
			u.offset += 6
		} else {
			// Best we can do is take a single byte and go to the beginning of
			// the loop to rescan for the next potential escape sequence
			b[i] = u.buf[u.offset]
			n++
			u.offset++
		}

		i++
	}

	if u.offset >= len(u.buf) {
		return n, io.EOF
	}

	return n, nil
}
