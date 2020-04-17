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

	"github.com/amreo/ercole-services/model"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/amreo/ercole-services/utils"
)

//go:generate mockgen -source ../database/database.go -destination=fake_database.go -package=service
//go:generate mockgen -source service.go -destination=fake_service.go -package=service
//go:generate mockgen -source emailer.go -destination=fake_emailer.go -package=service

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

var emptyHostData model.HostDataMap = model.HostDataMap{
	"Hostname": "",
	"Extra": map[string]interface{}{
		"Databases": primitive.A{},
	},
	"Info": map[string]interface{}{
		"CPUCores": 0,
	},
}

var hostData1 model.HostDataMap = model.HostDataMap{
	"_id":       utils.Str2oid("5dc3f534db7e81a98b726a52"),
	"Hostname":  "superhost1",
	"Archived":  false,
	"CreatedAt": utils.Str2priTime("2019-11-05T14:02:03Z"),
	"Extra": map[string]interface{}{
		"Databases": primitive.A{},
	},
	"Info": map[string]interface{}{
		"CPUCores": 0,
	},
}

var hostData2 model.HostDataMap = model.HostDataMap{
	"_id":       utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	"Hostname":  "superhost1",
	"Archived":  true,
	"CreatedAt": utils.Str2priTime("2019-11-05T12:02:03Z"),
	"Extra": map[string]interface{}{
		"Databases": primitive.A{},
	},
	"Info": map[string]interface{}{
		"CPUCores": 0,
	},
}

var hostData3 model.HostDataMap = model.HostDataMap{
	"_id":       utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	"Hostname":  "superhost1",
	"Archived":  true,
	"CreatedAt": utils.Str2priTime("2019-11-05T16:02:03Z"),
	"Extra": map[string]interface{}{
		"Databases": primitive.A{
			map[string]interface{}{
				"Name":     "acd",
				"Licenses": primitive.A{},
				"Features": primitive.A{},
			},
		},
	},
	"Info": map[string]interface{}{
		"CPUCores": 0,
	},
}

var hostData4 model.HostDataMap = model.HostDataMap{
	"_id":       utils.Str2oid("5dca7a8faebf0b7c2e5daf42"),
	"Hostname":  "superhost1",
	"Archived":  true,
	"CreatedAt": utils.Str2priTime("2019-11-05T18:02:03Z"),
	"Extra": map[string]interface{}{
		"Databases": primitive.A{
			map[string]interface{}{
				"Name": "acd",
				"Licenses": primitive.A{
					map[string]interface{}{
						"Name":  "Oracle ENT",
						"Count": 10,
					},
					map[string]interface{}{
						"Name":  "Driving",
						"Count": 100},
				},
				"Features": primitive.A{
					map[string]interface{}{
						"Name":   "Driving",
						"Status": true},
				},
			},
		},
	},
	"Info": map[string]interface{}{
		"CPUCores": 0,
	},
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
