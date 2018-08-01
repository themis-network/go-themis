package dpos

import (
	"math/big"
)

type set struct {
	size   int
	bucket []list
}

type list struct {
	size  int
	first *node
}
type node struct {
	data []byte
	next *node
}

func (s *set) insert(a []byte) {
	s.size++
	x := big.NewInt(0)
	x.SetBytes(a)
	x.Mod(x, big.NewInt(int64(len(s.bucket))))
	num := x.Uint64()
	if s.bucket[num].size == 0 {
		s.bucket[num].first = &node{a, nil}
		s.bucket[num].size++
	} else {
		temp := s.bucket[num].first
		for i := 0; i != s.bucket[num].size; i++ {
			temp = temp.next
		}
		temp.next = &node{a, nil}
		s.bucket[num].size++
	}
}
func (s *set) find(a []byte) bool {
	x := big.NewInt(0)
	x.SetBytes(a)
	x.Mod(x, big.NewInt(int64(len(s.bucket))))
	num := x.Uint64()
	if s.bucket[num].size == 0 {
		return false
	} else {
		temp := s.bucket[num].first
		for i := 0; i != s.bucket[num].size; i++ {
			if compare(temp.data, a) {
				return true
			}
			temp = temp.next
		}
		return false
	}
}

func compare(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i != len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
