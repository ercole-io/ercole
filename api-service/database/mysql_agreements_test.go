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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (m *MongodbSuite) TestAddMySQLAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).DeleteMany(context.TODO(), bson.M{})

	agreement := model.MySQLAgreement{
		ID:               utils.Str2oid("000000000000000000000001"),
		Type:             "type",
		NumberOfLicenses: 43,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	m.T().Run("should_insert", func(t *testing.T) {
		err := m.db.AddMySQLAgreement(agreement)
		m.Require().NoError(err)
	})
}

func (m *MongodbSuite) TestUpdateMySQLAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).DeleteMany(context.TODO(), bson.M{})

	agreement := model.MySQLAgreement{
		ID:               utils.Str2oid("000000000000000000000001"),
		Type:             "type",
		NumberOfLicenses: 43,
		Clusters:         []string{},
		Hosts:            []string{},
	}

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.UpdateMySQLAgreement(agreement)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrNotFound)
	})

	m.T().Run("should_update", func(t *testing.T) {
		_, err := m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).
			InsertOne(
				context.TODO(),
				agreement,
			)
		require.NoError(t, err)

		err = m.db.UpdateMySQLAgreement(agreement)
		assert.NoError(t, err)
	})
}

func (m *MongodbSuite) TestGetMySQLAgreements() {
	m.T().Run("should_load_all", func(t *testing.T) {
		defer m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).DeleteMany(context.TODO(), bson.M{})

		agreements := []model.MySQLAgreement{
			{
				ID:               utils.Str2oid("aaaaaaaaaaaaaaaaaaaaaaaa"),
				Type:             "",
				NumberOfLicenses: 0,
				Clusters:         []string{},
				Hosts:            []string{},
			},
			{
				ID:               utils.Str2oid("bbbbbbbbbbbbbbbbbbbbbbbb"),
				Type:             "",
				NumberOfLicenses: 0,
				Clusters:         []string{},
				Hosts:            []string{},
			},
		}
		agreementsInt := []interface{}{
			agreements[0],
			agreements[1],
		}
		_, err := m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).InsertMany(context.TODO(), agreementsInt)
		require.Nil(m.T(), err)

		actual, err := m.db.GetMySQLAgreements()
		m.Require().NoError(err)

		assert.Equal(t, agreements, actual)
	})

	m.T().Run("should_load_empty", func(t *testing.T) {
		actual, err := m.db.GetMySQLAgreements()
		m.Require().NoError(err)

		agreements := make([]model.MySQLAgreement, 0)
		assert.Equal(t, agreements, actual)
	})
}

func (m *MongodbSuite) TestDeleteMySQLAgreement() {
	defer m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).DeleteMany(context.TODO(), bson.M{})

	id := utils.Str2oid("000000000000000000000001")

	m.T().Run("error not found", func(t *testing.T) {
		err := m.db.DeleteMySQLAgreement(id)
		var aerr *utils.AdvancedError
		assert.ErrorAs(t, err, &aerr)
		assert.ErrorIs(t, aerr.Err, utils.ErrNotFound)
	})

	m.T().Run("should_delete", func(t *testing.T) {
		agreement := model.MySQLAgreement{
			ID:               utils.Str2oid("000000000000000000000001"),
			Type:             "type",
			NumberOfLicenses: 43,
			Clusters:         []string{},
			Hosts:            []string{},
		}
		_, err := m.db.Client.Database(m.dbname).Collection(mySQLAgreementCollection).
			InsertOne(
				context.TODO(),
				agreement,
			)
		require.NoError(t, err)

		err = m.db.DeleteMySQLAgreement(agreement.ID)
		assert.NoError(t, err)
	})
}
