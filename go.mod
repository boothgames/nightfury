module github.com/boothgames/nightfury

go 1.12

replace github.com/ugorji/go => github.com/ugorji/go/codec v1.1.7

require (
	github.com/fatih/color v1.7.0
	github.com/gin-gonic/gin v1.9.1
	github.com/golang/mock v1.3.1
	github.com/google/go-cmp v0.5.5
	github.com/mattn/go-colorable v0.1.4 // indirect
	github.com/mitchellh/go-homedir v1.1.0
	github.com/mitchellh/mapstructure v1.1.2
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.8.3
	go.etcd.io/bbolt v1.3.3
	gopkg.in/olahol/melody.v1 v1.0.0-20170518105555-d52139073376
)
