package critbit

import (
	"bytes"
	"encoding/binary"
	"math/bits"
	"unsafe"
)

// Key represents a binary key with arbitrary bit length.
// It consists of a byte slice containing the key data and
// the number of significant bits in the key.
type Key struct {
	// Data contains the binary representation of the key
	Data []byte
	// Nbits specifies the number of significant bits in the key
	Nbits int
}

// BitsKey creates a Key from a byte slice with a specific number of bits.
// This allows for keys that are not byte-aligned.
//
// Example:
//
//	key := BitsKey([]byte{0b10110000}, 5) // Only first 5 bits are significant
func BitsKey(b []byte, nbits int) Key {
	return Key{Data: b, Nbits: nbits}
}

// BytesKey creates a Key from a byte slice where all bits are significant.
// This is equivalent to BitsKey(b, len(b)*8).
func BytesKey(b []byte) Key {
	return Key{Data: b, Nbits: len(b) * 8}
}

// Uint64Key creates a Key from a uint64 value stored in big-endian format.
// The resulting key will have 64 significant bits.
func Uint64Key(n uint64) Key {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, n)
	return BytesKey(b)
}

// Uint32Key creates a Key from a uint32 value stored in big-endian format.
// The resulting key will have 32 significant bits.
func Uint32Key(n uint32) Key {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	return BytesKey(b)
}

// Uint16Key creates a Key from a uint16 value stored in big-endian format.
// The resulting key will have 16 significant bits.
func Uint16Key(n uint16) Key {
	b := make([]byte, 2)
	binary.BigEndian.PutUint16(b, n)
	return BytesKey(b)
}

// Uint8Key creates a Key from a uint8 value.
// The resulting key will have 8 significant bits.
func Uint8Key(n uint8) Key {
	return BytesKey([]byte{n})
}

// StringKey creates a Key from a string.
// All bytes of the string are considered significant bits.
func StringKey(s string) Key {
	return BytesKey(unsafe.Slice(unsafe.StringData(s), len(s)))
}

// Equal reports whether k and b represent the same key.
// Two keys are equal if they have the same number of significant bits
// and identical data content.
func (k Key) Equal(b Key) bool {
	if k.Nbits != b.Nbits {
		return false
	}
	return bytes.Equal(k.Data, b.Data)
}

// Critbit returns the position of the first critical bit
// where keys k and b differ.
// This is the core algorithm used by crit-bit trees to determine
// branching points.
//
// Returns:
//
//   - -1 if the keys are identical
//   - A bit position (encoded) where the keys first differ
//
// The returned value encodes both the byte offset and bit position:
//
//   - For data differences: (byte_offset << 4) | (bit_offset << 1) | 1
//   - For length differences: shorter_length << 1
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

	// Compare full bytes
	off := 0
	for ; off < moff; off++ {
		d := k.Data[off] ^ b.Data[off]
		if d != 0 {
			msbbit := bits.LeadingZeros8(d)
			return off<<4 | msbbit<<1 | 1
		}
	}

	// Compare partial byte if needed
	if mod > 0 {
		d := k.Data[off] ^ b.Data[off]
		d &= ^byte(0) << (8 - mod)
		if d != 0 {
			msbbit := bits.LeadingZeros8(d)
			return off<<4 | msbbit<<1 | 1
		}
	}

	// Keys are identical in data, check lengths
	if k.Nbits == b.Nbits {
		return -1
	} else if k.Nbits < b.Nbits {
		return k.Nbits << 1
	} else {
		return b.Nbits << 1
	}
}

// Direction determines which branch to take at a given bit position
// during tree traversal. This is used by the crit-bit tree to decide
// whether to go left (0) or right (1) at an internal node.
//
// The bit parameter encodes the critical bit position as returned
// by Critbit.
// Returns 0 for left branch, 1 for right branch.
func (k Key) Direction(bit int) int {
	koff := (k.Nbits + 7) >> 3
	cbit := bit >> 1
	coff := cbit >> 3

	// Check if the bit position is within the key data
	if coff < koff {
		if k.Data[coff]&(0x80>>(cbit&0x7)) != 0 {
			return 1
		}
	}

	// Handle length-based critical bits
	if bit&1 != 0 {
		return 0
	}
	if cbit < k.Nbits {
		return 1
	}
	return 0
}

// HasPrefix reports whether key k has p as a prefix.
// This is useful for longest prefix matching operations.
//
// A key has a prefix if:
//  1. The prefix is not longer than the key
//  2. All bits of the prefix match the corresponding bits in the key
func (k Key) HasPrefix(p Key) bool {
	if k.Nbits < p.Nbits {
		return false
	}

	// Compare full bytes
	n := p.Nbits >> 3
	off := 0
	for ; off < n; off++ {
		if k.Data[off] != p.Data[off] {
			return false
		}
	}

	// Compare partial byte if needed
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
