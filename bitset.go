package bitset

import (
	"fmt"
	"math/bits"
)

var unit_bit_len = int64(64)

type BitSet struct {
	set map[int64]uint64
}

func NewBitSet(content map[int64]uint64) *BitSet {
	if content == nil {
		content = map[int64]uint64{}
	}
	return &BitSet{set: content}
}

func (b *BitSet) From(context map[int64]uint64) {
	b.set = context
}

func (b *BitSet) Clone() *BitSet {
	newSet := make(map[int64]uint64, len(b.set))
	for k, v := range b.set {
		newSet[k] = v
	}
	return NewBitSet(newSet)
}

func (b *BitSet) Storage() map[int64]uint64 {
	return b.set
}

func (b *BitSet) Test(key int64) bool {
	outIdx, innerIdx := key_2_idx(key)
	val := b.set[outIdx]
	if val == 0 {
		return false
	}
	return (val & (1 << innerIdx)) != 0
}

func (b *BitSet) Set(key int64) (int64, uint64) {
	outIdx, innerIdx := key_2_idx(key)
	val := b.set[outIdx]
	val |= (1 << innerIdx)
	b.set[outIdx] = val
	return outIdx, val
}

func (b *BitSet) Clear(key int64) (int64, uint64) {
	outIdx, innerIdx := key_2_idx(key)
	val := b.set[outIdx]
	if val == 0 {
		return outIdx, val
	}

	val &= ^(1 << innerIdx)
	if val == 0 {
		delete(b.set, outIdx)
	} else {
		b.set[outIdx] = val
	}
	return outIdx, val
}

func (b *BitSet) Flip(key int64) (int64, uint64) {
	outIdx, innerIdx := key_2_idx(key)
	val := b.set[outIdx]
	raw := (val & (1 << innerIdx)) != 0
	if raw {
		val &= ^(1 << innerIdx)
	} else {
		val |= (1 << innerIdx)
	}
	if val == 0 {
		delete(b.set, outIdx)
	} else {
		b.set[outIdx] = val
	}
	return outIdx, val
}

func (b *BitSet) Count() int {
	count := 0
	for _, v := range b.set {
		count += bits.OnesCount64(uint64(v))
	}
	return count
}

func (b *BitSet) Len() int {
	return len(b.set)
}

func (b *BitSet) Any() bool {
	return !b.None()
}

func (b *BitSet) None() bool {
	return len(b.set) == 0
}

func (b *BitSet) Intersection(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSet(nil)
	}

	otherMap := other.set
	var deleteKeys []int64
	for k, v := range b.set {
		if otherV, ok := otherMap[k]; ok {
			newV := v & otherV
			if newV != 0 {
				rlt.set[k] = newV
			} else if isInplace {
				if deleteKeys == nil {
					deleteKeys = make([]int64, 0, 1)
				}
				deleteKeys = append(deleteKeys, k)
			}
		} else if isInplace {
			if deleteKeys == nil {
				deleteKeys = make([]int64, 0, 1)
			}
			deleteKeys = append(deleteKeys, k)
		}
	}

	if isInplace && deleteKeys != nil {
		for _, k := range deleteKeys {
			delete(rlt.set, k)
		}
	}
	return rlt
}

func (b *BitSet) RemoveIntersection(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSet(nil)
	}

	inter := b.Intersection(other)
	var deleteKeys []int64
	for k, v := range b.set {
		if interV, ok := inter.set[k]; ok {
			newVal := v &^ interV
			rlt.set[k] = newVal
			if newVal == 0 && isInplace {
				if deleteKeys == nil {
					deleteKeys = make([]int64, 0, 1)
				}
				deleteKeys = append(deleteKeys, k)
			}
		} else if !isInplace {
			rlt.set[k] = v
		}
	}

	if isInplace && deleteKeys != nil {
		for _, k := range deleteKeys {
			delete(b.set, k)
		}
	}
	return rlt
}

func (b *BitSet) Union(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSet(nil)
	}

	otherMap := other.set
	for k, v := range b.set {
		otherV := otherMap[k]
		if otherV != 0 {
			rlt.set[k] = v | otherV
		} else {
			rlt.set[k] = v
		}
	}

	for k, v := range otherMap {
		if _, ok := b.set[k]; !ok {
			rlt.set[k] = v
		}
	}
	return rlt
}

func (b *BitSet) Iterate(do func(int64) bool) {
	for outIdx, v := range b.set {
		for {
			innerIdx := bits.TrailingZeros64(v)
			if innerIdx == 64 {
				break
			}

			useInner := int64(innerIdx)
			useOuter := outIdx
			if outIdx < 0 {
				useOuter = outIdx + 1
				useInner = -int64(innerIdx)
			}

			if !do(useOuter*unit_bit_len + useInner) {
				break
			}
			v = v & ^(1 << innerIdx)
		}
	}
}

func (b *BitSet) String() string {
	setStr := fmt.Sprintf("%v", b.set)
	nums := []int64{}
	b.Iterate(func(v int64) bool {
		nums = append(nums, v)
		return true
	})
	return fmt.Sprintf("Raw Map: %s => nums: %v", setStr, nums)
}

func key_2_idx(key int64) (outer int64, inner int) {
	outer, inner = key/unit_bit_len, int(key%unit_bit_len)
	if key < 0 {
		outer = outer - 1
		inner = -inner
	}
	return
}
