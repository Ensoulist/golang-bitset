package bitset

import "math/bits"

type ISyncSource interface {
	Get(int64) (int64, bool)
	Set(int64, int64)
	Delete(int64)
}

type ISyncIteratable interface {
	ISyncSource
	Iterate(func(int64, int64) bool)
}

type BitSetSync struct {
	set ISyncSource
}

func NewBitSetSync(content ISyncSource) *BitSetSync {
	if content == nil {
		panic("NewBitSetSync content is nil")
	}
	return &BitSetSync{set: content}
}

func (b *BitSetSync) From(context ISyncSource) {
	b.set = context
}

func (b *BitSetSync) Storage() ISyncSource {
	return b.set
}

func (b *BitSetSync) Test(key int64) bool {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	if val == 0 {
		return false
	}
	return (val & (1 << innerIdx)) != 0
}

func (b *BitSetSync) Set(key int64) (int64, int64) {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	val |= (1 << innerIdx)
	b.set.Set(outIdx, val)
	return outIdx, val
}

func (b *BitSetSync) Clear(key int64) (int64, int64) {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	if val == 0 {
		return outIdx, val
	}

	val &= ^(1 << innerIdx)
	if val == 0 {
		b.set.Delete(outIdx)
	} else {
		b.set.Set(outIdx, val)
	}
	return outIdx, val
}

func (b *BitSetSync) Flip(key int64) (int64, int64) {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	val |= (1 << innerIdx)
	if val == 0 {
		b.set.Delete(outIdx)
		return outIdx, 0
	} else {
		b.set.Set(outIdx, val)
		return outIdx, val
	}
}

func (b *BitSetSync) Count() int {
	iteratable, ok := b.set.(ISyncIteratable)
	if !ok {
		panic("BitSetSync set is not ISyncIteratable")
	}
	count := 0
	iteratable.Iterate(func(_ int64, v int64) bool {
		count += bits.OnesCount64(uint64(v))
		return true
	})
	return count
}

func (b *BitSetSync) Iterate(do func(int64) bool) {
	iteratable, ok := b.set.(ISyncIteratable)
	if !ok {
		panic("BitSetSync set is not ISyncIteratable")
	}
	iteratable.Iterate(func(outIdx int64, v int64) bool {
		for {
			innerIdx := bits.TrailingZeros64(uint64(v))
			if innerIdx == 64 {
				break
			}
			if !do(int64(outIdx)*unit_bit_len + int64(innerIdx)) {
				break
			}
			v = v & ^(1 << innerIdx)
		}
		return true
	})
}
