module github.com/ProtocolONE/geoip-service

require (
	github.com/InVisionApp/go-health v2.1.0+incompatible
	github.com/InVisionApp/go-logger v1.0.1 // indirect
	github.com/ProtocolONE/go-micro-plugins v0.4.0
	github.com/golang/protobuf v1.5.2
	github.com/kelseyhightower/envconfig v1.3.0
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v0.0.0-20200119172437-4fe21aa238fd
	github.com/oschwald/geoip2-golang v1.2.1
	github.com/oschwald/maxminddb-golang v1.4.0 // indirect
	github.com/prometheus/client_golang v1.3.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/Clever/pathio.v3 v3.7.1
	gopkg.in/DATA-DOG/go-sqlmock.v1 v1.3.0 // indirect
)

replace (
	github.com/asim/go-bson => github.com/paysuper/go-micro-bson v0.0.0-20220702072159-89518495d2a4
	github.com/micro/go-log => github.com/paysuper/go-micro-log v0.0.0-20220702070844-04763368acc8
	github.com/micro/go-micro => github.com/paysuper/go-micro v0.0.0-20220210193104-32a80cb1af1c
	github.com/micro/go-plugins => github.com/paysuper/go-micro-plugins v0.0.0-20220702083743-93bc924f2d9f
	github.com/micro/go-rcache => github.com/paysuper/go-micro-rcache v0.0.0-20220702070444-665c82f4b9d5
	github.com/micro/h2c => github.com/paysuper/go-micro-h2c v0.0.0-20220702065649-c8b8547b076e
	github.com/micro/util => github.com/paysuper/go-micro-util v0.0.0-20220702070652-63b31644d7b0
)

go 1.13
