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

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/api-service/apimodel"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestInsertOracleDatabaseAgreement_Success() {
	agg := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "dfdfgdfg",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}
	_, err := m.db.InsertOracleDatabaseAgreement(agg)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("agreements_oracle_database").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("agreements_oracle_database").FindOne(context.TODO(), bson.M{
		"_id": agg.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OracleDatabaseAgreement
	val.Decode(&out)

	assert.Equal(m.T(), agg, out)
}

func (m *MongodbSuite) TestListOracleDatabaseAgreements() {
	defer m.db.Client.Database(m.dbname).Collection("agreements_oracle_database").DeleteMany(context.TODO(), bson.M{})
	agg1 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}
	agg2 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8539"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{},
		ItemDescription: "fgfgd",
		Metrics:         "Computer Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}
	agg3 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed853A"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{},
		ItemDescription: "fgfgd",
		Metrics:         "Named User Plus Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}
	_, err := m.db.InsertOracleDatabaseAgreement(agg1)
	require.NoError(m.T(), err)
	_, err = m.db.InsertOracleDatabaseAgreement(agg2)
	require.NoError(m.T(), err)
	_, err = m.db.InsertOracleDatabaseAgreement(agg3)
	require.NoError(m.T(), err)

	out, err := m.db.ListOracleDatabaseAgreements()
	m.Require().NoError(err)

	assert.Equal(m.T(), []apimodel.OracleDatabaseAgreementsFE{
		{
			ID:          utils.Str2oid("5dcad8933b243f80e2ed8538"),
			AgreementID: "abcde",
			CSI:         "435435",
			CatchAll:    true,
			Count:       345,
			Hosts: []apimodel.OracleDatabaseAgreementsAssociatedHostFE{
				{
					Hostname:                  "foo",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
				},
				{
					Hostname:                  "bar",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
				},
			},
			ItemDescription: "fgfgd",
			Metrics:         "Processor Perpetual",
			PartID:          "678867",
			ReferenceNumber: "567768",
			Unlimited:       true,
			AvailableCount:  345,
			LicensesCount:   345,
			UsersCount:      0,
		},
		{
			ID:              utils.Str2oid("5dcad8933b243f80e2ed8539"),
			AgreementID:     "abcde",
			CSI:             "435435",
			CatchAll:        true,
			Count:           345,
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ItemDescription: "fgfgd",
			Metrics:         "Computer Perpetual",
			PartID:          "678867",
			ReferenceNumber: "567768",
			Unlimited:       true,
			AvailableCount:  345,
			LicensesCount:   345,
			UsersCount:      0,
		},
		{
			ID:              utils.Str2oid("5dcad8933b243f80e2ed853A"),
			AgreementID:     "abcde",
			CSI:             "435435",
			CatchAll:        true,
			Count:           345,
			Hosts:           []apimodel.OracleDatabaseAgreementsAssociatedHostFE{},
			ItemDescription: "fgfgd",
			Metrics:         "Named User Plus Perpetual",
			PartID:          "678867",
			ReferenceNumber: "567768",
			Unlimited:       true,
			AvailableCount:  345,
			LicensesCount:   0,
			UsersCount:      345,
		},
	}, out)

}

func (m *MongodbSuite) TestListOracleDatabaseLicensingObjects() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	out, err := m.db.ListOracleDatabaseLicensingObjects()
	m.Require().NoError(err)

	assert.ElementsMatch(m.T(), []apimodel.OracleDatabaseLicensingObjects{
		{
			LicenseName:   "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			Count:         70,
			OriginalCount: 70,
		},
		{
			LicenseName:   "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			Count:         1.5,
			OriginalCount: 1.5,
		},
		{
			LicenseName:   "Diagnostics Pack",
			Name:          "test-db3",
			Type:          "host",
			Count:         0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseName:   "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			Count:         0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseName:   "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			Count:         70,
			OriginalCount: 70,
		},
	}, out)
}

func (m *MongodbSuite) TestFindOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection("agreements_oracle_database").DeleteMany(context.TODO(), bson.M{})
	agg1 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID:     "abcde",
		CSI:             "435435",
		CatchAll:        true,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		ItemDescription: "fgfgd",
		Metrics:         "Processor Perpetual",
		PartID:          "678867",
		ReferenceNumber: "567768",
		Unlimited:       true,
	}

	_, err := m.db.InsertOracleDatabaseAgreement(agg1)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		out, err := m.db.FindOracleDatabaseAgreement(utils.Str2oid("5dcad8933b243f80e2ed8538"))
		require.NoError(t, err)
		assert.Equal(t, agg1, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		_, err := m.db.FindOracleDatabaseAgreement(utils.Str2oid("5dcad8933b243f80e2ed8539"))
		require.Equal(t, utils.AerrOracleDatabaseAgreementNotFound, err)
	})
}
