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
			in:  `escaped \u003c`,
			out: "escaped <",
		},
		{
			in:  `\u003c\u003e`,
			out: "<>",
		},
		{
			in:  `\u0010`,
			out: `\u0010`,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.out, func(t *testing.T) {
			e := NewHTMLUnescaper([]byte(tt.in))
			buf := bytes.NewBuffer(nil)
			io.Copy(buf, e)

			out := buf.Bytes()
			if bytes.Compare(out, []byte(tt.out)) != 0 {
				t.Errorf("NewHTMLUnescaper(%s) got = %s want = %s", tt.in, out, tt.out)
			}
		})
	}
}
