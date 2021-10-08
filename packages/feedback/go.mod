module github.com/odahu/odahu-flow/packages/feedback

go 1.14

require (
	github.com/bmizerany/assert v0.0.0-20160611221934-b7ed37b82869 // indirect
	github.com/fluent/fluent-logger-golang v1.4.0
	github.com/gin-gonic/gin v1.7.0
	github.com/onsi/ginkgo v1.14.0 // indirect
	github.com/philhofer/fwd v1.1.1 // indirect
	github.com/prometheus/client_golang v1.7.1
	github.com/spf13/cobra v1.0.0
	github.com/spf13/viper v1.7.0
	github.com/stretchr/testify v1.6.1
	github.com/tinylib/msgp v1.1.0 // indirect
	github.com/zsais/go-gin-prometheus v0.1.0
	gopkg.in/yaml.v2 v2.3.0
	odahu-commons v0.0.0
	sigs.k8s.io/controller-runtime v0.6.1
)

// Use local commons package to release them at the same time
replace odahu-commons v0.0.0 => ../commons
