package critbit

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"unsafe"
)

type Key struct {
	Data  []byte
	Nbits int
}

func BitsKey(b []byte, nbits int) Key {
	return Key{Data: b, Nbits: nbits}
}

func BytesKey(b []byte) Key {
	return Key{Data: b, Nbits: len(b) * 8}
}

func Uint64Key(n uint64) Key {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return BytesKey(b)
}

func Uint32Key(n uint32) Key {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return BytesKey(b)
}

func Uint16Key(n uint16) Key {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return BytesKey(b)
}

func Uint8Key(n uint8) Key {
	return BytesKey([]byte{n})
}

func StringKey(s string) Key {
	return BytesKey(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (k Key) Equal(b Key) bool {
	if k.Nbits != b.Nbits {
		return false
	}
	return bytes.Equal(k.Data, b.Data)
}

func (k Key) Critbit(b Key) int {
	koff := k.Nbits >> 3
	boff := b.Nbits >> 3

	var moff, mod int
	if koff < boff {
		moff = koff
		mod = k.Nbits & 7
	} else {
		moff = boff
		mod = b.Nbits & 7
	}

	off := 0
	for ; off < moff; off++ {
		d := k.Data[off] ^ b.Data[off]
		if d != 0 {
			msbbit := bits.LeadingZeros8(d)
			return off<<4 | msbbit<<1 | 1
		}
	}
	if mod > 0 {
		d := k.Data[off] ^ b.Data[off]
		d &= ^byte(0) << (8 - mod)
		if d != 0 {
			msbbit := bits.LeadingZeros8(d)
			return off<<4 | msbbit<<1 | 1
		}
	}
	if k.Nbits == b.Nbits {
		return -1
	} else if k.Nbits < b.Nbits {
		return k.Nbits << 1
	} else {
		return b.Nbits << 1
	}
}

func (k Key) Direction(bit int) int {
	koff := (k.Nbits + 7) >> 3
	cbit := bit >> 1
	coff := cbit >> 3
	if coff < koff {
		if k.Data[coff]&(0x80>>(cbit&0x7)) != 0 {
			return 1
		}
	}
	if bit&1 != 0 {
		return 0
	}
	if cbit < k.Nbits {
		return 1
	}
	return 0
}

func (k Key) HasPrefix(p Key) bool {
	if k.Nbits < p.Nbits {
		return false
	}
	n := p.Nbits >> 3
	off := 0
	for ; off < n; off++ {
		if k.Data[off] != p.Data[off] {
			return false
		}
	}
	mod := p.Nbits & 7
	if mod > 0 {
		d := k.Data[off] ^ p.Data[off]
		d &= ^byte(0) << (8 - mod)
		if d != 0 {
			return false
		}
	}
	return true
}
