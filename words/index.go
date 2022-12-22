package words

import (
	"golang.org/x/exp/rand"
	"unsafe"
)

func GetWords(n int, index int) []string {
	var s []string
	for i := 0; i < n; i++ {
		num := index << i % 20000
		s = append(s, Word[num])
	}
	return s
}

const (
	letterBytes   = "abcdefghijklmnopqrstuvwxyz"
	letterIdxMask = 1<<6 - 1 // All 1-bits, as many as letterIdxBits
)

func RandStr(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), 10; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), 10
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= 6
		remain--
	}
	return *(*string)(unsafe.Pointer(&b))
}
