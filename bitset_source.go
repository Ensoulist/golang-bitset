package bitset

type IDataSource interface {
	Get(int64) (uint64, bool)
	Set(int64, uint64)
	Delete(int64)
	Len() int
}

type IDataIteratable interface {
	IDataSource
	New() IDataIteratable
	Clone() IDataIteratable
	Iterate(func(int64, uint64) bool)
}

type BaseMapSource map[int64]uint64

func (b BaseMapSource) Get(key int64) (uint64, bool) {
	val, ok := b[key]
	return val, ok
}

func (b BaseMapSource) Set(key int64, val uint64) {
	b[key] = val
}

func (b BaseMapSource) Delete(key int64) {
	delete(b, key)
}

func (b BaseMapSource) Len() int {
	return len(b)
}

func (b BaseMapSource) New() IDataIteratable {
	return &BaseMapSource{}
}

func (b BaseMapSource) Clone() IDataIteratable {
	clone := b.New()
	for k, v := range b {
		clone.Set(k, v)
	}
	return clone
}

func (b BaseMapSource) Iterate(fn func(int64, uint64) bool) {
	for k, v := range b {
		if !fn(k, v) {
			break
		}
	}
}

func NewBitSet(mp map[int64]uint64) *BitSet {
	if mp == nil {
		mp = map[int64]uint64{}
	}
	return NewBitSetFromSource(BaseMapSource(mp))
}
