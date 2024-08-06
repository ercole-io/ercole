// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetOracleDatabasePdbPoliciesAuditFlag_Green(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			APIService: config.APIService{
				OracleDatabasePoliciesAudit: []string{
					"policy0",
				},
			},
		},
	}

	expected := map[string][]string{"GREEN": {"policy0"}}

	db.EXPECT().PdbExist("hostname", "dbname", "pdbname").Return(true, nil)

	db.EXPECT().FindOracleDatabasePdbPoliciesAudit("hostname", "dbname", "pdbname").Return(
		&dto.OraclePoliciesAudit{List: []string{"policy0", "policy1"}}, nil)

	res, err := as.GetOracleDatabasePdbPoliciesAuditFlag("hostname", "dbname", "pdbname")

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOracleDatabasePdbPoliciesAuditFlag_Red(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			APIService: config.APIService{
				OracleDatabasePoliciesAudit: []string{
					"policy0",
				},
			},
		},
	}

	expected := map[string][]string{"RED": {"policy0"}}

	db.EXPECT().PdbExist("hostname", "dbname", "pdbname").Return(true, nil)

	db.EXPECT().FindOracleDatabasePdbPoliciesAudit("hostname", "dbname", "pdbname").Return(
		&dto.OraclePoliciesAudit{List: []string{"policy1"}}, nil)

	res, err := as.GetOracleDatabasePdbPoliciesAuditFlag("hostname", "dbname", "pdbname")

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetOracleDatabasePdbPoliciesAuditFlag_NotAvailable(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := map[string][]string{"N/A": nil}

	db.EXPECT().PdbExist("hostname", "dbname", "pdbname").Return(true, nil)

	db.EXPECT().FindOracleDatabasePdbPoliciesAudit("hostname", "dbname", "pdbname").Return(
		&dto.OraclePoliciesAudit{List: []string{"policy1"}}, nil)

	res, err := as.GetOracleDatabasePdbPoliciesAuditFlag("hostname", "dbname", "pdbname")

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestListOracleDatabasePdbPoliciesAudit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := []dto.OraclePdbPoliciesAuditListResponse{
		{
			Hostname:                "hostname",
			DbName:                  "dbname",
			PdbName:                 "pdbname",
			Flag:                    "green",
			PoliciesAuditConfigured: []string{"policy0"},
			PoliciesAudit:           []string{"policy0"},
			Matched:                 []string{"policy0"},
			NotMatched:              []string{},
		},
	}

	db.EXPECT().ListOracleDatabasePdbPoliciesAudit().Return(expected, nil)

	res, err := as.ListOracleDatabasePdbPoliciesAudit()

	require.NoError(t, err)
	assert.Equal(t, expected, res)
}
