package bitset

import (
	"testing"
)

func TestBitSet_Set(t *testing.T) {
	b := NewBitSet(map[int64]uint64{})
	cases := map[int64]bool{
		-64: true,
		-63: true,
		-2:  true,
		-1:  true,
		0:   true,
		1:   true,
		2:   true,
		3:   false,
		61:  false,
		62:  true,
		63:  true,
		64:  true,
		65:  true,
		66:  false,
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

func TestBitSet_Clear(t *testing.T) {
	b := NewBitSet(map[int64]uint64{})
	cases := map[int64]bool{
		-64: true,
		-63: true,
		-2:  true,
		-1:  true,
		0:   true,
		1:   true,
		2:   true,
		3:   false,
		61:  false,
		62:  true,
		63:  true,
		64:  true,
		65:  true,
		66:  false,
	}
	for key, expect := range cases {
		if expect {
			b.Set(key)
		}
	}
	t.Log(b)

	for key, expect := range cases {
		if expect {
			b.Clear(key)
		}
	}

	t.Log(b)

	if len(b.Storage()) > 0 {
		t.Error("expect empty map")
	}
}

func TestBitSet_Flip(t *testing.T) {
	b := NewBitSet(map[int64]uint64{})
	cases := map[int64]bool{
		-65: false,
		-64: true,
		-63: true,
		-2:  true,
		-1:  true,
		0:   true,
		1:   true,
		2:   true,
		3:   false,
		61:  false,
		62:  true,
		63:  true,
		64:  true,
		65:  true,
		66:  false,
	}
	for key, expect := range cases {
		if expect {
			b.Set(key)
		}
	}
	t.Log(b)
	for key, _ := range cases {
		b.Flip(key)
	}
	t.Log(b)

	for key, expect := range cases {
		if b.Test(key) == expect {
			t.Errorf("%d, expect %v, got %v", key, !expect, b.Test(key))
		}
	}

	if b.Count() != 4 {
		t.Errorf("expect count 3 vs %d", b.Count())
	}
}

func TestBitSet_Intersection(t *testing.T) {
	b1 := NewBitSet(nil)
	b1Val := []int64{-127, -63, -64, -1, 0, 1, 2, 63, 64}
	for _, v := range b1Val {
		b1.Set(v)
	}
	t.Log(b1)

	b2 := NewBitSet(nil)
	b2Val := []int64{-63, -64, 0, 62, 63, 64, 65}
	for _, v := range b2Val {
		b2.Set(v)
	}
	t.Log(b2)

	rlt := b1.Intersection(b2)
	t.Log(rlt)
	rlt.Iterate(func(v int64) bool {
		if !b2.Test(v) || !b1.Test(v) {
			t.Errorf("missing %d", v)
		}
		return true
	})

	b1.Intersection(b2, true)
	t.Log(b1)
}

func TestBitSet_RemoveIntersection(t *testing.T) {
	b1 := NewBitSet(nil)
	b1Val := []int64{-127, -63, -64, -1, 0, 1, 2, 63, 64}
	for _, v := range b1Val {
		b1.Set(v)
	}
	t.Log(b1)

	b2 := NewBitSet(nil)
	b2Val := []int64{-63, -64, 0, 62, 63, 64, 65}
	for _, v := range b2Val {
		b2.Set(v)
	}
	t.Log(b2)

	rlt := b1.RemoveIntersection(b2)
	t.Log(rlt)
	rlt.Iterate(func(v int64) bool {
		if b2.Test(v) {
			t.Errorf("expect remove, got %d", v)
		}
		return true
	})

	b1.RemoveIntersection(b2, true)
	t.Log(b1)
}
