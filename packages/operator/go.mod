module github.com/odahu/odahu-flow/packages/operator

go 1.14

require (
	github.com/DATA-DOG/go-sqlmock v1.5.0
	github.com/Jeffail/gabs v1.4.0 // indirect
	github.com/Masterminds/squirrel v1.4.0
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/aspenmesh/istio-client-go v0.0.0-20190426173040-3e73c27b9ace
	github.com/aws/aws-sdk-go v1.32.11
	github.com/awslabs/amazon-ecr-credential-helper v0.3.1
	github.com/banzaicloud/bank-vaults/pkg/sdk v0.3.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/emicklei/go-restful v2.9.5+incompatible
	github.com/fluent/fluent-logger-golang v1.4.0
	github.com/gin-gonic/gin v1.6.2
	github.com/go-logr/logr v0.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang-migrate/migrate v3.5.4+incompatible
	github.com/google/uuid v1.1.1
	github.com/hashicorp/go-memdb v1.0.4 // indirect
	github.com/hashicorp/vault v1.4.2
	github.com/hashicorp/vault/api v1.0.5-0.20200317185738-82f498082f02
	github.com/hashicorp/vault/sdk v0.1.14-0.20200521112254-72d69106be02 // indirect
	github.com/hashicorp/yamux v0.0.0-20190923154419-df201c70410d // indirect
	github.com/keybase/go-crypto v0.0.0-20190828182435-a05457805304 // indirect
	github.com/lib/pq v1.10.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/hashstructure v1.0.0
	github.com/onsi/gomega v1.10.1
	github.com/ory/dockertest v3.3.5+incompatible
	github.com/pborman/uuid v1.2.0
	github.com/pierrec/lz4 v2.3.0+incompatible // indirect
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.10.0
	github.com/rakyll/statik v0.1.6
	github.com/rclone/rclone v1.53.0
	github.com/satori/go.uuid v1.2.0
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/swag v1.5.0
	github.com/tektoncd/pipeline v0.13.1-0.20200625065359-44f22a067b75
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/vektra/mockery/v2 v2.7.5 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190809123943-df4f5c81cb3b // indirect
	github.com/xeipuuv/gojsonschema v1.1.0
	github.com/zsais/go-gin-prometheus v0.1.0
	go.uber.org/multierr v1.5.0
	go.uber.org/zap v1.15.0
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de
	golang.org/x/net v0.0.0-20200707034311-ab3426394381
	golang.org/x/tools v0.0.0-20200820180210-c8f393745106
	gopkg.in/src-d/go-git.v4 v4.13.1
	gopkg.in/yaml.v2 v2.3.0
	istio.io/api v0.0.0-20200512234804-e5412c253ffe
	k8s.io/api v0.18.7-rc.0
	k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/kubernetes v1.14.7
	knative.dev/networking v0.0.0-20200812200006-4d518e76538a
	knative.dev/pkg v0.0.0-20200812224206-44c860147a87
	knative.dev/serving v0.17.0
	odahu-commons v0.0.0
	sigs.k8s.io/controller-runtime v0.6.1
)

replace (
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	// Use local commons package to release them at the same time
	odahu-commons v0.0.0 => ../commons
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.5.9
)
