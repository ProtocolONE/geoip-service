geoip-service
=============

A fast [Go-Micro](https://github.com/micro/go-micro) GRPC based microservice for looking up MaxMind GeoIP2 and GeoLite2 database.

# Prerequisites
Requires a [go installation](https://golang.org/dl/).

A Database (choose one):
* [Free GeoLite 2 database](http://dev.maxmind.com/geoip/geoip2/geolite2/). Download the "MaxMind DB binary, gzipped" and unpack it.
* [GeoIP2 Downloadable Database](http://dev.maxmind.com/geoip/geoip2/downloadable/). This is a more detailed database and should be compatible, but since I don't have access to that, I have been unable to verify this.

## Running the service

This service works as [Go-Micro](https://github.com/micro/go-micro) microservice. You may want to 
setup your own registry with `MICRO_REGISTRY`/`MICRO_REGISTRY_ADDRESS` or use other go-micro flags.  

Download it 

`go get github.com/ProtocolONE/geoip-service`

If you need it uou can rebuild proto file with protoc 
```
protoc --proto_path=. --micro_out=. --go_out=. geoip.proto
```
 
Setup environment variable `MAXMIND_GEOIP_DB_PATH` with path to the maxmind database path.

By default service will be executed with declared by `MICRO_REGISTRY` registry and GRPC as a transport.

## Using Docker
The docker file in this project used to launch geoip-service in Protocol One environment. You may change it in any 
way you need it.

#Using the service

Once the service is running you can use go-micro to make requests

```go
package main

import (
    "context"
    "fmt"
    "github.com/ProtocolONE/geoip-service/pkg"
    "github.com/ProtocolONE/geoip-service/pkg/proto"
    "github.com/micro/go-micro"
)

func main() {
    // create a new service
    service := micro.NewService()

    // parse command line flags
    service.Init()

    // Create new greeter client
    client := proto.NewGeoIpService(geoip.ServiceName, service.Client())

    // Call it
    rsp, err := client.GetIpData(context.TODO(), &proto.GeoIpDataRequest{IP: "8.8.8.8"})
    if err != nil {
        fmt.Println(err)
    }

    // Print response
    fmt.Println(rsp)
}
```