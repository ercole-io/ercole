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

	"github.com/ercole-io/ercole/api-service/dto"
	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

var partSample = model.OracleDatabasePart{
	PartID:          "ID00001",
	ItemDescription: "ItemDesc 1",
	Cost:            42,
	Metric:          model.AgreementPartMetricProcessorPerpetual,
	Aliases:         []string{},
}

var agreementSample = model.OracleDatabaseAgreement{
	ID:          utils.Str2oid("5dcad8933b243f80e2ed8538"),
	AgreementID: "AID001",
	CSI:         "csi001",
	Parts: []model.AssociatedPart{
		{
			ID:                 utils.Str2oid("5dcad8933b243f80e2ed8551"),
			OracleDatabasePart: partSample,
			ReferenceNumber:    "R00001",
			Unlimited:          true,
			Count:              345,
			CatchAll:           true,
			Hosts:              []string{"foo", "bar"},
		},
	},
}

func (m *MongodbSuite) TestInsertOracleDatabaseAgreement_Success() {

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).FindOne(context.TODO(), bson.M{
		"_id": agreementSample.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OracleDatabaseAgreement
	val.Decode(&out)

	assert.Equal(m.T(), agreementSample, out)
}

func (m *MongodbSuite) TestInsertOracleDatabaseAgreement_DuplicateError() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	_, err = m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.Error(m.T(), err, "Should not accept two agreements with same ID")
}

func (m *MongodbSuite) TestGetOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreement(agreementSample.AgreementID)
		require.NoError(t, err)
		assert.Equal(t, agreementSample, *out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreement("this id doesn't exists")
		require.Nil(t, out)
		require.Equal(t, utils.AerrOracleDatabaseAgreementNotFound, err)
	})
}

func (m *MongodbSuite) TestGetOracleDatabaseAgreementByAssociatedPart() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreementByAssociatedPart(agreementSample.Parts[0].ID)
		require.NoError(t, err)
		assert.Equal(t, agreementSample, *out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreementByAssociatedPart(utils.Str2oid("5dcad8933b243f80e2ed0000"))
		require.Nil(t, out)
		require.Equal(t, utils.AerrOracleDatabaseAgreementNotFound, err)
	})
}

func (m *MongodbSuite) TestUpdateOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		agreementSampleUpdated := model.OracleDatabaseAgreement{
			ID:          utils.Str2oid("5dcad8933b243f80e2ed8538"),
			AgreementID: "AID001",
			CSI:         "000001",
			Parts: []model.AssociatedPart{
				{
					ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
					OracleDatabasePart: partSample,
					ReferenceNumber:    "000002",
					Unlimited:          true,
					Count:              345,
					CatchAll:           true,
					Hosts:              []string{"foo", "bar"},
				},
			},
		}

		err := m.db.UpdateOracleDatabaseAgreement(agreementSampleUpdated)
		require.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).FindOne(context.TODO(), bson.M{
			"_id": agreementSampleUpdated.ID,
		})
		require.NoError(m.T(), val.Err())

		var out model.OracleDatabaseAgreement
		val.Decode(&out)

		assert.Equal(m.T(), agreementSampleUpdated, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		agreementSampleUpdated := model.OracleDatabaseAgreement{
			ID: utils.Str2oid("doesn't exist"),
		}
		err := m.db.UpdateOracleDatabaseAgreement(agreementSampleUpdated)

		require.Equal(t, utils.AerrOracleDatabaseAgreementNotFound, err)
	})
}

func (m *MongodbSuite) TestRemoveOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	out, err := m.db.GetOracleDatabaseAgreement(agreementSample.AgreementID)
	require.NoError(m.T(), err)
	assert.Equal(m.T(), agreementSample, *out)

	err = m.db.RemoveOracleDatabaseAgreement(utils.Str2oid("5dcad8933b243f80e2ed8538"))
	require.NoError(m.T(), err)

	_, err = m.db.GetOracleDatabaseAgreement(agreementSample.AgreementID)
	require.Equal(m.T(), utils.AerrOracleDatabaseAgreementNotFound, err)

	err = m.db.RemoveOracleDatabaseAgreement(utils.Str2oid("5dcad8933b243f80e2ed8538"))
	require.Equal(m.T(), utils.AerrOracleDatabaseAgreementNotFound, err)
}

func (m *MongodbSuite) TestListOracleDatabaseAgreements_OnePart() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})
	agreementSample := model.OracleDatabaseAgreement{
		ID:          utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID: "AID001",
		CSI:         "csi001",
		Parts: []model.AssociatedPart{
			{
				ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				OracleDatabasePart: partSample,
				ReferenceNumber:    "R00001",
				CatchAll:           true,
				Count:              345,
				Hosts:              []string{"foo", "bar"},
				Unlimited:          true,
			}},
	}
	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	out, err := m.db.ListOracleDatabaseAgreements()
	m.Require().NoError(err)

	assert.Equal(m.T(), []dto.OracleDatabaseAgreementFE{
		{
			ID:          utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID: "AID001",
			CSI:         "csi001",

			PartID:          "ID00001",
			ItemDescription: "ItemDesc 1",
			Metric:          model.AgreementPartMetricProcessorPerpetual,

			ReferenceNumber: "R00001",
			Unlimited:       true,
			Count:           345,
			CatchAll:        true,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					Hostname: "foo",
				},
				{
					Hostname: "bar",
				},
			},

			AvailableCount: 345,
			LicensesCount:  345,
			UsersCount:     0,
		},
	}, out)
}

func (m *MongodbSuite) TestListOracleDatabaseAgreements_MultipleParts() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsColl).DeleteMany(context.TODO(), bson.M{})

	agreementSample := model.OracleDatabaseAgreement{
		ID:          utils.Str2oid("5dcad8933b243f80e2ed8538"),
		AgreementID: "agr001",
		CSI:         "csi001",
		Parts: []model.AssociatedPart{
			{

				ID:                 utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				OracleDatabasePart: partSample,
				ReferenceNumber:    "R00001",
				CatchAll:           true,
				Count:              345,
				Hosts:              []string{"foo", "bar"},
				Unlimited:          true,
			}},
	}
	agreementSample2 := model.OracleDatabaseAgreement{
		ID:          utils.Str2oid("5dcad8933b243f80e2ed8539"),
		AgreementID: "agr002",
		CSI:         "csi002",
		Parts: []model.AssociatedPart{
			{
				ID:                 utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				OracleDatabasePart: partSample,
				ReferenceNumber:    "R00002",
				CatchAll:           true,
				Count:              111,
				Hosts:              []string{"pippo", "clarabella"},
				Unlimited:          false,
			},
			{
				ID: utils.Str2oid("cccccccccccccccccccccccc"),
				OracleDatabasePart: model.OracleDatabasePart{
					PartID:          "partID2",
					ItemDescription: "partDescr2",
					Metric:          model.AgreementPartMetricNamedUserPlusPerpetual,
					Cost:            24,
					Aliases:         []string{"Prova", "Aliases"},
				},
				ReferenceNumber: "R00003",
				CatchAll:        false,
				Count:           222,
				Hosts:           []string{"topolino", "minni"},
				Unlimited:       true,
			},
		},
	}

	_, err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)
	_, err = m.db.InsertOracleDatabaseAgreement(agreementSample2)
	require.NoError(m.T(), err)

	out, err := m.db.ListOracleDatabaseAgreements()
	m.Require().NoError(err)

	assert.Equal(m.T(), []dto.OracleDatabaseAgreementFE{
		{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:     "agr001",
			CSI:             "csi001",
			PartID:          "ID00001",
			ItemDescription: "ItemDesc 1",
			Metric:          "Processor Perpetual",
			ReferenceNumber: "R00001",
			Unlimited:       true,
			Count:           345,
			CatchAll:        true,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					Hostname:                  "foo",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0},
				{
					Hostname:                  "bar",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				},
			},
			AvailableCount: 345,
			LicensesCount:  345,
			UsersCount:     0},
		{
			ID:              utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
			AgreementID:     "agr002",
			CSI:             "csi002",
			PartID:          "ID00001",
			ItemDescription: "ItemDesc 1",
			Metric:          "Processor Perpetual",
			ReferenceNumber: "R00002",
			Unlimited:       false,
			Count:           111,
			CatchAll:        true,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{

					Hostname:                  "pippo",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				},
				{
					Hostname:                  "clarabella",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				}},
			AvailableCount: 111,
			LicensesCount:  111,
			UsersCount:     0,
		},
		{
			ID:              utils.Str2oid("cccccccccccccccccccccccc"),
			AgreementID:     "agr002",
			CSI:             "csi002",
			PartID:          "partID2",
			ItemDescription: "partDescr2",
			Metric:          "Named User Plus Perpetual",
			ReferenceNumber: "R00003",
			Unlimited:       true,
			Count:           222,
			CatchAll:        false,
			Hosts: []dto.OracleDatabaseAgreementAssociatedHostFE{
				{
					Hostname:                  "topolino",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				},
				{
					Hostname:                  "minni",
					CoveredLicensesCount:      0,
					TotalCoveredLicensesCount: 0,
					ConsumedLicensesCount:     0,
				},
			},
			AvailableCount: 222,
			LicensesCount:  0,
			UsersCount:     222}},

		out)
}

func (m *MongodbSuite) TestListHostUsingOracleDatabaseLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(utils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	out, err := m.db.ListHostUsingOracleDatabaseLicenses()
	m.Require().NoError(err)

	assert.ElementsMatch(m.T(), []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseName:   "Diagnostics Pack",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
			OriginalCount: 70,
		},
		{
			LicenseName:   "Real Application Clusters",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
			OriginalCount: 1.5,
		},
		{
			LicenseName:   "Diagnostics Pack",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseName:   "Oracle ENT",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseName:   "Oracle ENT",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
			OriginalCount: 70,
		},
	}, out)
}