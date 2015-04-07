// Package cyrutf determine cyrrilic charset by statistics of byte-pairs
package cyrutf

import (
	"errors"
	"math"

	. "github.com/srf/cyrutf/pairs"
)

// Const to string translation
var ux = map[byte]string{
	CP866:  "cp866",
	CP1251: "windows-1251",
	ISO:    "iso-8859-5", // iso-8859-5
	KOI8:   "koi8",
}

// DetermineEncoding determines cyrillic charset and return string-name, charset-stat and error.
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
		return "", nil, errors.New("Charset not found")
	}
	return charset, stat, nil
}

// Calc count pairs of bytes for cyrillic charsets and return statistics as map of
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
