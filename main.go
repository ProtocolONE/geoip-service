package main

import (
	"fmt"
	"github.com/InVisionApp/go-health"
	"github.com/InVisionApp/go-health/handlers"
	"github.com/ProtocolONE/geoip-service/pkg"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/kelseyhightower/envconfig"
	"github.com/micro/go-micro"
	k8s "github.com/micro/kubernetes/go/micro"
	"github.com/oschwald/geoip2-golang"
	"log"
	"net/http"
	"time"
)

type Config struct {
	GeoIpDbPath     string `envconfig:"MAXMIND_GEOIP_DB_PATH" required:"true"`
	KubernetesHost  string `envconfig:"KUBERNETES_SERVICE_HOST" required:"false"`
	HealthCheckPort int    `envconfig:"HEALTH_CHECK_PORT" required:"false" default:"8080"`
}

type customHealthCheck struct{}

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

	go prepareHealthCheck(cfg)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

func prepareHealthCheck(cfg *Config) {
	h := health.New()
	err := h.AddChecks([]*health.Config{
		{
			Name:     "health-check",
			Checker:  &customHealthCheck{},
			Interval: time.Duration(1) * time.Second,
			Fatal:    true,
		},
	})

	if err != nil {
		log.Fatal("Health check register failed")
	}

	log.Printf("Health check listening on :%d", cfg.HealthCheckPort)

	if err = h.Start(); err != nil {
		log.Fatal("Health check start failed")
	}

	http.HandleFunc("/health", handlers.NewJSONHandlerFunc(h, nil))

	if err = http.ListenAndServe(fmt.Sprintf(":%d", cfg.HealthCheckPort), nil); err != nil {
		log.Fatal("Health check listen failed")
	}
}

func (c *customHealthCheck) Status() (interface{}, error) {
	return "ok", nil
}
