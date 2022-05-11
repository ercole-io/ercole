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
	"github.com/ercole-io/ercole/v2/utils"
)

var object1 model.OciObjects = model.OciObjects{
	ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f601"),
	ProfileID: "TestProfile1",
	Objects: []model.OciObject{
		{
			ObjectName:   "objectName1",
			ObjectNumber: 10,
		},
		{
			ObjectName:   "objectName2",
			ObjectNumber: 7,
		}},
	CreatedAt: time.Date(2022, 5, 5, 0, 0, 0, 0, time.UTC),
	Error:     "Error 1",
}

var object2 model.OciObjects = model.OciObjects{
	ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f602"),
	ProfileID: "TestProfile2",
	Objects: []model.OciObject{
		{
			ObjectName:   "objectName3",
			ObjectNumber: 4,
		},
		{
			ObjectName:   "objectName4",
			ObjectNumber: 1,
		}},
	CreatedAt: time.Date(2022, 5, 4, 0, 0, 0, 0, time.UTC),
	Error:     "Error 2",
}

var object3 model.OciObjects = model.OciObjects{
	ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f603"),
	ProfileID: "TestProfile3",
	Objects: []model.OciObject{
		{
			ObjectName:   "objectName5",
			ObjectNumber: 3,
		},
		{
			ObjectName:   "objectName6",
			ObjectNumber: 33,
		},
		{
			ObjectName:   "objectName7",
			ObjectNumber: 333,
		}},
	CreatedAt: time.Date(2022, 4, 5, 0, 0, 0, 0, time.UTC),
	Error:     "Error 3",
}

var object4 model.OciObjects = model.OciObjects{
	ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f604"),
	ProfileID: "TestProfile1",
	Objects: []model.OciObject{
		{
			ObjectName:   "objectName1",
			ObjectNumber: 11,
		},
		{
			ObjectName:   "objectName2",
			ObjectNumber: 77,
		}},
	CreatedAt: time.Date(2022, 5, 2, 0, 0, 0, 0, time.UTC),
	Error:     "Error 1",
}

func (m *MongodbSuite) TestInsertOciObject_Success() {
	err := m.db.AddOciObjects(object1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_objects").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_objects").FindOne(context.TODO(), bson.M{
		"_id": object1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OciObjects
	val.Decode(&out)

	assert.Equal(m.T(), object1, out)
}

func (m *MongodbSuite) TestGetOciObjects_Success() {
	var objects []model.OciObjects

	err := m.db.AddOciObjects(object1)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object2)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object3)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object4)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_objects").DeleteMany(context.TODO(), bson.M{})
	objects, err = m.db.GetOciObjects()
	require.NoError(m.T(), err)

	expected := []model.OciObjects{
		{
			ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f601"),
			ProfileID: "TestProfile1",
			Objects: []model.OciObject{
				{
					ObjectName:   "objectName1",
					ObjectNumber: 10,
				},
				{
					ObjectName:   "objectName2",
					ObjectNumber: 7,
				}},
			CreatedAt: time.Date(2022, 5, 5, 0, 0, 0, 0, time.UTC),
			Error:     "Error 1",
		},
		{
			ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f602"),
			ProfileID: "TestProfile2",
			Objects: []model.OciObject{
				{
					ObjectName:   "objectName3",
					ObjectNumber: 4,
				},
				{
					ObjectName:   "objectName4",
					ObjectNumber: 1,
				}},
			CreatedAt: time.Date(2022, 5, 4, 0, 0, 0, 0, time.UTC),
			Error:     "Error 2",
		},
		{
			ProfileID: "TestProfile3",
			ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f603"),
			Objects: []model.OciObject{
				{
					ObjectName:   "objectName5",
					ObjectNumber: 3,
				},
				{
					ObjectName:   "objectName6",
					ObjectNumber: 33,
				},
				{
					ObjectName:   "objectName7",
					ObjectNumber: 333,
				}},
			CreatedAt: time.Date(2022, 4, 5, 0, 0, 0, 0, time.UTC),
			Error:     "Error 3",
		},
	}
	assert.ElementsMatch(m.T(), expected, objects)
}

func (m *MongodbSuite) TestDeleteOldOciObjects_Success() {
	var objects []model.OciObjects

	err := m.db.AddOciObjects(object1)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object2)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object3)
	require.NoError(m.T(), err)
	err = m.db.AddOciObjects(object4)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_objects").DeleteMany(context.TODO(), bson.M{})

	err = m.db.DeleteOldOciObjects(time.Date(2022, 5, 1, 0, 0, 0, 0, time.UTC))
	require.NoError(m.T(), err)

	objects, err = m.db.GetOciObjects()
	require.NoError(m.T(), err)

	expected := []model.OciObjects{
		{
			ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f601"),
			ProfileID: "TestProfile1",
			Objects: []model.OciObject{
				{
					ObjectName:   "objectName1",
					ObjectNumber: 10,
				},
				{
					ObjectName:   "objectName2",
					ObjectNumber: 7,
				}},
			CreatedAt: time.Date(2022, 5, 5, 0, 0, 0, 0, time.UTC),
			Error:     "Error 1",
		},
		{
			ID:        utils.Str2oid("5ea6a65bb2e36eb58da2f602"),
			ProfileID: "TestProfile2",
			Objects: []model.OciObject{
				{
					ObjectName:   "objectName3",
					ObjectNumber: 4,
				},
				{
					ObjectName:   "objectName4",
					ObjectNumber: 1,
				}},
			CreatedAt: time.Date(2022, 5, 4, 0, 0, 0, 0, time.UTC),
			Error:     "Error 2",
		},
	}
	assert.ElementsMatch(m.T(), expected, objects)
}
