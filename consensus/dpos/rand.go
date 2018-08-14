// Reference http://xoshiro.di.unimi.it/ xoshiro256**

package dpos

import (
	"sort"
	
	"github.com/themis-network/go-themis/common"
)

// for rand top producers
type sortNum struct {
	serial int
	num uint64
}

type sortNumSlice []*sortNum

func (s sortNumSlice) Len() int {
	return len(s)
}

func (s sortNumSlice) Less(i,j int) bool {
	return s[i].num > s[j].num
}

func (s sortNumSlice) Swap(i,j int) {
	s[i],s[j] = s[j],s[i]
}

func (s sortNumSlice) GetTop(number uint64) sortNumSlice {
	sort.Sort(s)
	return s[0:number]
}

func (s sortNumSlice) Copy() sortNumSlice {
	result := sortNumSlice{}
	for _, v := range s {
		tmp := *v
		result = append(result, &tmp)
	}
	
	return result
}

// High performance random number generator
type Random struct {
	s [4] uint64
}

func NewRandom(seed uint64) *Random{
	res := &Random{}
	res.init(seed)
	return res
}

func (r *Random) Next() uint64 {
	result := rotl(r.s[1] * 5, 7) * 9
	
	t := r.s[1] << 17
	
	r.s[2] ^= r.s[0]
	r.s[3] ^= r.s[1]
	r.s[1] ^= r.s[2]
	r.s[0] ^= r.s[3]
	
	r.s[2] ^= t
	r.s[3] = rotl(r.s[3], 45)
	
	return result
}

// Shuffle topProducers in pseudo-random order
func (r *Random) Shuffle(producers []common.Address) {
	size := len(producers)
	for i, j := uint64(0), uint64(size); i < uint64(len(producers)); i, j = i + 1, j - 1 {
		before := r.Next() % j
		tmp := producers[before]
		producers[before] = producers[i]
		producers[i] = tmp
	}
}

func (r *Random) ResetSeed(seed uint64) {
	r.init(seed)
}

func (r *Random) init(seed uint64) {
	s0 := seed + 0x9e3779b97f4a7c15
	s1 := (s0 ^ (s0 >> 30)) * 0xbf58476d1ce4e5b9
	s2 := (s1 ^ (s1 >> 27)) * 0x94d049bb133111eb
	s3 := s2 ^ (s2 >> 31)
	
	r.s[0] = s0
	r.s[1] = s1
	r.s[2] = s2
	r.s[3] = s3
}


func rotl(x uint64, k uint64) uint64 {
	return (x << k) | (x >> (64 - k))
}