// Copyright (c) 2019 Sorint.lab S.p.A.
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"errors"
	"fmt"
	"reflect"
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

type alertSimilarTo struct{ al model.Alert }

func (sa *alertSimilarTo) Matches(x interface{}) bool {
	if val, ok := x.(model.Alert); !ok {
		return false
	} else if val.AlertCode != sa.al.AlertCode {
		return false
	} else {
		return reflect.DeepEqual(sa.al.OtherInfo, val.OtherInfo)
	}
}

func (sa *alertSimilarTo) String() string {
	return fmt.Sprintf("is similar to %v", sa.al)
}
