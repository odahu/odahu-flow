package v1

import "github.com/odahu/odahu-flow/packages/operator/pkg/apiserver/workers/v1/types"


func specsEqual(ke types.ExternalServiceEntity, se types.StorageEntity) (bool, error) {
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
