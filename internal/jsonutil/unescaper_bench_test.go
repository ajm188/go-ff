package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func listFiles(b *testing.B, dir string) []string {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		b.Fatal(err)
	}

	filenames := make([]string, 0, len(entries))
	for _, f := range entries {
		if f.IsDir() {
			continue
		}

		filenames = append(filenames, f.Name())
	}

	return filenames
}

type buffer struct {
	buf []byte
}

func (b *buffer) Read(d []byte) (n int, err error) {
	n = copy(d, b.buf)

	if n >= len(b.buf) {
		err = io.EOF
	}

	return n, err
}

func BenchmarkHTMLUnescaper(b *testing.B) {
	filenames := listFiles(b, "testdata")

	for _, filename := range filenames {
		b.Run(filename, func(b *testing.B) {
			data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
			if err != nil {
				b.Fatal(err)
			}

			buf := bytes.NewBuffer(nil)
			json.HTMLEscape(buf, data)
			data = buf.Bytes()

			b.ResetTimer()

			b.Run("HTMLUnescaper", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					buf := bytes.NewBuffer(nil)
					unescaper := NewHTMLUnescaper(data)
					_, err := io.Copy(buf, unescaper)
					if err != nil {
						b.Error(err)
					}
				}
			})

			data2 := make([]byte, len(data))
			copy(data2, data)

			b.Run("io.Copy", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					buf := bytes.NewBuffer(nil)
					_, err := io.Copy(buf, bytes.NewBuffer(data2))
					if err != nil {
						b.Error(err)
					}
				}
			})

			b.Run("copy passthrough", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					buf := bytes.NewBuffer(nil)
					_, err := io.Copy(buf, &buffer{data})
					if err != nil {
						b.Error(err)
					}
				}
			})
		})
	}
}

func BenchmarkHTMlUnescaperMarshal(b *testing.B) {
	filenames := listFiles(b, "testdata")

	for _, filename := range filenames {
		b.Run(filename, func(b *testing.B) {
			data, err := ioutil.ReadFile(filepath.Join("testdata", filename))
			if err != nil {
				b.Fatal(err)
			}

			var m []map[string]interface{}
			if err := json.Unmarshal(data, &m); err != nil {
				b.Fatal(err)
			}

			b.ResetTimer()

			b.Run("json.Encode", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					buf := bytes.NewBuffer(nil)
					encoder := json.NewEncoder(buf)
					encoder.SetEscapeHTML(false)

					if err := encoder.Encode(m); err != nil {
						b.Fatal(err)
					}
				}
			})

			b.Run("json.Marshal", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					if _, err := json.Marshal(m); err != nil {
						b.Fatal(err)
					}
				}
			})

			b.Run("HTMLUnescaper", func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					data, err := json.Marshal(m)
					if err != nil {
						b.Fatal(err)
					}

					if _, err := ioutil.ReadAll(NewHTMLUnescaper(data)); err != nil {
						b.Fatal(err)
					}
				}
			})
		})
	}
}
