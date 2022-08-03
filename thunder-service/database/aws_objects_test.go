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

var objAws1 model.AwsObject = model.AwsObject{
	SeqValue:    1,
	ProfileID:   profileId1,
	CreatedAt:   time.Date(2022, 5, 26, 0, 0, 1, 0, time.UTC),
	ProfileName: "profile",
	ObjectsCount: []model.AwsObjectCount{
		{Name: "obj1", Count: 1},
	},
}

func (m *MongodbSuite) TestGetAwsObjectsBySeqValue_Success() {
	defer m.db.Client.Database(m.dbname).Collection("aws_objects").DeleteMany(context.TODO(), bson.M{})
	err := m.db.AddAwsObject(objAws1, "aws_objects")
	require.NoError(m.T(), err)

	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	err = m.db.AddAwsProfile(awsProfile4)
	require.NoError(m.T(), err)

	results, err := m.db.GetAwsObjectsBySeqValue(1)
	require.NoError(m.T(), err)

	expected := []model.AwsObject{objAws1}
	assert.ElementsMatch(m.T(), expected, results)
}
