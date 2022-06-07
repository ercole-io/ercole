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

var rec1 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        999,
	ProfileID:       "TestProfile1",
	Category:        "TestCategory1",
	Suggestion:      "Suggestion1",
	CompartmentID:   "CompartmentID1",
	CompartmentName: "CompartmentName1",
	Name:            "Name1",
	ResourceID:      "ResourceID1",
	ObjectType:      "ObjectType1",
	Details: []model.RecDetail{
		{
			Name:  "NameA",
			Value: "ValueA",
		},
		{
			Name:  "NameB",
			Value: "ValueB",
		},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
}

var rec2 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        888,
	ProfileID:       "TestProfile2",
	Category:        "TestCategory2",
	Suggestion:      "Suggestion2",
	CompartmentID:   "CompartmentID2",
	CompartmentName: "CompartmentName2",
	Name:            "Name2",
	ResourceID:      "ResourceID2",
	ObjectType:      "ObjectType2",
	Details: []model.RecDetail{
		{
			Name:  "NameC",
			Value: "ValueC",
		},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 2, 0, time.UTC),
}

var rec3 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        777,
	ProfileID:       "TestProfile3",
	Category:        "TestCategory3",
	Suggestion:      "Suggestion3",
	CompartmentID:   "CompartmentID3",
	CompartmentName: "CompartmentName3",
	Name:            "Name3",
	ResourceID:      "ResourceID3",
	ObjectType:      "ObjectType3",
	Details: []model.RecDetail{
		{
			Name:  "NameD",
			Value: "ValueD",
		},
		{
			Name:  "NameE",
			Value: "ValueE",
		},
		{
			Name:  "NameF",
			Value: "ValueF",
		},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 3, 0, time.UTC),
}

var rec4 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        999,
	ProfileID:       "TestProfile4",
	Category:        "TestCategory4",
	Suggestion:      "Suggestion4",
	CompartmentID:   "CompartmentID4",
	CompartmentName: "CompartmentName4",
	Name:            "Name4",
	ResourceID:      "ResourceID4",
	ObjectType:      "ObjectType4",
	Details: []model.RecDetail{
		{
			Name:  "NameG",
			Value: "ValueG",
		},
		{
			Name:  "NameH",
			Value: "ValueH",
		},
	},
	CreatedAt: time.Date(2022, 5, 26, 0, 0, 4, 0, time.UTC),
}

var rec5 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        999,
	ProfileID:       "TestProfile5",
	Category:        "TestCategory5",
	Suggestion:      "Suggestion5",
	CompartmentID:   "CompartmentID5",
	CompartmentName: "CompartmentName5",
	Name:            "Name5",
	ResourceID:      "ResourceID5",
	ObjectType:      "ObjectType5",
	Details: []model.RecDetail{
		{
			Name:  "NameI",
			Value: "ValueI",
		},
		{
			Name:  "NameL",
			Value: "ValueL",
		},
		{
			Name:  "NameM",
			Value: "ValueM",
		},
	},
	CreatedAt: time.Date(2022, 5, 20, 0, 0, 5, 0, time.UTC),
}

var rec6 model.OciRecommendation = model.OciRecommendation{
	SeqValue:        999,
	ProfileID:       "TestProfile6",
	Category:        "TestCategory6",
	Suggestion:      "Suggestion6",
	CompartmentID:   "CompartmentID6",
	CompartmentName: "CompartmentName6",
	Name:            "Name6",
	ResourceID:      "ResourceID6",
	ObjectType:      "ObjectType6",
	Details: []model.RecDetail{
		{
			Name:  "NameN",
			Value: "ValueN",
		},
	},
	CreatedAt: time.Date(2022, 4, 26, 0, 0, 6, 0, time.UTC),
}

func (m *MongodbSuite) TestAddOciRecommendation_Success() {
	err := m.db.AddOciRecommendation(rec1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_recommendations").FindOne(context.TODO(), bson.M{
		"profileID": rec1.ProfileID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OciRecommendation
	val.Decode(&out)

	assert.Equal(m.T(), rec1, out)
}

func (m *MongodbSuite) TestAddOciRecommendations_Success() {
	var recs []model.OciRecommendation

	recs = append(recs, rec1, rec2, rec3)
	err := m.db.AddOciRecommendations(recs)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})
	val, err := m.db.Client.Database(m.dbname).Collection("oci_recommendations").Find(context.TODO(), bson.M{})
	require.NoError(m.T(), err)

	var results []model.OciRecommendation

	ctx := context.TODO()
	defer val.Close(ctx)

	for val.Next(ctx) {
		var out model.OciRecommendation
		err := val.Decode(&out)
		require.NoError(m.T(), err)
		results = append(results, out)
	}
	require.NoError(m.T(), val.Err())

	assert.Equal(m.T(), recs, results)
}

func (m *MongodbSuite) TestGetOciRecommendations_Success() {
	var recs []model.OciRecommendation
	var results []model.OciRecommendation
	var profiles = []string{"TestProfile1", "TestProfile4", "TestProfile5", "TestProfile6"}

	recs = append(recs, rec1, rec2, rec3, rec4, rec5, rec6)
	err := m.db.AddOciRecommendations(recs)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})
	results, err = m.db.GetOciRecommendations(profiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendation{
		{
			SeqValue:        999,
			ProfileID:       "TestProfile1",
			Category:        "TestCategory1",
			Suggestion:      "Suggestion1",
			CompartmentID:   "CompartmentID1",
			CompartmentName: "CompartmentName1",
			Name:            "Name1",
			ResourceID:      "ResourceID1",
			ObjectType:      "ObjectType1",
			Details: []model.RecDetail{
				{
					Name:  "NameA",
					Value: "ValueA",
				},
				{
					Name:  "NameB",
					Value: "ValueB",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
		},
		{
			SeqValue:        999,
			ProfileID:       "TestProfile4",
			Category:        "TestCategory4",
			Suggestion:      "Suggestion4",
			CompartmentID:   "CompartmentID4",
			CompartmentName: "CompartmentName4",
			Name:            "Name4",
			ResourceID:      "ResourceID4",
			ObjectType:      "ObjectType4",
			Details: []model.RecDetail{
				{
					Name:  "NameG",
					Value: "ValueG",
				},
				{
					Name:  "NameH",
					Value: "ValueH",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 4, 0, time.UTC),
		},
		{
			SeqValue:        999,
			ProfileID:       "TestProfile5",
			Category:        "TestCategory5",
			Suggestion:      "Suggestion5",
			CompartmentID:   "CompartmentID5",
			CompartmentName: "CompartmentName5",
			Name:            "Name5",
			ResourceID:      "ResourceID5",
			ObjectType:      "ObjectType5",
			Details: []model.RecDetail{
				{
					Name:  "NameI",
					Value: "ValueI",
				},
				{
					Name:  "NameL",
					Value: "ValueL",
				},
				{
					Name:  "NameM",
					Value: "ValueM",
				},
			},
			CreatedAt: time.Date(2022, 5, 20, 0, 0, 5, 0, time.UTC),
		},
		{
			SeqValue:        999,
			ProfileID:       "TestProfile6",
			Category:        "TestCategory6",
			Suggestion:      "Suggestion6",
			CompartmentID:   "CompartmentID6",
			CompartmentName: "CompartmentName6",
			Name:            "Name6",
			ResourceID:      "ResourceID6",
			ObjectType:      "ObjectType6",
			Details: []model.RecDetail{
				{
					Name:  "NameN",
					Value: "ValueN",
				},
			},
			CreatedAt: time.Date(2022, 4, 26, 0, 0, 6, 0, time.UTC),
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestGetOciRecommendationsByProfiles_Success() {
	var recs []model.OciRecommendation
	var results []model.OciRecommendation
	var profiles = []string{"TestProfile1", "TestProfile2"}

	recs = append(recs, rec1, rec2, rec3, rec4)
	err := m.db.AddOciRecommendations(recs)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})
	results, err = m.db.GetOciRecommendationsByProfiles(profiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendation{
		{
			SeqValue:        999,
			ProfileID:       "TestProfile1",
			Category:        "TestCategory1",
			Suggestion:      "Suggestion1",
			CompartmentID:   "CompartmentID1",
			CompartmentName: "CompartmentName1",
			Name:            "Name1",
			ResourceID:      "ResourceID1",
			ObjectType:      "ObjectType1",
			Details: []model.RecDetail{
				{
					Name:  "NameA",
					Value: "ValueA",
				},
				{
					Name:  "NameB",
					Value: "ValueB",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
		},
		{
			SeqValue:        888,
			ProfileID:       "TestProfile2",
			Category:        "TestCategory2",
			Suggestion:      "Suggestion2",
			CompartmentID:   "CompartmentID2",
			CompartmentName: "CompartmentName2",
			Name:            "Name2",
			ResourceID:      "ResourceID2",
			ObjectType:      "ObjectType2",
			Details: []model.RecDetail{
				{
					Name:  "NameC",
					Value: "ValueC",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 2, 0, time.UTC),
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestDeleteOldOciRecommendations_Success() {
	var recs []model.OciRecommendation
	var results []model.OciRecommendation
	var profiles = []string{"TestProfile1", "TestProfile4"}

	recs = append(recs, rec1, rec2, rec3, rec4, rec5, rec6)
	err := m.db.AddOciRecommendations(recs)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})

	err = m.db.DeleteOldOciRecommendations(time.Date(2022, 5, 25, 0, 0, 0, 0, time.UTC))
	require.NoError(m.T(), err)

	results, err = m.db.GetOciRecommendations(profiles)
	require.NoError(m.T(), err)

	expected := []model.OciRecommendation{
		{
			SeqValue:        999,
			ProfileID:       "TestProfile1",
			Category:        "TestCategory1",
			Suggestion:      "Suggestion1",
			CompartmentID:   "CompartmentID1",
			CompartmentName: "CompartmentName1",
			Name:            "Name1",
			ResourceID:      "ResourceID1",
			ObjectType:      "ObjectType1",
			Details: []model.RecDetail{
				{
					Name:  "NameA",
					Value: "ValueA",
				},
				{
					Name:  "NameB",
					Value: "ValueB",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
		},
		{
			SeqValue:        999,
			ProfileID:       "TestProfile4",
			Category:        "TestCategory4",
			Suggestion:      "Suggestion4",
			CompartmentID:   "CompartmentID4",
			CompartmentName: "CompartmentName4",
			Name:            "Name4",
			ResourceID:      "ResourceID4",
			ObjectType:      "ObjectType4",
			Details: []model.RecDetail{
				{
					Name:  "NameG",
					Value: "ValueG",
				},
				{
					Name:  "NameH",
					Value: "ValueH",
				},
			},
			CreatedAt: time.Date(2022, 5, 26, 0, 0, 4, 0, time.UTC),
		},
	}
	assert.ElementsMatch(m.T(), expected, results)
}

func (m *MongodbSuite) TestGetLastOciSeqValue_Success() {
	var recs []model.OciRecommendation
	var result uint64

	recs = append(recs, rec1, rec2, rec3, rec4, rec5, rec6)
	err := m.db.AddOciRecommendations(recs)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_recommendations").DeleteMany(context.TODO(), bson.M{})
	result, err = m.db.GetLastOciSeqValue()
	require.NoError(m.T(), err)

	var expected uint64
	expected = 999

	assert.Equal(m.T(), expected, result)
}
