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
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func (m *MongodbSuite) TestHistoricizeLicensesCompliance() {
	defer m.db.Client.Database(m.dbname).Collection("oracle_database_licenses_history").DeleteMany(context.TODO(), bson.M{})

	m.db.TimeNow = func() time.Time { return utils.P("2020-12-05T14:02:03+02:00") }

	expectedDateDay1 := utils.PDT("2020-12-05T00:00:00+02:00")
	m.T().Run("First insert, success", func(t *testing.T) {
		licenses := []dto.LicenseCompliance{
			{
				LicenseTypeID:   "L47247",
				ItemDescription: "Oracle Real Application Testing",
				Metric:          "Processor Perpetual",
				Consumed:        0,
				Covered:         0,
				Compliance:      1,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "A90611",
				ItemDescription: "Oracle Database Enterprise Edition",
				Metric:          "Processor Perpetual",
				Consumed:        2.5,
				Covered:         2.5,
				Compliance:      1,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "A90620",
				ItemDescription: "Oracle Partitioning",
				Metric:          "Processor Perpetual",
				Consumed:        3,
				Covered:         3,
				Compliance:      1,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "",
				ItemDescription: "MySQL Enterprise per cluster",
				Metric:          "",
				Consumed:        42,
				Covered:         84,
				Compliance:      0.5,
				Unlimited:       false,
			},
		}
		err := m.db.HistoricizeLicensesCompliance(licenses)
		require.NoError(m.T(), err)

		cur, err := m.db.Client.Database(m.db.Config.Mongodb.DBName).
			Collection("oracle_database_licenses_history").
			Find(context.TODO(), bson.D{})
		require.NoError(m.T(), err)

		ctx := context.TODO()
		defer cur.Close(ctx)

		var actual []map[string]interface{}
		for cur.Next(ctx) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}

			delete(result, "_id")
			actual = append(actual, result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		expected := []map[string]interface{}{
			{"history": primitive.A{map[string]interface{}{"consumed": 0.0, "covered": 0.0, "date": expectedDateDay1}}, "licenseTypeID": "L47247"},
			{"history": primitive.A{map[string]interface{}{"consumed": 2.5, "covered": 2.5, "date": expectedDateDay1}}, "licenseTypeID": "A90611"},
			{"history": primitive.A{map[string]interface{}{"consumed": 3.0, "covered": 3.0, "date": expectedDateDay1}}, "licenseTypeID": "A90620"},
			{"history": primitive.A{map[string]interface{}{"consumed": 42.0, "covered": 84.0, "date": expectedDateDay1}}, "licenseTypeID": "", "itemDescription": "MySQL Enterprise per cluster"},
		}

		assert.ElementsMatch(m.T(), expected, actual)
	})

	m.T().Run("Second insert, next day, success", func(t *testing.T) {
		m.db.TimeNow = func() time.Time { return utils.P("2020-12-06T15:02:03+02:00") }

		licenses := []dto.LicenseCompliance{
			{
				LicenseTypeID:   "L47247",
				ItemDescription: "Oracle Real Application Testing",
				Metric:          "Processor Perpetual",
				Consumed:        0.5,
				Covered:         5,
				Compliance:      1,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "A90611",
				ItemDescription: "Oracle Database Enterprise Edition",
				Metric:          "Processor Perpetual",
				Consumed:        4.5,
				Covered:         2.5,
				Compliance:      0,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "PID001",
				ItemDescription: "Another one",
				Metric:          "",
				Consumed:        3,
				Covered:         3,
				Compliance:      1,
				Unlimited:       false,
			},
			{
				LicenseTypeID:   "",
				ItemDescription: "MySQL Enterprise per cluster",
				Metric:          "",
				Consumed:        43,
				Covered:         86,
				Compliance:      0.5,
				Unlimited:       false,
			},
		}
		err := m.db.HistoricizeLicensesCompliance(licenses)
		require.NoError(m.T(), err)

		cur, err := m.db.Client.Database(m.db.Config.Mongodb.DBName).
			Collection("oracle_database_licenses_history").
			Find(context.TODO(), bson.D{})
		require.NoError(m.T(), err)

		ctx := context.TODO()
		defer cur.Close(ctx)

		var actual []map[string]interface{}
		for cur.Next(ctx) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}

			delete(result, "_id")
			actual = append(actual, result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		expectedDateDay2 := utils.PDT("2020-12-06T00:00:00+02:00")

		expected := []map[string]interface{}{
			{"history": primitive.A{map[string]interface{}{"consumed": 0.0, "covered": 0.0, "date": expectedDateDay1}, map[string]interface{}{"consumed": 0.5, "covered": 5.0, "date": expectedDateDay2}}, "licenseTypeID": "L47247"},
			{"history": primitive.A{map[string]interface{}{"consumed": 2.5, "covered": 2.5, "date": expectedDateDay1}, map[string]interface{}{"consumed": 4.5, "covered": 2.5, "date": expectedDateDay2}}, "licenseTypeID": "A90611"},
			{"history": primitive.A{map[string]interface{}{"consumed": 3.0, "covered": 3.0, "date": expectedDateDay1}}, "licenseTypeID": "A90620"},
			{"history": primitive.A{map[string]interface{}{"consumed": 3.0, "covered": 3.0, "date": expectedDateDay2}}, "licenseTypeID": "PID001"},
			{"history": primitive.A{map[string]interface{}{"consumed": 42.0, "covered": 84.0, "date": expectedDateDay1}, map[string]interface{}{"consumed": 43.0, "covered": 86.0, "date": expectedDateDay2}}, "licenseTypeID": "", "itemDescription": "MySQL Enterprise per cluster"},
		}

		assert.ElementsMatch(m.T(), expected, actual)
	})

	m.T().Run("Third insert, same day, success", func(t *testing.T) {
		m.db.TimeNow = func() time.Time { return utils.P("2020-12-06T15:02:03+02:00") }

		licenses := []dto.LicenseCompliance{
			{
				LicenseTypeID:   "A90611",
				ItemDescription: "Oracle Database Enterprise Edition",
				Metric:          "Processor Perpetual",
				Consumed:        42.52,
				Covered:         2.5,
				Compliance:      0,
				Unlimited:       false,
			},
		}
		err := m.db.HistoricizeLicensesCompliance(licenses)
		require.NoError(m.T(), err)

		cur, err := m.db.Client.Database(m.db.Config.Mongodb.DBName).
			Collection("oracle_database_licenses_history").
			Find(context.TODO(), bson.D{})
		require.NoError(m.T(), err)

		ctx := context.TODO()
		defer cur.Close(ctx)

		var actual []map[string]interface{}
		for cur.Next(ctx) {
			var result map[string]interface{}
			err := cur.Decode(&result)
			if err != nil {
				log.Fatal(err)
			}

			delete(result, "_id")
			actual = append(actual, result)
		}
		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		expectedDateDay2 := utils.PDT("2020-12-06T00:00:00+02:00")

		expected := []map[string]interface{}{
			{"history": primitive.A{map[string]interface{}{"consumed": 0.0, "covered": 0.0, "date": expectedDateDay1}, map[string]interface{}{"consumed": 0.5, "covered": 5.0, "date": expectedDateDay2}}, "licenseTypeID": "L47247"},
			{"history": primitive.A{map[string]interface{}{"consumed": 2.5, "covered": 2.5, "date": expectedDateDay1}, map[string]interface{}{"consumed": 42.52, "covered": 2.5, "date": expectedDateDay2}}, "licenseTypeID": "A90611"},
			{"history": primitive.A{map[string]interface{}{"consumed": 3.0, "covered": 3.0, "date": expectedDateDay1}}, "licenseTypeID": "A90620"},
			{"history": primitive.A{map[string]interface{}{"consumed": 3.0, "covered": 3.0, "date": expectedDateDay2}}, "licenseTypeID": "PID001"},
			{"history": primitive.A{map[string]interface{}{"consumed": 42.0, "covered": 84.0, "date": expectedDateDay1}, map[string]interface{}{"consumed": 43.0, "covered": 86.0, "date": expectedDateDay2}}, "licenseTypeID": "", "itemDescription": "MySQL Enterprise per cluster"},
		}

		assert.ElementsMatch(m.T(), expected, actual)
	})
}
