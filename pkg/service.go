package geoip

import (
	"context"
	"fmt"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/oschwald/geoip2-golang"
	"net"
)

const (
	Version     = "latest"
	ServiceName = "p1geoip"
)

type Service struct {
	GeoReader  *geoip2.Reader
	OpsCounter func()
}

func (c *Service) GetIpData(ctx context.Context, req *proto.GeoIpDataRequest, res *proto.GeoIpDataResponse) error {
	ip := net.ParseIP(req.IP)
	if ip == nil {
		return fmt.Errorf("%s is invalid IP adress", req.IP)
	}

	data, err := c.GeoReader.City(ip)
	if err != nil {
		return err
	}

	res.Continent = &proto.GeoIpContinent{
		GeoNameID: uint32(data.Continent.GeoNameID),
		Code:      data.Continent.Code,
		Names:     data.Continent.Names,
	}

	res.City = &proto.GeoIpCity{
		GeoNameID: uint32(data.City.GeoNameID),
		Names:     data.City.Names,
	}

	res.Country = &proto.GeoIpCountry{
		GeoNameID:         uint32(data.Country.GeoNameID),
		IsoCode:           data.Country.IsoCode,
		IsInEuropeanUnion: data.Country.IsInEuropeanUnion,
		Names:             data.Country.Names,
	}

	/*
	   This is used when the IP address belongs to something like a military base. The represented_country is
	   the country that the base represents. This can be useful for managing content licensing, among other uses.
	*/
	res.RepresentedCountry = &proto.GeoIpRepresentedCountry{
		GeoNameID:         uint32(data.RepresentedCountry.GeoNameID),
		IsoCode:           data.RepresentedCountry.IsoCode,
		IsInEuropeanUnion: data.RepresentedCountry.IsInEuropeanUnion,
		Names:             data.RepresentedCountry.Names,
		Type:              data.RepresentedCountry.Type,
	}

	res.Location = &proto.GeoIpLocation{
		AccuracyRadius: uint32(data.Location.AccuracyRadius),
		Latitude:       data.Location.Latitude,
		Longitude:      data.Location.Longitude,
		MetroCode:      uint32(data.Location.MetroCode),
		TimeZone:       data.Location.TimeZone,
	}

	res.Postal = &proto.GeoIpPostal{Code: data.Postal.Code}

	if len(data.Subdivisions) > 0 {
		res.Subdivisions = make([]*proto.GeoIpSubdivision, len(data.Subdivisions))

		for index, subdiv := range data.Subdivisions {
			res.Subdivisions[index] = &proto.GeoIpSubdivision{
				GeoNameID: uint32(subdiv.GeoNameID),
				IsoCode:   subdiv.IsoCode,
				Names:     subdiv.Names,
			}
		}
	}

	c.OpsCounter()
	return nil
}
