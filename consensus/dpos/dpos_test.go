package dpos

import (
	"fmt"
	"testing"
	"time"
	"sort"
)

func TestRandom(t *testing.T) {
	testLen := 5
	seed := uint64(time.Now().Unix())
	if seed <= 0 {
		t.Errorf("seed less than zero.")
	}

	rand := NewRandom()
	arr := rand.GenRandomForTest(seed, testLen)

	for i := 0; i < testLen; i++ {
		fmt.Println((*arr)[i])
	}
}

const (
	testLength int = 26
)

func TestSortDeterminacy(t *testing.T) {
	var intArr = [testLength]int{100,101,102,103,104,105,106,107,108,109, 110,111,112,113,114,115,116,117,118,119, 120,121,122,123,124,125}
	var uint64Arr = [testLength]uint64{10,2,10,2,10,2,10,2,10,2, 1,2,1,2,1,2,1,2,1,2, 0,10,0,10,0,10}
	
	var originNums sortNumSlice
	for i := 0; i < testLength ; i++ {
		var tmp = &sortNum {
			serial: intArr[i],
			num: uint64Arr[i],
		}
		originNums = append(originNums, tmp)
	}
	var sortNumsA sortNumSlice
	var sortNumsB sortNumSlice
	sortNumsA = originNums
	sort.Sort(sortNumsA)
	for i := 0; i<testLength; i++ {
		sortNumsB = originNums
		sort.Sort(sortNumsB)
		//fmt.Println(len(sortNumsA), len(sortNumsB))
		for j := 0; j<testLength; j++ {
			if sortNumsA[j].serial != sortNumsB[j].serial {
				t.Errorf("test TestSortDeterminacy() failed.")
			}	
		}
	}
	fmt.Println("TestSortDeterminacy() successed.")
}
