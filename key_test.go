package critbit

import (
	"testing"
)

func TestKey_Critbit(t *testing.T) {
	tests := []struct {
		name string
		k    Key
		b    Key
		bit  int
	}{
		{
			name: "uint64 same",
			k:    Uint64Key(1),
			b:    Uint64Key(1),
			bit:  -1,
		},
		{
			name: "uint32 same",
			k:    Uint32Key(1),
			b:    Uint32Key(1),
			bit:  -1,
		},
		{
			name: "uint16 same",
			k:    Uint16Key(1),
			b:    Uint16Key(1),
			bit:  -1,
		},
		{
			name: "uint8 same",
			k:    Uint8Key(1),
			b:    Uint8Key(1),
			bit:  -1,
		},
		{
			name: "string same",
			k:    StringKey("same string"),
			b:    StringKey("same string"),
			bit:  -1,
		},
		{
			name: "uint32 diff data 0bit",
			k:    Uint32Key(1 << 31),
			b:    Uint32Key(0),
			bit:  1,
		},
		{
			name: "uint32 diff data 1bit",
			k:    Uint32Key(1 << 30),
			b:    Uint32Key(0),
			bit:  3,
		},
		{
			name: "uint32 diff data 2bit",
			k:    Uint32Key(1 << 29),
			b:    Uint32Key(0),
			bit:  5,
		},
		{
			name: "uint32 diff data 3bit",
			k:    Uint32Key(1 << 28),
			b:    Uint32Key(0),
			bit:  7,
		},
		{
			name: "uint32 diff data 4bit",
			k:    Uint32Key(1 << 27),
			b:    Uint32Key(0),
			bit:  9,
		},
		{
			name: "uint32 diff data 30bit",
			k:    Uint32Key(1 << 1),
			b:    Uint32Key(0),
			bit:  61,
		},
		{
			name: "uint32 diff data 31bit",
			k:    Uint32Key(1),
			b:    Uint32Key(0),
			bit:  63,
		},
		{
			name: "1bit same",
			k:    BitsKey([]byte{0b1000_0000}, 1),
			b:    BitsKey([]byte{0b1000_0000}, 1),
			bit:  -1,
		},
		{
			name: "2bit same",
			k:    BitsKey([]byte{0b1100_0000}, 2),
			b:    BitsKey([]byte{0b1100_0000}, 2),
			bit:  -1,
		},
		{
			name: "1bit diff data 0bit",
			k:    BitsKey([]byte{0b1000_0000}, 1),
			b:    BitsKey([]byte{0b0000_0000}, 1),
			bit:  1,
		},
		{
			name: "2bit diff nbit 1bit",
			k:    BitsKey([]byte{0b1000_0000}, 1),
			b:    BitsKey([]byte{0b1000_0000}, 2),
			bit:  2,
		},
		{
			name: "2bit diff nbit 1bit",
			k:    BitsKey([]byte{0b1000_0000}, 2),
			b:    BitsKey([]byte{0b1000_0000}, 1),
			bit:  2,
		},
		{
			name: "10bit diff nbit 9bit",
			k:    BitsKey([]byte{0, 0b1000_0000}, 9),
			b:    BitsKey([]byte{0, 0b1000_0000}, 10),
			bit:  18,
		},
		{
			name: "10bit diff nbit 9bit",
			k:    BitsKey([]byte{0, 0b1000_0000}, 10),
			b:    BitsKey([]byte{0, 0b1000_0000}, 9),
			bit:  18,
		},
		{
			name: "10bit diff nbit 2bit",
			k:    BitsKey([]byte{0b1000_0000}, 2),
			b:    BitsKey([]byte{0b1000_0000, 0}, 10),
			bit:  4,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.k.Critbit(tc.b)
			if got != tc.bit {
				t.Errorf("want %v; but got %v", tc.bit, got)
			}
		})
	}
}
