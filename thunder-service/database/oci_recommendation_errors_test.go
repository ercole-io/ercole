// Copyright (c) 2021 Sorint.lab S.p.A.
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
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package database

import (
	"context"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
)

var error1 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  999,
	ProfileID: "TestProfile1",
	Category:  "TestCategory1",
	CreatedAt: time.Date(2022, 5, 25, 0, 0, 1, 0, time.UTC),
	Error:     "Error1",
}

var error2 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  888,
	ProfileID: "TestProfile2",
	Category:  "TestCategory2",
	CreatedAt: time.Date(2022, 5, 25, 0, 0, 2, 0, time.UTC),
	Error:     "Error2",
}

var error3 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  777,
	ProfileID: "TestProfile3",
	Category:  "TestCategory3",
	CreatedAt: time.Date(2022, 5, 25, 0, 0, 3, 0, time.UTC),
	Error:     "Error3",
}

var error4 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  999,
	ProfileID: "TestProfile4",
	Category:  "TestCategory4",
	CreatedAt: time.Date(2022, 5, 25, 0, 0, 4, 0, time.UTC),
	Error:     "Error4",
}

var error5 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  999,
	ProfileID: "TestProfile5",
	Category:  "TestCategory5",
	CreatedAt: time.Date(2022, 5, 19, 0, 0, 4, 0, time.UTC),
	Error:     "Error5",
}

var error6 model.OciRecommendationError = model.OciRecommendationError{
	SeqValue:  999,
	ProfileID: "TestProfile6",
	Category:  "TestCategory6",
	CreatedAt: time.Date(2022, 4, 25, 0, 0, 4, 0, time.UTC),
	Error:     "Error6",
}

func (m *MongodbSuite) TestAddOciRecommendationError_Success() {
	err := m.db.AddOciRecommendationError(error1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").FindOne(context.TODO(), bson.M{
		"profileID": error1.ProfileID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OciRecommendationError
	val.Decode(&out)

	assert.Equal(m.T(), error1, out)
}

func (m *MongodbSuite) TestAddOciRecommendationErrors_Success() {
	var errors []model.OciRecommendationError

	errors = append(errors, error1, error2, error3)
	err := m.db.AddOciRecommendationErrors(errors)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").DeleteMany(context.TODO(), bson.M{})
	val, err := m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").Find(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	var results []model.OciRecommendationError

	ctx := context.TODO()
	defer val.Close(ctx)

	for val.Next(ctx) {
		var out model.OciRecommendationError
		err := val.Decode(&out)
		require.NoError(m.T(), err)
		results = append(results, out)
	}
	require.NoError(m.T(), val.Err())

	assert.Equal(m.T(), errors, results)
}

func (m *MongodbSuite) TestGetOciRecommendationErrors_Success() {
	var errors []model.OciRecommendationError
	var results []model.OciRecommendationError
	var strProfiles = []string{"TestProfile1", "TestProfile4"}

	errors = append(errors, error1, error2, error3, error4)
	err := m.db.AddOciRecommendationErrors(errors)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").DeleteMany(context.TODO(), bson.M{})
	results, err = m.db.GetOciRecommendationErrors(strProfiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendationError{
		{
			SeqValue:  999,
			ProfileID: "TestProfile1",
			Category:  "TestCategory1",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 1, 0, time.UTC),
			Error:     "Error1",
		},
		{
			SeqValue:  999,
			ProfileID: "TestProfile4",
			Category:  "TestCategory4",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 4, 0, time.UTC),
			Error:     "Error4",
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestGetOciRecommendationErrorsByProfiles_Success() {
	var errors []model.OciRecommendationError
	var results []model.OciRecommendationError
	var profiles = []string{"TestProfile1", "TestProfile2"}

	errors = append(errors, error1, error2, error3, error4)
	err := m.db.AddOciRecommendationErrors(errors)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").DeleteMany(context.TODO(), bson.M{})
	results, err = m.db.GetOciRecommendationErrorsByProfiles(profiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendationError{
		{
			SeqValue:  999,
			ProfileID: "TestProfile1",
			Category:  "TestCategory1",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 1, 0, time.UTC),
			Error:     "Error1",
		},
		{
			SeqValue:  888,
			ProfileID: "TestProfile2",
			Category:  "TestCategory2",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 2, 0, time.UTC),
			Error:     "Error2",
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestDeleteOldOciRecommendationErrors_Success() {
	var errors []model.OciRecommendationError
	var results []model.OciRecommendationError
	var strProfiles = []string{"TestProfile1", "TestProfile2", "TestProfile3", "TestProfile4"}

	errors = append(errors, error1, error2, error3, error4, error5, error6)
	err := m.db.AddOciRecommendationErrors(errors)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendation_errors").DeleteMany(context.TODO(), bson.M{})

	err = m.db.DeleteOldOciRecommendationErrors(time.Date(2022, 5, 24, 0, 0, 0, 0, time.UTC))
	require.NoError(m.T(), err)

	results, err = m.db.GetOciRecommendationErrors(strProfiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendationError{
		{
			SeqValue:  999,
			ProfileID: "TestProfile1",
			Category:  "TestCategory1",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 1, 0, time.UTC),
			Error:     "Error1",
		},
		{
			SeqValue:  999,
			ProfileID: "TestProfile4",
			Category:  "TestCategory4",
			CreatedAt: time.Date(2022, 5, 25, 0, 0, 4, 0, time.UTC),
			Error:     "Error4",
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}
