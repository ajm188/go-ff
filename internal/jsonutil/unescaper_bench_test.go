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
		})
	}
}
