package bitset

import (
	"math/bits"
	"testing"
)

func TestBitSet(t *testing.T) {
	innerIdx := bits.TrailingZeros64(uint64(1))
	print(innerIdx)
	b := NewBitSet(map[int64]int64{})
	cases := map[int64]bool{
		// -64: true,
		// -63: true,
		// -2:  true,
		// -1:  true,
		0:  true,
		1:  true,
		2:  true,
		3:  false,
		61: false,
		62: true,
		63: true,
		64: true,
		65: true,
		66: false,
	}
	for key, expect := range cases {
		if expect {
			b.Set(key)
		}
	}
	t.Log(b)

	for key, expect := range cases {
		if b.Test(key) != expect {
			t.Errorf("expect %v, got %v", expect, b.Test(key))
		}
	}
}
