package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/samber/lo"
	"math/rand"
	"strings"
)

// StrMd5 str to md5
func StrMd5(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	rs := hex.EncodeToString(h.Sum(nil))
	return rs
}

func RandInt(min, max int) int {
	if min >= max || min == 0 && max == 0 {
		return max
	}
	return rand.Intn(max-min) + min
}

const Char = "abcdefghijklmnopqrstuvwxyz0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandChars() string {
	return strings.Join(lo.Shuffle[string](strings.Split(Char, "")), "")
}
func RandCharsLen(i int) string {
	if i > 36 || i <= 0 {
		return RandChars()
	}
	s := RandChars()
	return s[0:i]
}

func RandChar(size int) string {
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(Char[rand.Int63()%62])
	}
	return s.String()
}

func MicsSlice(origin []int, count int) []int {
	tmpOrigin := lo.Shuffle[int](origin)
	return tmpOrigin[:count]
}

func RandomStr(str string) string {
	bs := []byte(str)
	for i := len(bs) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		bs[i], bs[num] = bs[num], bs[i]
	}
	return string(bs)
}
