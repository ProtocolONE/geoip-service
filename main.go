package main

import (
	"github.com/ProtocolONE/geoip-service/pkg"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/kelseyhightower/envconfig"
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/oschwald/geoip2-golang"
	"log"
)

type Config struct {
	GeoIpDbPath    string `envconfig:"MAXMIND_GEOIP_DB_PATH" required:"true"`
	KubernetesHost string `envconfig:"KUBERNETES_SERVICE_HOST" required:"false"`
}

func main() {
	cfg := &Config{}

	if err := envconfig.Process("", cfg); err != nil {
		log.Fatalf("Config init failed with error: %s\n", err)
	}

	db, err := geoip2.Open(cfg.GeoIpDbPath)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	dbMeta := db.Metadata()
	dbInfo := "Loaded database info:\n" +
		"\tFilename: %s\n" +
		"\tVersion: %d.%d\n" +
		"\tType: %s\n"

	log.Printf(dbInfo, cfg.GeoIpDbPath, dbMeta.BinaryFormatMajorVersion, dbMeta.BinaryFormatMinorVersion, dbMeta.DatabaseType)

	var service micro.Service

	options := []micro.Option{
		micro.Name(geoip.ServiceName),
		micro.Version(geoip.Version),
	}

	if cfg.KubernetesHost == "" {
		service = micro.NewService(options...)
		log.Println("Initialize micro service")
	} else {
		service = k8s.NewService(options...)
		log.Println("Initialize k8s service")
	}

	service.Init()

	err = proto.RegisterGeoIpServiceHandler(service.Server(), &geoip.Service{GeoReader: db})
	if err != nil {
		log.Fatal(err)
	}

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
