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
)

var recAws1 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   999,
	ProfileID:  "TestProfile1",
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

var recAws2 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   888,
	ProfileID:  "TestProfile2",
	Category:   "TestCategory2",
	Suggestion: "Suggestion2",
	Name:       "Name2",
	ResourceID: "ResourceID2",
	ObjectType: "ObjectType2",
	Details: []map[string]interface{}{
		{"NameC": "ValueC"},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 2, 0, time.UTC),
}

var recAws3 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   777,
	ProfileID:  "TestProfile3",
	Category:   "TestCategory3",
	Suggestion: "Suggestion3",
	Name:       "Name3",
	ResourceID: "ResourceID3",
	ObjectType: "ObjectType3",
	Details: []map[string]interface{}{
		{"NameD": "ValueD"},
		{"NameE": "ValueE"},
		{"NameF": "ValueF"},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 3, 0, time.UTC),
}

var recAws4 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   999,
	ProfileID:  "TestProfile4",
	Category:   "TestCategory4",
	Suggestion: "Suggestion4",
	Name:       "Name4",
	ResourceID: "ResourceID4",
	ObjectType: "ObjectType4",
	Details: []map[string]interface{}{
		{"NameG": "ValueG"},
		{"NameH": "ValueH"},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 4, 0, time.UTC),
}

var recAws5 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   999,
	ProfileID:  "TestProfile5",
	Category:   "TestCategory5",
	Suggestion: "Suggestion5",
	Name:       "Name5",
	ResourceID: "ResourceID5",
	ObjectType: "ObjectType5",
	Details: []map[string]interface{}{
		{"NameI": "ValueI"},
		{"NameL": "ValueL"},
		{"NameM": "ValueM"},
	},
	CreatedAt: time.Date(2022, 5, 20, 0, 0, 5, 0, time.UTC),
}

var recAws6 model.AwsRecommendation = model.AwsRecommendation{
	SeqValue:   999,
	ProfileID:  "TestProfile6",
	Category:   "TestCategory6",
	Suggestion: "Suggestion6",
	Name:       "Name6",
	ResourceID: "ResourceID6",
	ObjectType: "ObjectType6",
	Details: []map[string]interface{}{
		{"NameN": "ValueN"},
	},
	CreatedAt: time.Date(2022, 4, 26, 0, 0, 6, 0, time.UTC),
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

func (m *MongodbSuite) TestAwsOciRecommendations_Success() {
	var recs []interface{}

	recs = append(recs, recAws1, recAws2, recAws3)
	err := m.db.AddAwsObjects(recs, "aws_recommendations")
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	val, err := m.db.Client.Database(m.dbname).Collection("aws_recommendations").Find(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	var results []interface{}

	ctx := context.TODO()
	defer val.Close(ctx)

	for val.Next(ctx) {
		var out model.AwsRecommendation
		err := val.Decode(&out)
		require.NoError(m.T(), err)
		results = append(results, out)
	}
	require.NoError(m.T(), val.Err())

	assert.Equal(m.T(), recs, results)
}

func (m *MongodbSuite) TestGetLastAwsSeqValue_Success() {
	var recs []interface{}
	var result uint64

	recs = append(recs, recAws1, recAws2, recAws3, recAws4, recAws5, recAws6)
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
	var recs []interface{}
	var results []model.AwsRecommendation
	var profiles = []string{"TestProfile1", "TestProfile4", "TestProfile5", "TestProfile6"}

	recs = append(recs, recAws1, recAws2, recAws3, recAws4, recAws5, recAws6)
	err := m.db.AddAwsObjects(recs, "aws_recommendations")
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_recommendations").DeleteMany(context.TODO(), bson.M{})
	results, err = m.db.GetAwsRecommendations(profiles)
	require.NoError(m.T(), err)

	expected := []model.AwsRecommendation{
		{
			SeqValue:   999,
			ProfileID:  "TestProfile1",
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
		},
		{
			SeqValue:   999,
			ProfileID:  "TestProfile4",
			Category:   "TestCategory4",
			Suggestion: "Suggestion4",
			Name:       "Name4",
			ResourceID: "ResourceID4",
			ObjectType: "ObjectType4",
			Details: []map[string]interface{}{
				{"NameG": "ValueG"},
				{"NameH": "ValueH"},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 4, 0, time.UTC),
		},
		{
			SeqValue:   999,
			ProfileID:  "TestProfile5",
			Category:   "TestCategory5",
			Suggestion: "Suggestion5",
			Name:       "Name5",
			ResourceID: "ResourceID5",
			ObjectType: "ObjectType5",
			Details: []map[string]interface{}{
				{"NameI": "ValueI"},
				{"NameL": "ValueL"},
				{"NameM": "ValueM"},
			},
			CreatedAt: time.Date(2022, 5, 20, 0, 0, 5, 0, time.UTC),
		},
		{
			SeqValue:   999,
			ProfileID:  "TestProfile6",
			Category:   "TestCategory6",
			Suggestion: "Suggestion6",
			Name:       "Name6",
			ResourceID: "ResourceID6",
			ObjectType: "ObjectType6",
			Details: []map[string]interface{}{
				{"NameN": "ValueN"},
			},
			CreatedAt: time.Date(2022, 4, 26, 0, 0, 6, 0, time.UTC),
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}
