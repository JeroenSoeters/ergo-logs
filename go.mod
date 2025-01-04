module github.com/jeroensoeters/ergo-logs

go 1.23.4

require (
	ergo.services/ergo v1.999.300
	github.com/fsnotify/fsnotify v1.8.0
	github.com/onsi/gomega v1.36.1
	github.com/stretchr/testify v1.10.0
)

replace ergo.services/ergo => github.com/ergo-services/ergo v1.999.300

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/onsi/ginkgo/v2 v2.22.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.32.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
	golang.org/x/text v0.21.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
