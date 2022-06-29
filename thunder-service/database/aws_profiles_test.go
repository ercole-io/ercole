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

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

var strSecretAccessKey1 = "SecretAccessTestKey1"
var strSecretAccessKey2 = "SecretAccessTestKey2"
var strSecretAccessKey3 = "SecretAccessTestKey2"
var strSecretAccessKey1Upd = "SecretAccessTestKey1Upd"

var awsProfile1 model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId1",
	Region:          "eu-frankfurt-99",
	SecretAccessKey: &strSecretAccessKey1,
	ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:        false,
}

var awsProfile2 model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId2",
	Region:          "eu-frankfurt-22",
	SecretAccessKey: &strSecretAccessKey2,
	ID:              utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
	Selected:        false,
}

var awsProfile3 model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId3",
	Region:          "eu-frankfurt-33",
	SecretAccessKey: &strSecretAccessKey3,
	ID:              utils.Str2oid("5dd40bfb12f54dfda2b2c293"),
	Selected:        true,
}

var awsProfile1UpdWithKey model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId1",
	Region:          "eu-frankfurt-77",
	SecretAccessKey: &strSecretAccessKey1Upd,
	ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:        false,
}

var awsProfile1UpdWithoutKey model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId1",
	Region:          "eu-frankfurt-77",
	SecretAccessKey: nil,
	ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:        false,
}

var awsProfile1UpdWithoutKeyResult model.AwsProfile = model.AwsProfile{
	AccessKeyId:     "AccessKeyId1",
	Region:          "eu-frankfurt-77",
	SecretAccessKey: &strSecretAccessKey1,
	ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:        false,
}

func (m *MongodbSuite) TestInsertAwsProfile_Success() {
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.AwsProfile
	val.Decode(&out)

	assert.Equal(m.T(), awsProfile1, out)
}

func (m *MongodbSuite) TestGetAwsProfilesWithPrivateKey_Success() {
	var profiles []model.AwsProfile

	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAwsProfile(awsProfile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	profiles, err = m.db.GetAwsProfiles(false)
	require.NoError(m.T(), err)

	expected := []model.AwsProfile{
		{
			ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			AccessKeyId:     "AccessKeyId1",
			Region:          "eu-frankfurt-99",
			SecretAccessKey: &strSecretAccessKey1,
			Selected:        false,
		},
		{
			ID:              utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			AccessKeyId:     "AccessKeyId2",
			Region:          "eu-frankfurt-22",
			SecretAccessKey: &strSecretAccessKey2,
			Selected:        false,
		},
	}
	assert.EqualValues(m.T(), expected, profiles)
}

func (m *MongodbSuite) TestGetAwsProfilesWithoutPrivateKey_Success() {
	var awsProfiles []model.AwsProfile

	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAwsProfile(awsProfile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	awsProfiles, err = m.db.GetAwsProfiles(true)
	require.NoError(m.T(), err)

	expected := []model.AwsProfile{
		{
			ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			AccessKeyId:     "AccessKeyId1",
			Region:          "eu-frankfurt-99",
			SecretAccessKey: nil,
			Selected:        false,
		},
		{
			ID:              utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			AccessKeyId:     "AccessKeyId2",
			Region:          "eu-frankfurt-22",
			SecretAccessKey: nil,
			Selected:        false,
		},
	}
	assert.EqualValues(m.T(), expected, awsProfiles)
}

func (m *MongodbSuite) TestUpdateAwsProfileWithPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.AwsProfile
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateAwsProfile(awsProfile1UpdWithKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapAwsProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[awsProfile1UpdWithKey.ID], awsProfile1UpdWithKey)
}

func (m *MongodbSuite) TestUpdateAwsProfileWithoutPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.AwsProfile
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateAwsProfile(awsProfile1UpdWithoutKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapAwsProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[awsProfile1UpdWithKey.ID], awsProfile1UpdWithoutKeyResult)
}

func (m *MongodbSuite) TestdeleteAwsProfile_Success() {
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.DeleteAwsProfile(awsProfile1.ID)
	require.NoError(m.T(), err)

	val = m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})
	require.Error(m.T(), val.Err())
}

func (m *MongodbSuite) TestSelectAwsProfile_Success() {
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})

	err = m.db.SelectAwsProfile(awsProfile1.ID.Hex(), true)
	require.NoError(m.T(), err)

	expected := model.AwsProfile{
		ID:              utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
		AccessKeyId:     "AccessKeyId1",
		Region:          "eu-frankfurt-99",
		SecretAccessKey: &strSecretAccessKey1,
		Selected:        true,
	}

	val := m.db.Client.Database(m.dbname).Collection("aws_profiles").FindOne(context.TODO(), bson.M{
		"_id": awsProfile1.ID,
	})

	require.NoError(m.T(), val.Err())

	var out model.AwsProfile
	val.Decode(&out)

	assert.EqualValues(m.T(), expected, out)
}

func (m *MongodbSuite) TestGetSelectedAwsAwsProfiles_Success() {
	err := m.db.AddAwsProfile(awsProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAwsProfile(awsProfile2)
	require.NoError(m.T(), err)
	err = m.db.AddAwsProfile(awsProfile3)
	require.NoError(m.T(), err)

	defer m.db.Client.Database(m.dbname).Collection("aws_profiles").DeleteMany(context.TODO(), bson.M{})

	err = m.db.SelectAwsProfile(awsProfile1.ID.Hex(), true)
	require.NoError(m.T(), err)

	expected := []string{"5dd40bfb12f54dfda7b1c291", "5dd40bfb12f54dfda2b2c293"}

	selectedProfiles, err := m.db.GetSelectedAwsProfiles()
	require.NoError(m.T(), err)

	assert.Equal(m.T(), expected, selectedProfiles)
}
