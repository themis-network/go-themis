//reference http://xoshiro.di.unimi.it/ xoshiro256**

package dpos

import (
	"math/big"
	"sort"
	"github.com/themis-network/go-themis/common"
)
//for rand top producers
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

// TODO refactor; add resetSeed
type Random struct {
	s [4] uint64
}

func NewRandom() *Random{
	return &Random{
		s: [4]uint64{0, 0, 0, 0},
	}
}

func (r *Random)setSeed(seed uint64) {
	s0 := seed + 0x9e3779b97f4a7c15
	s1 := (s0 ^ (s0 >> 30)) * 0xbf58476d1ce4e5b9
	s2 := (s1 ^ (s1 >> 27)) * 0x94d049bb133111eb
	s3 := s2 ^ (s2 >> 31)
	r.s = [4]uint64{s0, s1, s2, s3}
}

func rotl(x uint64, k uint64) uint64 {
	return (x << k) | (x >> (64 - k))
}

func (r *Random)genRandom() uint64 {
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

//shffule topProducers in random order 
func (r *Random)Shuffle(producersAddr []common.Address, weightsBig[]*big.Int, amount *big.Int, seed uint64) ([]common.Address, error) {

	var sortWeights sortNumSlice
	for k, v := range weightsBig {
		var tmp = &sortNum {
			serial: k,
			num: (*v).Uint64(),
		}			
		sortWeights = append(sortWeights, tmp)
	}
	sort.Sort(sortWeights)
	var topProducers []common.Address
	var i int64
	for  i = 0; i < amount.Int64(); i++ {
		topProducers = append(topProducers, producersAddr[sortWeights[i].serial])
	}

	//rand top producers
	r.setSeed(seed)
	var randomNums sortNumSlice
	for i = 0; i < amount.Int64(); i++ {
		var tmp = &sortNum {
			serial: int(i),
			num: r.genRandom(),
		}
		randomNums = append(randomNums, tmp)
	}

	sort.Sort(randomNums)

	var newProducers []common.Address
	for i = 0; i < amount.Int64(); i++ {
		newProducers = append(newProducers, topProducers[randomNums[i].serial])
	}

	return newProducers, nil
}

func (r *Random)GenRandomForTest(seed uint64, amount int) *[]uint64{
	r.setSeed(seed)
	res := make([]uint64, amount, amount)
	for i := 0; i < amount; i++ {
		res = append(res, r.genRandom())
	}
	return &res
}