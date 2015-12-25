// Package cyrutf determine cyrrilic charset by statistics of dibytes
package cyrutf

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"strings"

	. "github.com/c4e8ece0/cyrutf/pairs"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/unicode/norm"
)

var ErrCharsetNotFound = errors.New("Charset not found")

// Const to string translation
var ux = map[byte]string{
	CP866:   "cp866",
	CP1251:  "windows-1251",
	ISO:     "iso-8859-5",
	KOI8:    "koi8",
	UTF16BE: "utf-16 be",
	UTF16LE: "utf-16 le",
	UTF8:    "utf-8",
}

// NewReader return io.Reader with utf-8 encoded data
func NewReader(r io.Reader) (io.Reader, error) {
	str, _ := ioutil.ReadAll(r)
	c, _, err := DetermineEncoding(str)
	enc := string(c)
	if err != nil {
		_, p, _ := charset.DetermineEncoding(str, "text/html") // works only on utf-8
		c = p
	}
	if enc == "windows-1252" {
		enc = "windows-1251"
	}
	if enc == "" {
		enc = "utf-8" // in the name of universe
	}
	return norm.NFC.Reader(charset.NewReaderLabel(enc, strings.NewReader(string(str))))

}

// Determine cyrillic charset and return string-name, charset-stat and error.
func DetermineEncoding(a []byte) (string, map[byte]float32, error) {
	stat := Calc(a)
	var max float32 = -1
	var charset string
	var found bool = false
	for enc, w := range stat {
		if w > max {
			found = true
			charset = ux[enc]
			max = w
		}
	}
	if !found {
		return "", nil, ErrCharsetNotFound
	}
	return charset, stat, nil
}

// Count pairs of bytes for cyrillic charsets and return statistics as map of
// internal_id => value
func Calc(a []byte) map[byte]float32 {
	l := len(a) - 1
	stat := make(map[byte]float32)
	cnt := make(map[byte]int)

	for i := 0; i < l; i++ {
		key := uint16(a[i])<<8 | uint16(a[i+1])
		if _, has := Pairs[key]; has {
			for enc, weight := range Pairs[key] {
				stat[enc] += float32(math.Log10(float64(weight)))
				cnt[enc]++
			}
		}
	}
	for k, v := range stat {
		stat[k] = v * float32(cnt[k])
	}
	return stat
}
