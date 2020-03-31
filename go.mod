module github.com/kubermatic/kubermatic-installer

require (
	github.com/Masterminds/semver v1.4.2
	github.com/icza/dyno v0.0.0-20180601094105-0c96289f9585
	github.com/pmezard/go-difflib v1.0.0
	github.com/sirupsen/logrus v1.4.2
	github.com/stretchr/testify v1.4.0
	github.com/urfave/cli v1.20.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.2.8
	k8s.io/api v0.18.0
	k8s.io/apimachinery v0.18.0
	sigs.k8s.io/yaml v1.2.0
)

go 1.13
