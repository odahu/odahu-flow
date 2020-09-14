package hash

import "github.com/mitchellh/hashstructure"

func Hash(i interface{}) (uint64, error) {
	return hashstructure.Hash(i, nil)
}
