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

var strPrivateTestKey1 = "PrivateTestKey1"
var strPrivateTestKey2 = "PrivateTestKey2"
var strPrivateTestKey1Upd = "PrivateTestKey1Upd"

var profile1 model.OciProfile = model.OciProfile{
	Profile:        "TestProfile1",
	TenancyOCID:    "ocid1.tenancy.test1",
	UserOCID:       "ocid1.user.test1",
	KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:3a:f9:f9:d1",
	Region:         "eu-frankfurt-99",
	PrivateKey:     &strPrivateTestKey1,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
}

var profile2 model.OciProfile = model.OciProfile{
	Profile:        "TestProfile2",
	TenancyOCID:    "ocid1.tenancy.test2",
	UserOCID:       "ocid1.user.test2",
	KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:32:f2:f2:d2",
	Region:         "eu-frankfurt-22",
	PrivateKey:     &strPrivateTestKey2,
	ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
}

var profile1UpdWithKey model.OciProfile = model.OciProfile{
	Profile:        "TestProfile1",
	TenancyOCID:    "ocid1.tenancy.testUpd",
	UserOCID:       "ocid1.user.testUpd",
	KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:37:f7:f7:d7",
	Region:         "eu-frankfurt-77",
	PrivateKey:     &strPrivateTestKey1Upd,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
}

var profile1UpdWithoutKey model.OciProfile = model.OciProfile{
	Profile:        "TestProfile1",
	TenancyOCID:    "ocid1.tenancy.testUpd",
	UserOCID:       "ocid1.user.testUpd",
	KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:37:f7:f7:d7",
	Region:         "eu-frankfurt-77",
	PrivateKey:     nil,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
}

var profile1UpdWithoutKeyResult model.OciProfile = model.OciProfile{
	Profile:        "TestProfile1",
	TenancyOCID:    "ocid1.tenancy.testUpd",
	UserOCID:       "ocid1.user.testUpd",
	KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:37:f7:f7:d7",
	Region:         "eu-frankfurt-77",
	PrivateKey:     &strPrivateTestKey1,
	ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
}

func (m *MongodbSuite) TestInsertOciProfile_Success() {
	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_configuration").FindOne(context.TODO(), bson.M{
		"_id": profile1.ID,
	})
	require.NoError(m.T(), val.Err())

	var out model.OciProfile
	val.Decode(&out)

	assert.Equal(m.T(), profile1, out)
}

func (m *MongodbSuite) TestGetOciProfilesWithPrivateKey_Success() {
	var profiles []model.OciProfile

	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	err = m.db.AddOciProfile(profile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	profiles, err = m.db.GetOciProfiles(false)
	require.NoError(m.T(), err)

	expected := []model.OciProfile{
		{
			ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			Profile:        "TestProfile1",
			TenancyOCID:    "ocid1.tenancy.test1",
			UserOCID:       "ocid1.user.test1",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:3a:f9:f9:d1",
			Region:         "eu-frankfurt-99",
			PrivateKey:     &strPrivateTestKey1,
		},
		{
			ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			Profile:        "TestProfile2",
			TenancyOCID:    "ocid1.tenancy.test2",
			UserOCID:       "ocid1.user.test2",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:32:f2:f2:d2",
			Region:         "eu-frankfurt-22",
			PrivateKey:     &strPrivateTestKey2,
		},
	}
	assert.EqualValues(m.T(), expected, profiles)
}

func (m *MongodbSuite) TestGetOciProfilesWithoutPrivateKey_Success() {
	var profiles []model.OciProfile

	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	err = m.db.AddOciProfile(profile2)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	profiles, err = m.db.GetOciProfiles(true)
	require.NoError(m.T(), err)

	expected := []model.OciProfile{
		{
			ID:             utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
			Profile:        "TestProfile1",
			TenancyOCID:    "ocid1.tenancy.test1",
			UserOCID:       "ocid1.user.test1",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:3a:f9:f9:d1",
			Region:         "eu-frankfurt-99",
			PrivateKey:     nil,
		},
		{
			ID:             utils.Str2oid("5dd40bfb12f54dfda2b2c292"),
			Profile:        "TestProfile2",
			TenancyOCID:    "ocid1.tenancy.test2",
			UserOCID:       "ocid1.user.test2",
			KeyFingerprint: "04:12:b5:62:75:e9:be:d2:0e:54:1e:de:32:f2:f2:d2",
			Region:         "eu-frankfurt-22",
			PrivateKey:     nil,
		},
	}
	assert.EqualValues(m.T(), expected, profiles)
}

func (m *MongodbSuite) TestUpdateOciProfileWithPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.OciProfile
	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_configuration").FindOne(context.TODO(), bson.M{
		"_id": profile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateOciProfile(profile1UpdWithKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapOciProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[profile1UpdWithKey.ID], profile1UpdWithKey)
}

func (m *MongodbSuite) TestUpdateOciProfileWithoutPrivateKey_Success() {
	var expected map[primitive.ObjectID]model.OciProfile
	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_configuration").FindOne(context.TODO(), bson.M{
		"_id": profile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.UpdateOciProfile(profile1UpdWithoutKey)
	require.NoError(m.T(), err)

	expected, err = m.db.GetMapOciProfiles()
	require.NoError(m.T(), err)

	assert.EqualValues(m.T(), expected[profile1UpdWithKey.ID], profile1UpdWithoutKeyResult)
}

func (m *MongodbSuite) TestdeleteOciProfile_Success() {
	err := m.db.AddOciProfile(profile1)
	require.NoError(m.T(), err)
	defer m.db.Client.Database(m.dbname).Collection("oci_configuration").DeleteMany(context.TODO(), bson.M{})
	val := m.db.Client.Database(m.dbname).Collection("oci_configuration").FindOne(context.TODO(), bson.M{
		"_id": profile1.ID,
	})
	require.NoError(m.T(), val.Err())

	err = m.db.DeleteOciProfile(profile1.ID)
	require.NoError(m.T(), err)

	val = m.db.Client.Database(m.dbname).Collection("oci_configuration").FindOne(context.TODO(), bson.M{
		"_id": profile1.ID,
	})
	require.Error(m.T(), val.Err())
}
