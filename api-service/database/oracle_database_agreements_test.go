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

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/ercole-io/ercole/v2/utils/mongoutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

var licenseTypeSample = model.OracleDatabaseLicenseType{
	ID:              "ID00001",
	ItemDescription: "ItemDesc 1",
	Cost:            42,
	Metric:          model.LicenseTypeMetricProcessorPerpetual,
	Aliases:         []string{},
}

var agreementSample = model.OracleDatabaseAgreement{
	ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
	AgreementID:     "AID001",
	CSI:             "csi001",
	LicenseTypeID:   licenseTypeSample.ID,
	ReferenceNumber: "R00001",
	Unlimited:       true,
	Count:           345,
	CatchAll:        true,
	Restricted:      true, // there shouldn't be CatchAll==true && Restricted && true, this is only for tests
	Hosts:           []string{"foo", "bar"},
}

func (m *MongodbSuite) TestInsertOracleDatabaseAgreement_Success() {
	aerr := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), aerr)
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

	val := m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).FindOne(context.TODO(), bson.M{
		"_id": agreementSample.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OracleDatabaseAgreement
	err := val.Decode(&out)
	assert.NoError(m.T(), err)

	assert.Equal(m.T(), agreementSample, out)
}

func (m *MongodbSuite) TestInsertOracleDatabaseAgreement_DuplicateError() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	err = m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.Error(m.T(), err, "Should not accept two agreements with same ID")
}

func (m *MongodbSuite) TestGetOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreement(agreementSample.ID)
		require.NoError(t, err)
		assert.Equal(t, agreementSample, *out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		out, err := m.db.GetOracleDatabaseAgreement(utils.Str2oid("xxxxxxxxxxxxxxxxxxxxxxxx"))
		require.Nil(t, out)
		require.Equal(t, utils.ErrOracleDatabaseAgreementNotFound, err)
	})
}

func (m *MongodbSuite) TestUpdateOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	m.T().Run("id_exist", func(t *testing.T) {
		agreementSampleUpdated := model.OracleDatabaseAgreement{
			ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
			AgreementID:     "AID001",
			CSI:             "000001",
			LicenseTypeID:   licenseTypeSample.ID,
			ReferenceNumber: "000002",
			Unlimited:       true,
			Count:           345,
			CatchAll:        true,
			Restricted:      false,
			Hosts:           []string{"foo", "bar"},
		}

		err := m.db.UpdateOracleDatabaseAgreement(agreementSampleUpdated)
		require.NoError(t, err)

		val := m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).FindOne(context.TODO(), bson.M{
			"_id": agreementSampleUpdated.ID,
		})
		require.NoError(m.T(), val.Err())

		var out model.OracleDatabaseAgreement
		err2 := val.Decode(&out)
		assert.NoError(t, err2)

		assert.Equal(m.T(), agreementSampleUpdated, out)
	})

	m.T().Run("id_not_exist", func(t *testing.T) {
		agreementSampleUpdated := model.OracleDatabaseAgreement{
			ID: utils.Str2oid("doesn't exist"),
		}
		err := m.db.UpdateOracleDatabaseAgreement(agreementSampleUpdated)

		require.Equal(t, utils.ErrOracleDatabaseAgreementNotFound, err)
	})
}

func (m *MongodbSuite) TestRemoveOracleDatabaseAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

	err := m.db.InsertOracleDatabaseAgreement(agreementSample)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseAgreement(agreementSample.ID)
	require.NoError(m.T(), err)

	err = m.db.RemoveOracleDatabaseAgreement(utils.Str2oid("5dcad8933b243f80e2ed8538"))
	require.Equal(m.T(), utils.ErrOracleDatabaseAgreementNotFound, err)
}

func (m *MongodbSuite) TestListOracleDatabaseAgreements() {
	defer m.db.Client.Database(m.dbname).Collection("oracle_database_license_types").DeleteMany(context.TODO(), bson.M{})
	licenseTypeSample1 := model.OracleDatabaseLicenseType{
		ID:              "ID00001",
		ItemDescription: "ItemDesc 1",
		Cost:            42,
		Metric:          model.LicenseTypeMetricProcessorPerpetual,
		Aliases:         []string{},
	}
	licenseTypeSample2 := model.OracleDatabaseLicenseType{
		ID:              "ID00002",
		ItemDescription: "ItemDesc 2",
		Cost:            24,
		Metric:          model.LicenseTypeMetricNamedUserPlusPerpetual,
		Aliases:         []string{},
	}
	_, err2 := m.db.Client.Database(m.dbname).Collection(oracleDbLicenseTypesCollection).
		InsertMany(context.TODO(), []interface{}{licenseTypeSample1, licenseTypeSample2})
	require.NoError(m.T(), err2)

	agreementSample := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
		AgreementID:     "agr001",
		CSI:             "csi001",
		LicenseTypeID:   licenseTypeSample1.ID,
		ReferenceNumber: "R00001",
		CatchAll:        true,
		Restricted:      false,
		Count:           345,
		Hosts:           []string{"foo", "bar"},
		Unlimited:       true,
	}
	agreementSample2 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
		AgreementID:     "agr002",
		CSI:             "csi002",
		LicenseTypeID:   licenseTypeSample1.ID,
		ReferenceNumber: "R00002",
		CatchAll:        false,
		Restricted:      true,
		Count:           111,
		Hosts:           []string{"pippo", "clarabella"},
		Unlimited:       false,
	}

	agreementSample3 := model.OracleDatabaseAgreement{
		ID:              utils.Str2oid("cccccccccccccccccccccccc"),
		AgreementID:     "agr002",
		CSI:             "csi002",
		LicenseTypeID:   licenseTypeSample2.ID,
		ReferenceNumber: "R00003",
		CatchAll:        false,
		Restricted:      false,
		Count:           222,
		Hosts:           []string{"topolino", "minni"},
		Unlimited:       true,
	}

	m.T().Run("Empty collection", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

		out, err := m.db.ListOracleDatabaseAgreements()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseAgreementFE{}, out)
	})

	m.T().Run("One agreement", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})
		err := m.db.InsertOracleDatabaseAgreement(agreementSample)
		require.NoError(m.T(), err)

		out, err := m.db.ListOracleDatabaseAgreements()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseAgreementFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				AgreementID:              "agr001",
				CSI:                      "csi001",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   model.LicenseTypeMetricProcessorPerpetual,
				ReferenceNumber:          "R00001",
				Unlimited:                true,
				CatchAll:                 true,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "foo"}, {Hostname: "bar"}},
				LicensesPerCore:          345,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 345,
				AvailableLicensesPerUser: 0,
			},
		}, out)
	})

	m.T().Run("Multiple associations agreement-licenseTypes", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(oracleDbAgreementsCollection).DeleteMany(context.TODO(), bson.M{})

		err := m.db.InsertOracleDatabaseAgreement(agreementSample)
		require.NoError(m.T(), err)
		err = m.db.InsertOracleDatabaseAgreement(agreementSample2)
		require.NoError(m.T(), err)
		err = m.db.InsertOracleDatabaseAgreement(agreementSample3)
		require.NoError(m.T(), err)

		out, err := m.db.ListOracleDatabaseAgreements()
		m.Require().NoError(err)

		assert.Equal(m.T(), []dto.OracleDatabaseAgreementFE{
			{
				ID:                       utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				AgreementID:              "agr001",
				CSI:                      "csi001",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   "Processor Perpetual",
				ReferenceNumber:          "R00001",
				Unlimited:                true,
				CatchAll:                 true,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "foo", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "bar", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          345,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 345,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				AgreementID:              "agr002",
				CSI:                      "csi002",
				LicenseTypeID:            "ID00001",
				ItemDescription:          "ItemDesc 1",
				Metric:                   "Processor Perpetual",
				ReferenceNumber:          "R00002",
				Unlimited:                false,
				CatchAll:                 false,
				Restricted:               true,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "pippo", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "clarabella", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          111,
				LicensesPerUser:          0,
				AvailableLicensesPerCore: 111,
				AvailableLicensesPerUser: 0,
			},
			{
				ID:                       utils.Str2oid("cccccccccccccccccccccccc"),
				AgreementID:              "agr002",
				CSI:                      "csi002",
				LicenseTypeID:            "ID00002",
				ItemDescription:          "ItemDesc 2",
				Metric:                   "Named User Plus Perpetual",
				ReferenceNumber:          "R00003",
				Unlimited:                true,
				CatchAll:                 false,
				Restricted:               false,
				Hosts:                    []dto.OracleDatabaseAgreementAssociatedHostFE{{Hostname: "topolino", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}, {Hostname: "minni", CoveredLicensesCount: 0, TotalCoveredLicensesCount: 0, ConsumedLicensesCount: 0}},
				LicensesPerCore:          0,
				LicensesPerUser:          222,
				AvailableLicensesPerCore: 0,
				AvailableLicensesPerUser: 222,
			}},

			out)
	})
}

func (m *MongodbSuite) TestListHostUsingOracleDatabaseLicenses() {
	defer m.db.Client.Database(m.dbname).Collection("hosts").DeleteMany(context.TODO(), bson.M{})
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_07.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_08.json"))
	m.InsertHostData(mongoutils.LoadFixtureMongoHostDataMap(m.T(), "../../fixture/test_apiservice_mongohostdata_17.json"))

	actual, err := m.db.ListHostUsingOracleDatabaseLicenses()
	m.Require().NoError(err)

	expected := []dto.HostUsingOracleDatabaseLicenses{
		{
			LicenseTypeID: "A90649",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
			OriginalCount: 70,
		},
		{
			LicenseTypeID: "A90619",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  1.5,
			OriginalCount: 1.5,
		},
		{
			LicenseTypeID: "A90649",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseTypeID: "A90611",
			Name:          "test-db3",
			Type:          "host",
			LicenseCount:  0.5,
			OriginalCount: 0.5,
		},
		{
			LicenseTypeID: "A90611",
			Name:          "Puzzait",
			Type:          "cluster",
			LicenseCount:  70,
			OriginalCount: 70,
		},
	}

	assert.ElementsMatch(m.T(), expected, actual)
}
