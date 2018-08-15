package dpos

import (
    "testing"

    "github.com/themis-network/go-themis/common"
)

const (
    testLength = 26
)

var (
    originalAddress = []common.Address{
        common.BytesToAddress([]byte{0}),
        common.BytesToAddress([]byte{1}),
        common.BytesToAddress([]byte{2}),
        common.BytesToAddress([]byte{3}),
        common.BytesToAddress([]byte{4}),
        common.BytesToAddress([]byte{5}),
        common.BytesToAddress([]byte{6}),
        common.BytesToAddress([]byte{7}),
        common.BytesToAddress([]byte{8}),
        common.BytesToAddress([]byte{9}),
        common.BytesToAddress([]byte{10}),
    }
)

func TestSortNumSlice(t *testing.T) {
    var intArr = []int{100,101,102,103,104,105,106,107,108,109, 110,111,112,113,114,115,116,117,118,119, 120,121,122,123,124,125}
    var uint64Arr = []uint64{10,2,10,2,10,2,10,2,10,2, 1,2,1,2,1,2,1,2,1,2, 0,10,0,10,0,10}
    
    originalTable := sortNumSlice{}
    for index, weight := range uint64Arr {
        originalTable = append(originalTable, &sortNum {
            serial: intArr[index],
            num: weight,
        })
    }
    copyTable := originalTable.Copy()
    
    if !equals(originalTable.GetTop(testLength), copyTable.GetTop(testLength)) {
        t.Errorf("can not get same result with same original data")
    }
}

func TestRandomShuffleEqual(t *testing.T) {
    copyedAddress := make([]common.Address, 0)
    for _, v := range originalAddress {
        copyedAddress = append(copyedAddress, v)
    }
    
    seed := uint64(200)
    random := NewRandom(seed)
    random.Shuffle(copyedAddress)
    
    // The possibility of getting same order with original producers is very small.
    if compareProducers(originalAddress, copyedAddress) {
        t.Errorf("producer's order not changed")
    }
}

func TestRandomShuffle(t *testing.T) {
    copyedAddress := make([]common.Address, 0)
    for _, v := range originalAddress {
        copyedAddress = append(copyedAddress, v)
    }
    
    seed := uint64(100)
    random := NewRandom(seed)
    random.Shuffle(originalAddress)
    
    random.ResetSeed(seed)
    random.Shuffle(copyedAddress)
    
    if !compareProducers(originalAddress, copyedAddress) {
        t.Errorf("can't get same pseudo-random order of producers, first %v, sencond %v", originalAddress, copyedAddress)
    }
}

func BenchmarkRandomShuffle(b *testing.B) {
    random := NewRandom(200)
    
    for i := 0; i < b.N; i++ {
        random.Shuffle(originalAddress)
    }
}

func TestRandomResetSeed(t *testing.T) {
    seed := uint64(1)
    random := NewRandom(seed)
    
    length := 1000
    originalRes := make([]uint64, 0)
    dstRes := make([]uint64, 0)
    
    for i := 0; i < length; i++ {
        originalRes = append(originalRes, random.Next())
    }
    
    random.ResetSeed(seed)
    for i := 0; i < length; i++ {
        dstRes = append(dstRes, random.Next())
    }
    
    for i :=0; i < length; i++ {
        if originalRes[i] != dstRes[i] {
            t.Errorf("reset seed failed: get different result %v, first %v, second %v", i, originalRes[i], dstRes[i])
        }
    }
}

func BenchmarkRandom(b *testing.B) {
    random := NewRandom(0)
    
    for i := 0; i < b.N; i++ {
        random.Next()
    }
}

func equals(src sortNumSlice, dst sortNumSlice) bool {
    if len(src) != len(dst) {
        return false
    }
    
    for i, v := range src {
        if v.serial != dst[i].serial || v.num != dst[i].num {
            return false
        }
    }
    
    return true
}