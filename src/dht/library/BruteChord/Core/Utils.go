package Core

import (
	"math/rand/v2"
	"sort"
	"strconv"
)

var usedId = map[ChordHash]bool{}

func toString(A ChordHash) string {
	return strconv.Itoa(int(A))
}

func generateRandomKey() ChordHash {
	return rand.Int64() % (1 << NumberBits)
}

func Sort(ids []ChordHash) {
	sort.Slice(ids, func(i, j int) bool {
		return ids[i] < ids[j]
	})
}

func SortKeys(data Store) []ChordHash {
	keys := make([]ChordHash, 0)
	for key := range data {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func GenerateRandomBinaryId() ChordHash {
	var result ChordHash
	for {
		result = generateRandomKey()
		if _, exist := usedId[result]; !exist {
			usedId[result] = true
			break
		}
	}
	return result
}

// between : starting from L + 1 in a clockwise order, I can reach M before R + 1.
func between(L ChordHash, M ChordHash, R ChordHash) bool {
	if L < R {
		return L < M && M <= R
	} else {
		return L < M || M <= R
	}
}
