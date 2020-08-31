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

	"github.com/ercole-io/ercole/config"
	"github.com/ercole-io/ercole/model"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}
	as.LoadOracleDatabaseAgreementPartsList()

	assert.Equal(t, "L10001", as.OracleDatabaseAgreementParts[0].PartID)
	assert.Equal(t, "Oracle Database Enterprise Edition", as.OracleDatabaseAgreementParts[0].ItemDescription)
	assert.Equal(t, "Named User Plus Perpetual", as.OracleDatabaseAgreementParts[0].Metrics)
	assert.Equal(t, []string{"Oracle ENT"}, as.OracleDatabaseAgreementParts[0].Aliases)
	assert.Equal(t, "L103405", as.OracleDatabaseAgreementParts[2].PartID)
	assert.Equal(t, []string{"Oracle STD"}, as.OracleDatabaseAgreementParts[2].Aliases)

	//Known list of metrics check!
	for i, part := range as.OracleDatabaseAgreementParts {
		assert.Contains(t,
			[]string{"Processor Perpetual", "Named User Plus Perpetual", "Stream Perpetual", "Computer Perpetual"},
			part.Metrics,
			"There is a Oracle/Database agreement part with unknown metric #", i, part,
		)
	}
}

func TestGetOracleDatabaseAgreementPartsList_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
		OracleDatabaseAgreementParts: []model.OracleDatabaseAgreementPart{
			{},
		},
	}
	res, err := as.GetOracleDatabaseAgreementPartsList()
	require.NoError(t, err)
	assert.Equal(t, []model.OracleDatabaseAgreementPart{
		{},
	}, res)
}
