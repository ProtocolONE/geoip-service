package geoip_test

import (
	"context"
	"github.com/ProtocolONE/geoip-service/pkg"
	"github.com/ProtocolONE/geoip-service/pkg/proto"
	"github.com/oschwald/geoip2-golang"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net"
	"io/ioutil"
	"testing"
	pathio "gopkg.in/Clever/pathio.v3"
)

type ServiceTestSuite struct {
	suite.Suite
	service *geoip.Service
	db      *geoip2.Reader
}

func Test_GameService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) SetupTest() {
	pathioReader, err := pathio.Reader("../assets/GeoLite2-City.mmdb")
	if err != nil {
		suite.Fail("Unable to open database file")
	}

	geoipBuf, err := ioutil.ReadAll(pathioReader)
	if err != nil {
		suite.Fail("Unable to read database")
	}

	db, err := geoip2.FromBytes(geoipBuf)	
	if err != nil {
		suite.Fail("Unable to open database")
	}

	suite.db = db
	suite.service = &geoip.Service{GeoReader: db}
}

func (suite *ServiceTestSuite) TearDownTest() {
	if err := suite.db.Close(); err != nil {
		panic(err)
	}
}

func (suite *ServiceTestSuite) TestWrongIpAddress() {
	err := suite.service.GetIpData(context.TODO(), &proto.GeoIpDataRequest{IP: "---"}, &proto.GeoIpDataResponse{})
	assert.Error(suite.T(), err, "Return error with wrong IP")
}

func (suite *ServiceTestSuite) TestRightIpAddress() {
	testIp := "136.0.16.217"
	response := &proto.GeoIpDataResponse{}

	origin, err := suite.db.City(net.ParseIP(testIp))
	assert.NoError(suite.T(), err)

	err = suite.service.GetIpData(context.TODO(), &proto.GeoIpDataRequest{IP: testIp}, response)
	assert.NoError(suite.T(), err)

	assert.EqualValues(suite.T(), response.City.GeoNameID, origin.City.GeoNameID)
	assert.EqualValues(suite.T(), response.City.Names, origin.City.Names)

	assert.EqualValues(suite.T(), response.Continent.GeoNameID, origin.Continent.GeoNameID)
	assert.EqualValues(suite.T(), response.Continent.Code, origin.Continent.Code)
	assert.EqualValues(suite.T(), response.Continent.Names, origin.Continent.Names)

	assert.EqualValues(suite.T(), response.Country.GeoNameID, origin.Country.GeoNameID)
	assert.EqualValues(suite.T(), response.Country.Names, origin.Country.Names)
	assert.EqualValues(suite.T(), response.Country.IsoCode, origin.Country.IsoCode)
	assert.EqualValues(suite.T(), response.Country.IsInEuropeanUnion, origin.Country.IsInEuropeanUnion)

	assert.EqualValues(suite.T(), response.RepresentedCountry.GeoNameID, origin.RepresentedCountry.GeoNameID)
	assert.EqualValues(suite.T(), response.RepresentedCountry.Names, origin.RepresentedCountry.Names)
	assert.EqualValues(suite.T(), response.RepresentedCountry.IsoCode, origin.RepresentedCountry.IsoCode)
	assert.EqualValues(suite.T(), response.RepresentedCountry.IsInEuropeanUnion, origin.RepresentedCountry.IsInEuropeanUnion)
	assert.EqualValues(suite.T(), response.RepresentedCountry.Type, origin.RepresentedCountry.Type)

	assert.EqualValues(suite.T(), response.Postal.Code, origin.Postal.Code)

	assert.EqualValues(suite.T(), response.Location.AccuracyRadius, origin.Location.AccuracyRadius)
	assert.EqualValues(suite.T(), response.Location.Longitude, origin.Location.Longitude)
	assert.EqualValues(suite.T(), response.Location.Latitude, origin.Location.Latitude)
	assert.EqualValues(suite.T(), response.Location.TimeZone, origin.Location.TimeZone)
	assert.EqualValues(suite.T(), response.Location.MetroCode, origin.Location.MetroCode)

	for i, subdiv := range origin.Subdivisions {
		assert.EqualValues(suite.T(), response.Subdivisions[i].GeoNameID, subdiv.GeoNameID)
		assert.EqualValues(suite.T(), response.Subdivisions[i].IsoCode, subdiv.IsoCode)
		assert.EqualValues(suite.T(), response.Subdivisions[i].Names, subdiv.Names)
	}
}
