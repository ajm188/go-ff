package jsonutil

import (
	"bytes"
	"io"
	"testing"
)

func TestHTMLUnescaper(t *testing.T) {
	tests := []struct {
		in  string
		out string
	}{
		{
			in:  "hello",
			out: "hello",
		},
		{
			in:  "trailing newline\n",
			out: "trailing newline\n",
		},
		{
			in:  "escaped \u003c",
			out: "escaped <",
		},
		{
			in:  "\u003c\u003e",
			out: "<>",
		},
		{
			in:  "\u0010",
			out: "\u0010",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.out, func(t *testing.T) {
			buf := bytes.NewBuffer(nil)
			buf.WriteString(tt.in)

			e := NewHTMLUnescaper(buf.Bytes())
			buf2 := bytes.NewBuffer(nil)
			io.Copy(buf2, e)

			out := buf2.String()
			if out != tt.out {
				t.Errorf("NewHTMLUnescaper(%s) got = %s want = %s", tt.in, out, tt.out)
			}
		})
	}
}
