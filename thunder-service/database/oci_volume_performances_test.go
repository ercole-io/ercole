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
	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var volumePerformance1 model.OciVolumePerformance = model.OciVolumePerformance{
	ID:  utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Vpu: 0,
	Performances: []model.OciPerformance{
		{
			Size: 50,
			Values: model.OciPerfValues{
				MaxThroughput: 0.8,
				MaxIOPS:       100,
			},
		},
		{
			Size: 100,
			Values: model.OciPerfValues{
				MaxThroughput: 1.6,
				MaxIOPS:       200,
			},
		},
	},
}

var volumePerformance2 model.OciVolumePerformance = model.OciVolumePerformance{
	ID:  utils.Str2oid("5dd40bfb12f54dfda7b17777"),
	Vpu: 10,
	Performances: []model.OciPerformance{
		{
			Size: 50,
			Values: model.OciPerfValues{
				MaxThroughput: 24,
				MaxIOPS:       3000,
			},
		},
		{
			Size: 100,
			Values: model.OciPerfValues{
				MaxThroughput: 48,
				MaxIOPS:       6000,
			},
		},
	},
}

var volumePerformanceUpd1 model.OciVolumePerformance = model.OciVolumePerformance{
	ID:  utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Vpu: 0,
	Performances: []model.OciPerformance{
		{
			Size: 50,
			Values: model.OciPerfValues{
				MaxThroughput: 0.9,
				MaxIOPS:       110,
			},
		},
		{
			Size: 100,
			Values: model.OciPerfValues{
				MaxThroughput: 2.6,
				MaxIOPS:       202,
			},
		},
	},
}

func (m *MongodbSuite) TestAddOciVolumePerformance_Success() {
	err := m.db.AddOciVolumePerformance(volumePerformance1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_volume_performance").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_volume_performance").FindOne(context.TODO(), bson.M{
		"_id": volumePerformance1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OciVolumePerformance
	val.Decode(&out)

	assert.Equal(m.T(), volumePerformance1, out)
}

func (m *MongodbSuite) TestGetOciVolumePerformances_Success() {
	var volumePerformances []model.OciVolumePerformance

	err := m.db.AddOciVolumePerformance(volumePerformance1)
	require.NoError(m.T(), err)
	err = m.db.AddOciVolumePerformance(volumePerformance2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_volume_performance").DeleteMany(context.TODO(), bson.M{})
	volumePerformances, err = m.db.GetOciVolumePerformances()
	require.NoError(m.T(), err)

	expected := []model.OciVolumePerformance{
		{
			ID:  utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			Vpu: 0,
			Performances: []model.OciPerformance{
				{
					Size: 50,
					Values: model.OciPerfValues{
						MaxThroughput: 0.8,
						MaxIOPS:       100,
					},
				},
				{
					Size: 100,
					Values: model.OciPerfValues{
						MaxThroughput: 1.6,
						MaxIOPS:       200,
					},
				},
			},
		},
		{
			ID:  utils.Str2oid("5dd40bfb12f54dfda7b17777"),
			Vpu: 10,
			Performances: []model.OciPerformance{
				{
					Size: 50,
					Values: model.OciPerfValues{
						MaxThroughput: 24,
						MaxIOPS:       3000,
					},
				},
				{
					Size: 100,
					Values: model.OciPerfValues{
						MaxThroughput: 48,
						MaxIOPS:       6000,
					},
				},
			},
		},
	}
	assert.EqualValues(m.T(), expected, volumePerformances)
}

func (m *MongodbSuite) TestUpdateOciVolumePerformance_Success() {
	err := m.db.AddOciVolumePerformance(volumePerformance1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_volume_performance").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_volume_performance").FindOne(context.TODO(), bson.M{
		"_id": volumePerformance1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateOciVolumePerformance(volumePerformanceUpd1)
	require.NoError(m.T(), err)

	expected, err := m.db.GetOciVolumePerformances()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[0], volumePerformanceUpd1)
}

func (m *MongodbSuite) TestDeleteOciVolumePerformance() {
	err := m.db.AddOciVolumePerformance(volumePerformance1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_volume_performance").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_volume_performance").FindOne(context.TODO(), bson.M{
		"_id": volumePerformance1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.DeleteOciVolumePerformance(volumePerformance1.ID)
	require.NoError(m.T(), err)

	val = m.db.Client.Database(m.dbname).Collection("oci_volume_performance").FindOne(context.TODO(), bson.M{
		"_id": volumePerformance1.ID,
	})
	require.Error(m.T(), val.Err())
}

func (m *MongodbSuite) TestGetOciVolumePerformance_Success() {
	var volumePerformance *model.OciVolumePerformance

	err := m.db.AddOciVolumePerformance(volumePerformance1)
	require.NoError(m.T(), err)
	err = m.db.AddOciVolumePerformance(volumePerformance2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_volume_performance").DeleteMany(context.TODO(), bson.M{})
	volumePerformance, err = m.db.GetOciVolumePerformance(10, 100)
	require.NoError(m.T(), err)

	fmt.Println("FROM DB: ", volumePerformance)

	expected := model.OciVolumePerformance{
		ID:  utils.Str2oid("5dd40bfb12f54dfda7b17777"),
		Vpu: 10,
		Performances: []model.OciPerformance{
			{
				Size: 100,
				Values: model.OciPerfValues{
					MaxThroughput: 48,
					MaxIOPS:       6000,
				},
			},
		},
	}

	assert.EqualValues(m.T(), expected, *volumePerformance)
}
