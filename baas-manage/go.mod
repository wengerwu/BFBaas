module github.com/paybf/baasmanager/baas-gateway

go 1.12

require (
	cloud.google.com/go v0.26.0
	github.com/BurntSushi/toml v0.3.1 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/sessions v0.0.3 // indirect
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/go-xorm/core v0.6.0
	github.com/go-xorm/xorm v0.7.1
	github.com/goinggo/mapstructure v0.0.0-20140717182941-194205d9b4a9
	github.com/google/go-cmp v0.5.0 // indirect
	github.com/hashicorp/hcl v1.0.0 // indirect
	github.com/magiconair/properties v1.8.1 // indirect
	github.com/mitchellh/mapstructure v1.3.2 // indirect
	github.com/onsi/ginkgo v1.13.0 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7 // indirect
	github.com/paybf/baasmanager/baas-core v0.0.0
	github.com/pelletier/go-toml v1.2.0 // indirect
	github.com/shopspring/decimal v0.0.0-20180709203117-cd690d0c9e24
	github.com/spf13/afero v1.1.1 // indirect
	github.com/spf13/cast v1.3.1 // indirect
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.0.2
	github.com/syyongx/ii18n v0.0.0-20190531015407-03d063505fc9
	github.com/typa01/go-utils v0.0.0-20181126045345-a86b05b01c1e
	golang.org/x/crypto v0.0.0-20190404164418-38d8ce5564a5 // indirect
	google.golang.org/appengine v1.1.0 // indirect
)

replace github.com/paybf/baasmanager/baas-core v0.0.0 => ../baas-core
