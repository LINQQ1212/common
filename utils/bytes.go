package utils

import "encoding/binary"

func Itob(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func Btoi(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

/*

func Itob(val uint64) []byte {
	r := make([]byte, 8)
	for i := uint64(0); i < 8; i++ {
		r[i] = byte((val >> (i * 8)) & 0xff)
	}
	return r
}
func Btoi(val []byte) uint64 {
	r := uint64(0)
	for i := uint64(0); i < 8; i++ {
		r |= uint64(val[i]) << (8 * i)
	}
	return r
}


	func Itob(v uint64) []byte {
		buf := make([]byte, binary.MaxVarintLen64)
		n := binary.PutUvarint(buf, v)
		b := buf[:n]
		return b
	}

func Btoi(raw []byte) uint64 {
	x, _ := binary.Uvarint(raw)
	return x
}*/
