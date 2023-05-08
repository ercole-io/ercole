// Copyright (c) 2023 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestSaveExadata_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	asc := NewMockAlertSvcClientInterface(mockCtrl)
	hds := HostDataService{
		Config: config.Configuration{
			AlertService: config.AlertService{
				PublisherUsername: "publ1sh3r",
				PublisherPassword: "M0stS3cretP4ssw0rd",
				RemoteEndpoint:    "http://ercole.example.org"},
			DataService: config.DataService{
				LogInsertingHostdata: false,
			},
		},
		ServerVersion:  "2.34.1",
		Database:       db,
		AlertSvcClient: asc,
		TimeNow:        utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Log:            logger.NewLogger("TEST"),
	}
	exd := mongoutils.LoadFixtureExadata(t, "../../fixture/test_dataservice_exadata_v1_00.json")

	t.Run("New exadata", func(t *testing.T) {
		gomock.InOrder(
			db.EXPECT().FindExadataByRackID(exd.RackID).Return(nil, nil),
			db.EXPECT().AddExadata(exd).Return(nil),
		)

		err := hds.SaveExadata(&exd)
		require.NoError(t, err)
	})

	t.Run("Update exadata", func(t *testing.T) {
		newExd := exd

		newExd.Components = append(newExd.Components, model.OracleExadataComponent{
			HostType: "kvm_host",
			Hostname: "new_hostname",
		})

		gomock.InOrder(
			db.EXPECT().FindExadataByRackID(newExd.RackID).Return(nil, nil),
			db.EXPECT().AddExadata(newExd).Return(nil),
		)

		err := hds.SaveExadata(&newExd)
		require.NoError(t, err)
	})
}
