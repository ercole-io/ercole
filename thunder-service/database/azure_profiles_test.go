// Copyright (c) 2023 Sorint.lab S.p.A.
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

var strClientSecret1 = "ClientSecret1"
var strClientSecret2 = "ClientSecret2"
var strClientSecret3 = "ClientSecret2"
var strClientSecretUpd = "ClientSecretUpd"

var azureProfile1 model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId1",
	ClientId:       "ClientId1",
	SubscriptionId: "SubscriptionId1",
	Region:         "eu-frankfurt-99",
	ClientSecret:   &strClientSecret1,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:       false,
}

var azureProfile2 model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId2",
	ClientId:       "ClientId2",
	SubscriptionId: "SubscriptionId2",
	Region:         "eu-frankfurt-22",
	ClientSecret:   &strClientSecret2,
	ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
	Selected:       false,
}

var azureProfile3 model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId3",
	ClientId:       "ClientId3",
	SubscriptionId: "SubscriptionId3",
	Region:         "eu-frankfurt-33",
	ClientSecret:   &strClientSecret3,
	ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c293"),
	Selected:       true,
}

var azureProfile1UpdWithKey model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId1",
	ClientId:       "ClientId1",
	SubscriptionId: "SubscriptionId1",
	Region:         "eu-frankfurt-77",
	ClientSecret:   &strClientSecretUpd,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:       false,
}

var azureProfile1UpdWithoutKey model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId1",
	ClientId:       "ClientId1",
	SubscriptionId: "SubscriptionId1",
	Region:         "eu-frankfurt-77",
	ClientSecret:   nil,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:       false,
}

var azureProfile1UpdWithoutKeyResult model.AzureProfile = model.AzureProfile{
	TenantId:       "TenantId1",
	ClientId:       "ClientId1",
	SubscriptionId: "SubscriptionId1",
	Region:         "eu-frankfurt-77",
	ClientSecret:   &strClientSecret1,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	Selected:       false,
}

func (m *MongodbSuite) TestInsertAzureProfile_Success() {
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.AzureProfile
	val.Decode(&out)

	assert.Equal(m.T(), azureProfile1, out)
}

func (m *MongodbSuite) TestGetAzureProfilesWithPrivateKey_Success() {
	var profiles []model.AzureProfile

	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAzureProfile(azureProfile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	profiles, err = m.db.GetAzureProfiles(false)
	require.NoError(m.T(), err)

	expected := []model.AzureProfile{
		{
			TenantId:       "TenantId1",
			ClientId:       "ClientId1",
			SubscriptionId: "SubscriptionId1",
			Region:         "eu-frankfurt-99",
			ClientSecret:   &strClientSecret1,
			ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			Selected:       false,
		},
		{
			TenantId:       "TenantId2",
			ClientId:       "ClientId2",
			SubscriptionId: "SubscriptionId2",
			Region:         "eu-frankfurt-22",
			ClientSecret:   &strClientSecret2,
			ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			Selected:       false,
		},
	}
	assert.EqualValues(m.T(), expected, profiles)
}

func (m *MongodbSuite) TestGetAzureProfilesWithoutPrivateKey_Success() {
	var azureProfiles []model.AzureProfile

	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAzureProfile(azureProfile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	azureProfiles, err = m.db.GetAzureProfiles(true)
	require.NoError(m.T(), err)

	expected := []model.AzureProfile{
		{
			TenantId:       "TenantId1",
			ClientId:       "ClientId1",
			SubscriptionId: "SubscriptionId1",
			Region:         "eu-frankfurt-99",
			ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			Selected:       false,
		},
		{
			TenantId:       "TenantId2",
			ClientId:       "ClientId2",
			SubscriptionId: "SubscriptionId2",
			Region:         "eu-frankfurt-22",
			ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			Selected:       false,
		},
	}
	assert.EqualValues(m.T(), expected, azureProfiles)
}

func (m *MongodbSuite) TestUpdateAzureProfileWithPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.AzureProfile
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateAzureProfile(azureProfile1UpdWithKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapAzureProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[azureProfile1UpdWithKey.ID], azureProfile1UpdWithKey)
}

func (m *MongodbSuite) TestUpdateAzureProfileWithoutPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.AzureProfile
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateAzureProfile(azureProfile1UpdWithoutKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapAzureProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[azureProfile1UpdWithKey.ID], azureProfile1UpdWithoutKeyResult)
}

func (m *MongodbSuite) TestDeleteAzureProfile_Success() {
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.DeleteAzureProfile(azureProfile1.ID)
	require.NoError(m.T(), err)

	val = m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})
	require.Error(m.T(), val.Err())
}

func (m *MongodbSuite) TestSelectAzureProfile_Success() {
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})

	err = m.db.SelectAzureProfile(azureProfile1.ID.Hex(), true)
	require.NoError(m.T(), err)

	expected := model.AzureProfile{
		TenantId:       "TenantId1",
		ClientId:       "ClientId1",
		SubscriptionId: "SubscriptionId1",
		Region:         "eu-frankfurt-99",
		ClientSecret:   &strClientSecret1,
		ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
		Selected:       true,
	}

	val := m.db.Client.Database(m.dbname).Collection("azure_profiles").FindOne(context.TODO(), bson.M{
		"_id": azureProfile1.ID,
	})

	require.NoError(m.T(), val.Err())

	var out model.AzureProfile
	val.Decode(&out)

	assert.EqualValues(m.T(), expected, out)
}

func (m *MongodbSuite) TestGetSelectedAzureProfiles_Success() {
	err := m.db.AddAzureProfile(azureProfile1)
	require.NoError(m.T(), err)
	err = m.db.AddAzureProfile(azureProfile2)
	require.NoError(m.T(), err)
	err = m.db.AddAzureProfile(azureProfile3)
	require.NoError(m.T(), err)

	defer m.db.Client.Database(m.dbname).Collection("azure_profiles").DeleteMany(context.TODO(), bson.M{})

	err = m.db.SelectAzureProfile(azureProfile1.ID.Hex(), true)
	require.NoError(m.T(), err)

	expected := []primitive.ObjectID{utils.Str2oid("5dd40bfb12f54dfda7b1c291"), utils.Str2oid("5dd40bfb12f54dfda2b2c293")}

	selectedProfiles, err := m.db.GetSelectedAzureProfiles()
	require.NoError(m.T(), err)

	assert.Equal(m.T(), expected, selectedProfiles)
}
