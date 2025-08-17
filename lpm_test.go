package critbit

import (
	"net/netip"
	"testing"
)

func TestLongest(t *testing.T) {
	var m Tree[int]
	addrs := []string{
		"10.1.2.1/32",
		"10.1.2.0/24",
		"10.1.0.0/16",
		"10.0.0.0/8",
		"0.0.0.0/4",
		"0.0.0.0/8",
		"1.0.0.0/8",
		"0.0.0.0/7",
	}
	for i, addr := range addrs {
		p := netip.MustParsePrefix(addr)
		key := Key{
			Data:  p.Addr().AsSlice(),
			Nbits: p.Bits(),
		}
		m.Set(key, i)
	}

	tests := []struct {
		name   string
		prefix netip.Prefix
		v      int
		found  bool
	}{
		{
			name:   "10.1.0.0/16",
			prefix: netip.MustParsePrefix("10.1.0.0/16"),
			v:      2,
			found:  true,
		},
		{
			name:   "10.1.1.0/24",
			prefix: netip.MustParsePrefix("10.1.1.0/24"),
			v:      2,
			found:  true,
		},
		{
			name:   "10.1.2.8/30",
			prefix: netip.MustParsePrefix("10.1.2.8/30"),
			v:      1,
			found:  true,
		},
		{
			name:   "0.0.0.0/7",
			prefix: netip.MustParsePrefix("0.0.0.0/7"),
			v:      7,
			found:  true,
		},
		{
			name:   "8.0.0.0/5",
			prefix: netip.MustParsePrefix("8.0.0.0/5"),
			v:      4,
			found:  true,
		},
		{
			name:   "0.0.0.0/9",
			prefix: netip.MustParsePrefix("0.0.0.0/9"),
			v:      5,
			found:  true,
		},
		{
			name:   "0.0.0.0/3",
			prefix: netip.MustParsePrefix("0.0.0.0/3"),
			v:      0,
			found:  false,
		},
		{
			name:   "16.0.0.0/4",
			prefix: netip.MustParsePrefix("16.0.0.0/4"),
			v:      0,
			found:  false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			key := Key{
				Data:  tc.prefix.Addr().AsSlice(),
				Nbits: tc.prefix.Bits(),
			}
			got, found := m.Longest(key)
			if found != tc.found {
				t.Errorf("want %v; but got %v", tc.found, found)
			}
			if got != tc.v {
				t.Errorf("want %v; but got %v", tc.v, got)
			}
		})
	}
}
