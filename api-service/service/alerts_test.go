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

	"github.com/ercole-io/ercole/v2/config"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestSearchAlerts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "INFO",
		},
		{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "INFO",
		},
	}

	db.EXPECT().SearchAlerts(
		"aggregated-code-severity",
		[]string{"foo", "bar", "foobarx"}, "AlertCode", true,
		1, 1, "", "", model.AlertSeverityCritical, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2020-04-07T14:02:03Z"),
	).Return(
		expectedRes,
		nil,
	).Times(1)

	res, err := as.SearchAlerts(
		"aggregated-code-severity",
		"foo bar foobarx", "AlertCode", true,
		1, 1, "", "", model.AlertSeverityCritical, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2020-04-07T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchAlerts_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchAlerts(
		"aggregated-code-severity",
		[]string{"foo", "bar", "foobarx"}, "AlertCode", true,
		1, 1, "", "", model.AlertSeverityCritical, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(
		nil,
		aerrMock,
	).Times(1)

	res, err := as.SearchAlerts(
		"aggregated-code-severity",
		"foo bar foobarx", "AlertCode", true,
		1, 1, "", "", model.AlertSeverityCritical, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestAckAlerts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().UpdateAlertsStatus([]primitive.ObjectID{utils.Str2oid("5e8c234b24f648a08585bd44")}, model.AlertStatusAck).
		Return(nil).Times(1)

	err := as.AckAlerts([]primitive.ObjectID{utils.Str2oid("5e8c234b24f648a08585bd44")})
	require.NoError(t, err)
}

func TestAckAlerts_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().UpdateAlertsStatus([]primitive.ObjectID{utils.Str2oid("5e8c234b24f648a08585bd44")}, model.AlertStatusAck).
		Return(aerrMock).Times(1)

	err := as.AckAlerts([]primitive.ObjectID{utils.Str2oid("5e8c234b24f648a08585bd44")})
	require.Equal(t, aerrMock, err)
}

func TestAcknowledgeAlerts(t *testing.T) {
	testCases := []struct {
		filter dto.AlertsFilter
		expErr error
	}{
		{
			filter: dto.AlertsFilter{},
			expErr: nil,
		},
		{
			filter: dto.AlertsFilter{},
			expErr: aerrMock,
		},
	}

	for _, tc := range testCases {
		mockCtrl := gomock.NewController(t)
		defer func() {
			mockCtrl.Finish()
		}()

		db := NewMockMongoDatabaseInterface(mockCtrl)
		as := APIService{
			Database: db,
		}

		db.EXPECT().UpdateAlertsStatusByFilter(tc.filter, model.AlertStatusAck).Return(tc.expErr)

		actErr := as.AckAlertsByFilter(tc.filter)
		assert.Equal(t, tc.expErr, actErr)
	}
}

func TestSearchAlertsAsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}

	data := []map[string]interface{}{
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef82"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_LICENSE",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "A new Enterprise license has been enabled to ercsoldbx",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"hostname": "ercsoldbx",
			},
		},
		{
			"_id":                     utils.Str2oid("5f1943c97238d4bb6c98ef83"),
			"alertAffectedTechnology": "Oracle/Database",
			"alertCategory":           "LICENSE",
			"alertCode":               "NEW_OPTION",
			"alertSeverity":           "CRITICAL",
			"alertStatus":             "NEW",
			"date":                    utils.PDT("2020-07-23T10:01:13.746+02:00"),
			"description":             "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server",
			"hostname":                "ercsoldbx",
			"otherInfo": map[string]interface{}{
				"dbname": "ERCSOL19",
				"features": []string{
					"Diagnostics Pack",
				},
				"hostname": "ercsoldbx",
			},
		},
	}

	db.EXPECT().SearchAlerts(
		"all", []string{}, "", false, -1, -1, "Italy", "TST", "", "", utils.P("2020-06-10T11:54:59Z"), utils.P("2020-06-17T11:54:59Z")).
		Return(data, nil).Times(1)

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2019-12-05T14:02:03Z"),
	}

	from := utils.P("2020-06-10T11:54:59Z")
	to := utils.P("2020-06-17T11:54:59Z")

	actual, err := as.SearchAlertsAsXLSX(from, to, filter)
	require.NoError(t, err)
	assert.Equal(t, "LICENSE", actual.GetCellValue("Alerts", "A2"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", actual.GetCellValue("Alerts", "B2"))
	assert.Equal(t, "CRITICAL", actual.GetCellValue("Alerts", "C2"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Alerts", "D2"))
	assert.Equal(t, "NEW_LICENSE", actual.GetCellValue("Alerts", "E2"))
	assert.Equal(t, "A new Enterprise license has been enabled to ercsoldbx", actual.GetCellValue("Alerts", "F2"))

	assert.Equal(t, "LICENSE", actual.GetCellValue("Alerts", "A3"))
	assert.Equal(t, "2020-07-23 08:01:13.746 +0000 UTC", actual.GetCellValue("Alerts", "B3"))
	assert.Equal(t, "CRITICAL", actual.GetCellValue("Alerts", "C3"))
	assert.Equal(t, "ercsoldbx", actual.GetCellValue("Alerts", "D3"))
	assert.Equal(t, "NEW_OPTION", actual.GetCellValue("Alerts", "E3"))
	assert.Equal(t, "The database ERCSOL19 on ercsoldbx has enabled new features (Diagnostics Pack) on server", actual.GetCellValue("Alerts", "F3"))
}
