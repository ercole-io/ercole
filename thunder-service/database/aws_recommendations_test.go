// Copyright (c) 2022 Sorint.lab S.p.A.
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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var profileId1 = primitive.NewObjectID()

var privateKeyTest = "privateKey"

var awsProfile4 = model.AwsProfile{
	ID:              profileId1,
	AccessKeyId:     "accessKey4",
	SecretAccessKey: &privateKeyTest,
	Region:          "region4",
	Selected:        true,
	Name:            "profile4",
}

var recAws1 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   999,
	ProfileID:  profileId1,
	Category:   "TestCategory1",
	Suggestion: "Suggestion1",
	Name:       "Name1",
	ResourceID: "ResourceID1",
	ObjectType: "ObjectType1",
	Details: []map[string]interface{}{
		{"NameA": "ValueA"},
		{"NameB": "ValueA"},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
}

func (m *MongodbSuite) TestAddAwsRecommendation_Success() {
	err := m.db.AddAwsObject(recAws1, "aws_recommendations")
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("aws_recommendations").FindOne(context.TODO(), bson.M{
		"profileID": recAws1.ProfileID,
	})
	require.NoError(m.T(), val.Err())

	var out model.AwsRecommendation
	val.Decode(&out)

	assert.Equal(m.T(), recAws1, out)
}

func (m *MongodbSuite) TestGetLastAwsSeqValue_Success() {
	var recs []interface{}
	var result uint64

	recs = append(recs, recAws1)
	err := m.db.AddAwsObjects(recs, "aws_recommendations")
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	result, err = m.db.GetLastAwsSeqValue()
	require.NoError(m.T(), err)

	var expected uint64
	expected = 999

	assert.Equal(m.T(), expected, result)
}

func (m *MongodbSuite) TestGetAwsRecommendations_Success() {
	var results []model.AwsRecommendation

	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	err := m.db.AddAwsObject(recAws1, "aws_recommendations")
	require.NoError(m.T(), err)

	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	err = m.db.AddAwsProfile(awsProfile4)
	require.NoError(m.T(), err)

	results, err = m.db.GetAwsRecommendationsByProfiles([]primitive.ObjectID{profileId1})
	require.NoError(m.T(), err)

	expected := []model.AwsRecommendation{recAws1}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestGetAwsRecommendationsBySeqValue_Success() {
	var results []model.AwsRecommendation

	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	err := m.db.AddAwsObject(recAws1, "aws_recommendations")
	require.NoError(m.T(), err)

	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	err = m.db.AddAwsProfile(awsProfile4)
	require.NoError(m.T(), err)

	results, err = m.db.GetAwsRecommendationsBySeqValue(999)
	require.NoError(m.T(), err)

	expected := []model.AwsRecommendation{recAws1}
	assert.ElementsMatch(m.T(), expected, results)
}
