package bitset

import (
	"fmt"
	"math/bits"
)

var unit_bit_len = int64(64)

type BitSet struct {
	set IDataSource
}

func NewBitSetFromSource(source IDataSource) *BitSet {
	if source == nil {
		panic("NewBitSetFromSource source is nil")
	}
	return &BitSet{set: source}
}

func (b *BitSet) From(source IDataSource) {
	b.set = source
}

func (b *BitSet) Clone() *BitSet {
	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}
	newSource := iteratable.Clone()
	return NewBitSetFromSource(newSource)
}

func (b *BitSet) Storage() IDataSource {
	return b.set
}

func (b *BitSet) Test(key int64) bool {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	if val == 0 {
		return false
	}
	return (val & (1 << innerIdx)) != 0
}

func (b *BitSet) Set(key int64) (int64, uint64) {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	val |= (1 << innerIdx)
	b.set.Set(outIdx, val)
	return outIdx, val
}

func (b *BitSet) Clear(key int64) (int64, uint64) {
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

func (b *BitSet) Flip(key int64) (int64, uint64) {
	outIdx, innerIdx := key_2_idx(key)
	val, _ := b.set.Get(outIdx)
	raw := (val & (1 << innerIdx)) != 0
	if raw {
		val &= ^(1 << innerIdx)
	} else {
		val |= (1 << innerIdx)
	}
	if val == 0 {
		b.set.Delete(outIdx)
	} else {
		b.set.Set(outIdx, val)
	}
	return outIdx, val
}

func (b *BitSet) Count() int {
	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}
	count := 0
	iteratable.Iterate(func(_ int64, v uint64) bool {
		count += bits.OnesCount64(v)
		return true
	})
	return count
}

func (b *BitSet) Len() int {
	return b.set.Len()
}

func (b *BitSet) Any() bool {
	return !b.None()
}

func (b *BitSet) None() bool {
	return b.set.Len() == 0
}

func (b *BitSet) Intersection(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSetFromSource(iteratable.New())
	}

	otherMap := other.set
	var deleteKeys []int64
	iteratable.Iterate(func(k int64, v uint64) bool {
		if otherV, ok := otherMap.Get(k); ok {
			newV := v & otherV
			if newV != 0 {
				rlt.set.Set(k, newV)
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
		return true
	})

	if isInplace && deleteKeys != nil {
		for _, k := range deleteKeys {
			rlt.set.Delete(k)
		}
	}
	return rlt
}

func (b *BitSet) RemoveIntersection(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSetFromSource(iteratable.New())
	}

	inter := b.Intersection(other)
	var deleteKeys []int64
	iteratable.Iterate(func(k int64, v uint64) bool {
		if interV, ok := inter.set.Get(k); ok {
			newVal := v &^ interV
			rlt.set.Set(k, newVal)
			if newVal == 0 && isInplace {
				if deleteKeys == nil {
					deleteKeys = make([]int64, 0, 1)
				}
				deleteKeys = append(deleteKeys, k)
			}
		} else if !isInplace {
			rlt.set.Set(k, v)
		}
		return true
	})

	if isInplace && deleteKeys != nil {
		for _, k := range deleteKeys {
			b.set.Delete(k)
		}
	}
	return rlt
}

func (b *BitSet) Union(other *BitSet, inplace ...bool) *BitSet {
	isInplace := false
	if len(inplace) > 0 {
		isInplace = inplace[0]
	}

	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}
	otherIter, ok := other.set.(IDataIteratable)
	if !ok {
		panic("BitSet other set is not IDataIteratable")
	}

	var rlt *BitSet
	if isInplace {
		rlt = b
	} else {
		rlt = NewBitSetFromSource(iteratable.New())
	}

	otherMap := other.set
	iteratable.Iterate(func(k int64, v uint64) bool {
		otherV, _ := otherMap.Get(k)
		if otherV != 0 {
			rlt.set.Set(k, v|otherV)
		} else {
			rlt.set.Set(k, v)
		}
		return true
	})

	otherIter.Iterate(func(k int64, v uint64) bool {
		if _, ok := b.set.Get(k); !ok {
			rlt.set.Set(k, v)
		}
		return true
	})
	return rlt
}

func (b *BitSet) Iterate(do func(int64) bool) {
	iteratable, ok := b.set.(IDataIteratable)
	if !ok {
		panic("BitSet set is not IDataIteratable")
	}
	iteratable.Iterate(func(outIdx int64, v uint64) bool {
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
		return true
	})
}

func (b *BitSet) String() string {
	setStr := fmt.Sprintf("%v", b.set)
	nums := []int64{}
	b.Iterate(func(v int64) bool {
		nums = append(nums, v)
		return true
	})
	return fmt.Sprintf("Raw data: %s => setted: %v", setStr, nums)
}

func key_2_idx(key int64) (outer int64, inner int) {
	outer, inner = key/unit_bit_len, int(key%unit_bit_len)
	if key < 0 {
		outer = outer - 1
		inner = -inner
	}
	return
}
