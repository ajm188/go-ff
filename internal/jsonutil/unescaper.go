package jsonutil

import (
	"bytes"
	"io"
)

var escapeChars = map[string]byte{
	"003c": '<',
	"003e": '>',
	"0026": '&',
}

// HTMLUnescaper provides an io.Reader implementation that does the inverse of
// json.HTMLEscape as bytes are read off of it.
//
// N.B.: The current implementation only unescapes the <, >, and & characters.
// The U+2028 and U+2029 characters remain escaped.
type HTMLUnescaper struct {
	buf    []byte
	offset int
}

// NewHTMLUnescaper returns an HTMLUnescaper that unescapes bytes in buf when
// read from. The HTMLUnescaper takes ownership of, but does not mutate, the
// passed-in buffer; therefore callers may continue to use the original buffer
// when they are finished using the unescaper.
func NewHTMLUnescaper(buf []byte) *HTMLUnescaper {
	return &HTMLUnescaper{
		buf: buf,
	}
}

// Read implements io.Reader, unescaping HTML-escaped characters as they are
// read from the underlying buffer.
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
		hexcode := string(sextet[2:])

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
