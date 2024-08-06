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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetOraclePsqlMigrabilities_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric := "test_metric"
	schema := "test_schema"
	objectType := "test_objectType"

	expected := []model.PgsqlMigrability{
		{
			Metric:     &metric,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	db.EXPECT().FindPsqlMigrabilities("hostname01", "dbname01").Return(expected, nil)

	res, err := as.GetOraclePsqlMigrabilities("hostname01", "dbname01")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOraclePsqlMigrabilitiesSemaphore_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric1 := "PLSQL LINES"
	metric2 := "NO PLSQL LINES"
	schema := "test_schema"
	objectType := "test_objectType"

	green := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	yellow := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      1000,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	red := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      10001,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	db.EXPECT().FindPsqlMigrabilities("hostname01", "dbname01").Return(green, nil)
	res, err := as.GetOraclePsqlMigrabilitiesSemaphore("hostname01", "dbname01")
	require.NoError(t, err)
	assert.Equal(t, "green", res)

	db.EXPECT().FindPsqlMigrabilities("hostname02", "dbname02").Return(yellow, nil)
	res, err = as.GetOraclePsqlMigrabilitiesSemaphore("hostname02", "dbname02")
	require.NoError(t, err)
	assert.Equal(t, "yellow", res)

	db.EXPECT().FindPsqlMigrabilities("hostname03", "dbname03").Return(red, nil)
	res, err = as.GetOraclePsqlMigrabilitiesSemaphore("hostname03", "dbname03")
	require.NoError(t, err)
	assert.Equal(t, "red", res)
}

func TestGetOraclePdbPsqlMigrabilities_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric := "test_metric"
	schema := "test_schema"
	objectType := "test_objectType"

	expected := []model.PgsqlMigrability{
		{
			Metric:     &metric,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	db.EXPECT().FindPdbPsqlMigrabilities("hostname01", "dbname01", "pdbname01").Return(expected, nil)

	res, err := as.GetOraclePdbPsqlMigrabilities("hostname01", "dbname01", "pdbname01")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOraclePdbPsqlMigrabilitiesSemaphore_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric1 := "PLSQL LINES"
	metric2 := "NO PLSQL LINES"
	schema := "test_schema"
	objectType := "test_objectType"

	green := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	yellow := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      1000,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	red := []model.PgsqlMigrability{
		{
			Metric:     &metric1,
			Count:      10001,
			Schema:     &schema,
			ObjectType: &objectType,
		},
		{
			Metric:     &metric2,
			Count:      0,
			Schema:     &schema,
			ObjectType: &objectType,
		},
	}

	db.EXPECT().FindPdbPsqlMigrabilities("hostname01", "dbname01", "pdbname01").Return(green, nil)
	res, err := as.GetOraclePdbPsqlMigrabilitiesSemaphore("hostname01", "dbname01", "pdbname01")
	require.NoError(t, err)
	assert.Equal(t, "green", res)

	db.EXPECT().FindPdbPsqlMigrabilities("hostname02", "dbname02", "pdbname02").Return(yellow, nil)
	res, err = as.GetOraclePdbPsqlMigrabilitiesSemaphore("hostname02", "dbname02", "pdbname02")
	require.NoError(t, err)
	assert.Equal(t, "yellow", res)

	db.EXPECT().FindPdbPsqlMigrabilities("hostname03", "dbname03", "pdbname03").Return(red, nil)
	res, err = as.GetOraclePdbPsqlMigrabilitiesSemaphore("hostname03", "dbname03", "pdbname03")
	require.NoError(t, err)
	assert.Equal(t, "red", res)
}

func TestListOracleDatabasePsqlMigrabilities(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric := "test_metric"
	schema := "test_schema"
	objectType := "test_objectType"

	expected := []dto.OracleDatabasePgsqlMigrability{
		{
			Hostname: "hostname01",
			Dbname:   "dbname01",
			Flag:     "green",
			Metrics: []model.PgsqlMigrability{
				{
					Metric:     &metric,
					Count:      0,
					Schema:     &schema,
					ObjectType: &objectType,
				},
			},
		},
	}

	db.EXPECT().ListOracleDatabasePsqlMigrabilities().Return(expected, nil)

	res, err := as.ListOracleDatabasePsqlMigrabilities()
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestListOracleDatabasePdbPsqlMigrabilities(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
	}

	metric := "test_metric"
	schema := "test_schema"
	objectType := "test_objectType"

	expected := []dto.OracleDatabasePdbPgsqlMigrability{
		{
			Hostname: "hostname01",
			Dbname:   "dbname01",
			Pdbname:  "pdbname01",
			Flag:     "green",
			Metrics: []model.PgsqlMigrability{
				{
					Metric:     &metric,
					Count:      0,
					Schema:     &schema,
					ObjectType: &objectType,
				},
			},
		},
	}

	db.EXPECT().ListOracleDatabasePdbPsqlMigrabilities().Return(expected, nil)

	res, err := as.ListOracleDatabasePdbPsqlMigrabilities()
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}
