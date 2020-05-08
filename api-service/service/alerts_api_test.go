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
	"testing"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAlerts_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
		map[string]interface{}{
			"AffectedHosts": 12,
			"Code":          "NEW_SERVER",
			"Count":         12,
			"OldestAlert":   "2020-05-06T15:40:04.543+02:00",
			"Severity":      "NOTICE",
		},
	}

	db.EXPECT().SearchAlerts(
		"aggregated-code-severity",
		[]string{"foo", "bar", "foobarx"}, "AlertCode", true,
		1, 1, model.AlertSeverityMinor, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2020-04-07T14:02:03Z"),
	).Return(
		expectedRes,
		nil,
	).Times(1)

	res, err := as.SearchAlerts(
		"aggregated-code-severity",
		"foo bar foobarx", "AlertCode", true,
		1, 1, model.AlertSeverityMinor, model.AlertStatusNew,
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
		1, 1, model.AlertSeverityMinor, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	).Return(
		nil,
		aerrMock,
	).Times(1)

	res, err := as.SearchAlerts(
		"aggregated-code-severity",
		"foo bar foobarx", "AlertCode", true,
		1, 1, model.AlertSeverityMinor, model.AlertStatusNew,
		utils.P("2019-11-05T14:02:03Z"), utils.P("2019-12-05T14:02:03Z"),
	)

	require.Equal(t, aerrMock, err)
	assert.Nil(t, res)
}

func TestAckAlert_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().UpdateAlertStatus(utils.Str2oid("5e8c234b24f648a08585bd44"), model.AlertStatusAck).Return(nil).Times(1)

	err := as.AckAlert(utils.Str2oid("5e8c234b24f648a08585bd44"))
	require.NoError(t, err)
}

func TestAckAlert_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().UpdateAlertStatus(utils.Str2oid("5e8c234b24f648a08585bd44"), model.AlertStatusAck).Return(aerrMock).Times(1)

	err := as.AckAlert(utils.Str2oid("5e8c234b24f648a08585bd44"))
	require.Equal(t, aerrMock, err)
}
