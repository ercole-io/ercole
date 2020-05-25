// Copyright (c) 2019 Sorint.lab S.p.A.
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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"testing"

	"github.com/ercole-io/ercole/utils"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSearchAddms_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"Action":         "Run SQL Tuning Advisor on the SELECT statement with SQL_ID \"4ztz048yfq32s\".",
			"Benefit":        83.34,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "SQL statements consuming significant database time were found. These statements offer a good opportunity for performance improvement.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "SQL Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
		{
			"Action":         "Look at the \"Top SQL Statements\" finding for SQL statements consuming significant I/O on this segment. For example, the SELECT statement with SQL_ID \"4ztz048yfq32s\" is responsible for 100% of \"User I/O\" and \"Cluster\" waits for this segment.",
			"Benefit":        68.24,
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Finding":        "Individual database segments responsible for significant \"User I/O\" and \"Cluster\" waits were found.",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"Recommendation": "Segment Tuning",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	db.EXPECT().SearchAddms(
		[]string{"foo", "bar", "foobarx"}, "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchAddms(
		"foo bar foobarx", "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchAddms_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchAddms(
		[]string{"foo", "bar", "foobarx"}, "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchAddms(
		"foo bar foobarx", "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchSegmentAdvisors_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.82+02:00"),
			"Dbname":         "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Environment":    "SVIL",
			"Hostname":       "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":       "Germany",
			"PartitionName":  "",
			"Reclaimable":    "\u003c1",
			"Recommendation": "3d7e603f515ed171fc99bdb908f38fb2",
			"SegmentName":    "nascar1-f9b3703bf8b3cc7ae070cd28e7fed7b3",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":      utils.P("2020-04-07T08:52:59.872+02:00"),
			"Dbname":         "ERCOLE",
			"Environment":    "TST",
			"Hostname":       "test-db",
			"Location":       "Germany",
			"PartitionName":  "iyyiuyyoy",
			"Reclaimable":    "\u003c1",
			"Recommendation": "32b36a77e7481343ef175483c086859e",
			"SegmentName":    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			"SegmentOwner":   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			"SegmentType":    "TABLE",
			"_id":            utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	db.EXPECT().SearchSegmentAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Reclaimable",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchSegmentAdvisors(
		"foo bar foobarx", "Reclaimable",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchSegmentAdvisors_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchSegmentAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Reclaimable",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchSegmentAdvisors(
		"foo bar foobarx", "Reclaimable",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchPatchAdvisors_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.82+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Dbver":       "11.2.0.3.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "SVIL",
			"Hostname":    "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"CreatedAt":   utils.P("2020-04-07T08:52:59.872+02:00"),
			"Date":        utils.P("2012-04-16T02:00:00+02:00"),
			"Dbname":      "ERCOLE",
			"Dbver":       "12.2.0.1.0 Enterprise Edition",
			"Description": "PSU 11.2.0.3.2",
			"Environment": "TST",
			"Hostname":    "test-db",
			"Location":    "Germany",
			"Status":      "KO",
			"_id":         utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	db.EXPECT().SearchPatchAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD",
		utils.P("2019-12-05T14:02:03Z"), "OK",
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchPatchAdvisors(
		"foo bar foobarx", "Date",
		true, 1, 1, utils.P("2019-06-05T14:02:03Z"),
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"), "OK",
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchPatchAdvisors_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchPatchAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		"OK",
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchPatchAdvisors(
		"foo bar foobarx", "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		"OK",
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchDatabases_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []map[string]interface{}{
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "16",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.82+02:00"),
			"DatafileSize":     "61",
			"Dataguard":        false,
			"Environment":      "SVIL",
			"HA":               false,
			"Hostname":         "publicitate-36d06ca83eafa454423d2097f4965517",
			"Location":         "Germany",
			"Memory":           4.199,
			"Name":             "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"RAC":              false,
			"SegmentsSize":     "41",
			"Status":           "OPEN",
			"UniqueName":       "4wcqjn-ecf040bdfab7695ab332aef7401f185c",
			"Version":          "11.2.0.3.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd32"),
		},
		{
			"ArchiveLogStatus": false,
			"BlockSize":        "8192",
			"CPUCount":         "2",
			"Charset":          "AL32UTF8",
			"CreatedAt":        utils.P("2020-04-07T08:52:59.872+02:00"),
			"DatafileSize":     "6",
			"Dataguard":        false,
			"Environment":      "TST",
			"HA":               false,
			"Hostname":         "test-db",
			"Location":         "Germany",
			"Memory":           1.484,
			"Name":             "ERCOLE",
			"RAC":              false,
			"SegmentsSize":     "3",
			"Status":           "OPEN",
			"UniqueName":       "ERCOLE",
			"Version":          "12.2.0.1.0 Enterprise Edition",
			"Work":             "1",
			"_id":              utils.Str2oid("5e8c234b24f648a08585bd43"),
		},
	}

	db.EXPECT().SearchDatabases(
		false, []string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchDatabases(
		false, "foo bar foobarx", "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchDatabases_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchDatabases(
		false, []string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchDatabases(
		false, "foo bar foobarx", "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestListLicenses_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := []interface{}{
		map[string]interface{}{
			"Compliance": false,
			"Count":      0,
			"Used":       5,
			"_id":        "Oracle ENT",
		},
		map[string]interface{}{
			"Compliance": true,
			"Count":      0,
			"Used":       0,
			"_id":        "Oracle STD",
		},
	}

	db.EXPECT().ListLicenses(
		false, "Used",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.ListLicenses(
		false, "Used",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestListLicenses_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().ListLicenses(
		false, "Used",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.ListLicenses(
		false, "Used",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestGetLicense_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	expectedRes := map[string]interface{}{
		"Compliance": false,
		"Count":      0,
		"Hosts": []interface{}{
			map[string]interface{}{
				"Databases": []interface{}{
					"ERCOLE",
					"urcole",
				},
				"Hostname": "itl-csllab-112.sorint.localpippo",
			},
			map[string]interface{}{
				"Databases": []interface{}{
					"ERCOLE",
				},
				"Hostname": "test-db",
			},
			map[string]interface{}{
				"Databases": []interface{}{
					"rudeboy-fb3160a04ffea22b55555bbb58137f77",
					"007bond-f260462ca34bbd17deeda88f042e42a1",
					"jacket-d4a157354d91bfc68fce6f45546d8f3d",
					"allstate-9a6a2a820a3f61aeb345a834abf40fba",
					"4wcqjn-ecf040bdfab7695ab332aef7401f185c",
				},
				"Hostname": "publicitate-36d06ca83eafa454423d2097f4965517",
			},
		},
		"Used": 5,
		"_id":  "Oracle ENT",
	}

	db.EXPECT().GetLicense(
		"Oracle ENT", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.GetLicense(
		"Oracle ENT", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestGetLicense_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().GetLicense(
		"Oracle ENT", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.GetLicense(
		"Oracle ENT", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSetLicenseCount_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SetLicenseCount(
		"Oracle ENT", 1,
	).Return(nil).Times(1)

	err := as.SetLicenseCount("Oracle ENT", 1)

	require.NoError(t, err)
}

func TestSetLicenseCount_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SetLicenseCount(
		"Oracle ENT", 1,
	).Return(aerrMock).Times(1)

	err := as.SetLicenseCount("Oracle ENT", 1)

	require.Equal(t, aerrMock, err)
}

func TestSetLicensesCount_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SetLicenseCount("Oracle ENT", 10).Return(nil).Times(1)
	db.EXPECT().SetLicenseCount("Oracle STD", 20).Return(nil).Times(1)

	err := as.SetLicensesCount([]map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": 10,
		},
		{
			"_id":   "Oracle STD",
			"Count": 20,
		},
	})

	require.NoError(t, err)
}

func TestSetLicensesCount_SuccessEmpty(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	err := as.SetLicensesCount([]map[string]interface{}{})

	require.NoError(t, err)
}

func TestSetLicensesCount_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SetLicenseCount("Oracle ENT", 10).Return(aerrMock).Times(1)

	err := as.SetLicensesCount([]map[string]interface{}{
		{
			"_id":   "Oracle ENT",
			"Count": 10,
		},
		{
			"_id":   "Oracle STD",
			"Count": 20,
		},
	})

	require.Equal(t, aerrMock, err)
}
