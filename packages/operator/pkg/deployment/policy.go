package deployment

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/odahu/odahu-flow/packages/operator/pkg/deployment/bindata"  //nolint
	"text/template"
	"bytes"
)


func ReadDefaultPoliciesAndRender(roleName string)  (map[string]string, error) {
	policies := map[string]string{}

	for _, name := range bindata.AssetNames() {
		bts, err := bindata.Asset(name)
		if err != nil {
			return nil, err
		}

		tpl, err := template.New("_").Parse(string(bts))
		if err != nil {
			return nil, err
		}

		b := bytes.NewBuffer([]byte{})
		err = tpl.Execute(b, struct {
			Role string
		}{
			Role: roleName,
		})
		if err != nil {
			return nil, err
		}

		policies[name] = b.String()
	}
	return policies, nil
}

// Builds default polices for model injected into configmap
func BuildDefaultPolicyConfigMap(cmName string, cmNs string, policies map[string]string) *v1.ConfigMap {
	return &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cmName,
			Namespace: cmNs,
		},
		Data: policies,
	}
}