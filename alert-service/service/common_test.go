package service

import (
	"errors"
	"time"

	"github.com/amreo/ercole-services/model"

	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen -source ../database/database.go -destination=fake_database.go -package=service
//go:generate mockgen -source service.go -destination=fake_service.go -package=service

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

var hostData1 model.HostData = model.HostData{
	ID:        str2oid("5dc3f534db7e81a98b726a52"),
	Hostname:  "superhost1",
	Archived:  false,
	CreatedAt: p("2019-11-05T14:02:03Z"),
}

//p parse the string s and return the equivalent time
func p(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

//btc break the time continuum and return a function that return the time t
func btc(t time.Time) func() time.Time {
	return func() time.Time {
		return t
	}
}

func str2oid(str string) primitive.ObjectID {
	val, _ := primitive.ObjectIDFromHex(str)
	return val
}
