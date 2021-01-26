package deployment

import (
	"bytes"
	"github.com/odahu/odahu-flow/packages/operator/pkg/deployment/bindata" //nolint
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"path"
	"text/template"
)

const mlServerPoliciesDir = "ml_servers/"

func ReadDefaultPoliciesAndRender(roleName string, predictorsPolicyFilename string) (map[string]string, error) {
	policies := map[string]string{}

	for _, name := range bindata.AssetNames() {
		// Filter out policies related to other ML Servers
		dir, file := path.Split(name)
		if dir == mlServerPoliciesDir && file != predictorsPolicyFilename {
			continue
		}

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

		policies[file] = b.String()
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
