package utils

import "github.com/odahu/odahu-flow/packages/operator/pkg/controller/types"


func SpecsEqual(ke types.RuntimeEntity, se types.StorageEntity) (bool, error) {
	kh, err := ke.GetSpecHash()
	if err != nil {
		return false, err
	}
	sh, err := se.GetSpecHash()
	if err != nil {
		return false, err
	}

	return kh == sh, nil

}
