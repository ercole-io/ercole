// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"testing"
)

func TestAddLicensesToSecondaryDbs(t *testing.T) {
	//TOD0
	//	mockCtrl := gomock.NewController(t)
	//	defer mockCtrl.Finish()
	//	mongoDb := NewMockMongoDatabaseInterface(mockCtrl)
	//
	//	hds := HostDataService{
	//		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	//		Database: mongoDb,
	//		Config: config.Configuration{
	//			AlertService: config.AlertService{
	//				PublisherUsername: "publ1sh3r",
	//				PublisherPassword: "M0stS3cretP4ssw0rd",
	//				RemoteEndpoint:    "http://ercole.example.org",
	//			},
	//			DataService: config.DataService{
	//				EnablePatching:       true,
	//				LogInsertingHostdata: true,
	//			},
	//		},
	//		Version: "1.6.6",
	//		Log:     utils.NewLogger("TEST"),
	//	}
	//
	//	hdPrimary := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_22.json")
	//	hdSecondary := utils.LoadFixtureHostData(t, "../../fixture/test_dataservice_hostdata_v1_23.json")
	//
	//	db := hdSecondary.Features.Oracle.Database.Databases[0]
	//	assert.True(t, db.Status == model.OracleDatabaseStatusMounted && db.Role != model.OracleDatabaseRolePrimary)
	//
	//	hds.addLicensesToSecondaryDbs(hdSecondary.Info, &db)
}
