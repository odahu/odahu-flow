package hash

import "github.com/mitchellh/hashstructure"

func Hash(i interface{}) (uint64, error) {
	return hashstructure.Hash(i, nil)
}


func Equal(i interface{}, j interface{}) bool {
	iHash, err := Hash(i)
	if err != nil {
		panic(err)
	}
	jHash, err := Hash(j)
	if err != nil {
		panic(err)
	}
	return iHash == jHash
}
