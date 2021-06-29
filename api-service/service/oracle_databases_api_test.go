// Copyright (c) 2020 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSearchOracleDatabaseAddms_Success(t *testing.T) {
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

	db.EXPECT().SearchOracleDatabaseAddms(
		[]string{"foo", "bar", "foobarx"}, "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchOracleDatabaseAddms(
		"foo bar foobarx", "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchOracleDatabaseAddms_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchOracleDatabaseAddms(
		[]string{"foo", "bar", "foobarx"}, "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchOracleDatabaseAddms(
		"foo bar foobarx", "Benefit",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchOracleDatabaseSegmentAdvisorsXLSX_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
		Database: db,
	}
	data := []dto.OracleDatabaseSegmentAdvisor{
		{
			Reclaimable:    4.3,
			SegmentsSize:   50,
			Hostname:       "test-db3",
			Dbname:         "foobar3",
			SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			SegmentType:    "TABLE",
			PartitionName:  "iyyiuyyoy",
			Recommendation: "32b36a77e7481343ef175483c086859e",
		},
		{
			Reclaimable:    4.3,
			SegmentsSize:   0,
			Hostname:       "test-db3",
			Dbname:         "foobar3",
			SegmentOwner:   "Brittany-424f6a749eef846fa40a1ad1ee3d3674",
			SegmentName:    "pasta-973e4d1f937da4d9bc1b092f934ab0ec",
			SegmentType:    "TABLE",
			PartitionName:  "iyyiuyyoy",
			Recommendation: "32b36a77e7481343ef175483c086859e",
		},
	}

	db.EXPECT().SearchOracleDatabaseSegmentAdvisors(
		[]string{}, "",
		false, "Italy", "TST", utils.P("2019-12-05T14:02:03Z"),
	).Return(data, nil).Times(1)

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2019-12-05T14:02:03Z"),
	}

	actual, err := as.SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter)
	require.NoError(t, err)

	assert.Equal(t, "4.3", actual.GetCellValue("Segment_Advisor", "A2"))
	assert.Equal(t, "50", actual.GetCellValue("Segment_Advisor", "B2"))
	assert.Equal(t, "0.086", actual.GetCellValue("Segment_Advisor", "C2"))
	assert.Equal(t, "test-db3", actual.GetCellValue("Segment_Advisor", "D2"))
	assert.Equal(t, "foobar3", actual.GetCellValue("Segment_Advisor", "E2"))
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", actual.GetCellValue("Segment_Advisor", "F2"))
	assert.Equal(t, "pasta-973e4d1f937da4d9bc1b092f934ab0ec", actual.GetCellValue("Segment_Advisor", "G2"))
	assert.Equal(t, "TABLE", actual.GetCellValue("Segment_Advisor", "H2"))
	assert.Equal(t, "iyyiuyyoy", actual.GetCellValue("Segment_Advisor", "I2"))

	assert.Equal(t, "4.3", actual.GetCellValue("Segment_Advisor", "A3"))
	assert.Equal(t, "0", actual.GetCellValue("Segment_Advisor", "B3"))
	assert.Equal(t, "", actual.GetCellValue("Segment_Advisor", "C3"))
	assert.Equal(t, "test-db3", actual.GetCellValue("Segment_Advisor", "D3"))
	assert.Equal(t, "foobar3", actual.GetCellValue("Segment_Advisor", "E3"))
	assert.Equal(t, "Brittany-424f6a749eef846fa40a1ad1ee3d3674", actual.GetCellValue("Segment_Advisor", "F3"))
	assert.Equal(t, "pasta-973e4d1f937da4d9bc1b092f934ab0ec", actual.GetCellValue("Segment_Advisor", "G3"))
	assert.Equal(t, "TABLE", actual.GetCellValue("Segment_Advisor", "H3"))
	assert.Equal(t, "iyyiuyyoy", actual.GetCellValue("Segment_Advisor", "I3"))

}

func TestSearchOracleDatabaseSegmentAdvisorsXLSX_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchOracleDatabaseSegmentAdvisors(
		[]string{}, "",
		false, "Italy", "TST", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	filter := dto.GlobalFilter{
		Location:    "Italy",
		Environment: "TST",
		OlderThan:   utils.P("2019-12-05T14:02:03Z"),
	}

	res, err := as.SearchOracleDatabaseSegmentAdvisorsAsXLSX(filter)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchOracleDatabasePatchAdvisors_Success(t *testing.T) {
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

	db.EXPECT().SearchOracleDatabasePatchAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD",
		utils.P("2019-12-05T14:02:03Z"), "OK",
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchOracleDatabasePatchAdvisors(
		"foo bar foobarx", "Date",
		true, 1, 1, utils.P("2019-06-05T14:02:03Z"),
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"), "OK",
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchOracleDatabasePatchAdvisors_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchOracleDatabasePatchAdvisors(
		[]string{"foo", "bar", "foobarx"}, "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		"OK",
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchOracleDatabasePatchAdvisors(
		"foo bar foobarx", "Date",
		true, 1, 1,
		utils.P("2019-06-05T14:02:03Z"), "Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
		"OK",
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

func TestSearchOracleDatabases_Success(t *testing.T) {
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

	db.EXPECT().SearchOracleDatabases(
		false, []string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(expectedRes, nil).Times(1)

	res, err := as.SearchOracleDatabases(
		dto.SearchOracleDatabasesFilter{
			dto.GlobalFilter{
				"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
			},
			false, "foo bar foobarx", "Memory",
			true, 1, 1,
		},
	)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}

func TestSearchOracleDatabases_Fail(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
	}

	db.EXPECT().SearchOracleDatabases(
		false, []string{"foo", "bar", "foobarx"}, "Memory",
		true, 1, 1,
		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
	).Return(nil, aerrMock).Times(1)

	res, err := as.SearchOracleDatabases(

		dto.SearchOracleDatabasesFilter{
			dto.GlobalFilter{
				"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
			},
			false, "foo bar foobarx", "Memory",
			true, 1, 1,
		},
	)

	require.Nil(t, res)
	assert.Equal(t, aerrMock, err)
}

//TODO
//func TestSearchLicenses_Success(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	expectedRes := []interface{}{
//		map[string]interface{}{
//			"Compliance": false,
//			"Count":      0,
//			"Used":       5,
//			"_id":        "Oracle ENT",
//		},
//		map[string]interface{}{
//			"Compliance": true,
//			"Count":      0,
//			"Used":       0,
//			"_id":        "Oracle STD",
//		},
//	}

//	db.EXPECT().SearchLicenses(
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	).Return(expectedRes, nil).Times(1)

//	res, err := as.SearchLicenses(
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	)

//	require.NoError(t, err)
//	assert.Equal(t, expectedRes, res)
//}

//func TestSearchLicenses_Fail(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	db.EXPECT().SearchLicenses(
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	).Return(nil, aerrMock).Times(1)

//	res, err := as.SearchLicenses(
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	)

//	require.Nil(t, res)
//	assert.Equal(t, aerrMock, err)
//}

//TODO
//func TestSearchOracleDatabaseUsedLicenses_Success(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	expectedRes := []interface{}{
//		map[string]interface{}{
//			"Compliance": false,
//			"Count":      0,
//			"Used":       5,
//			"_id":        "Oracle ENT",
//		},
//		map[string]interface{}{
//			"Compliance": true,
//			"Count":      0,
//			"Used":       0,
//			"_id":        "Oracle STD",
//		},
//	}

//	db.EXPECT().SearchOracleDatabaseUsedLicenses(
//		"Used",
//		true, 1, 1,
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	).Return(expectedRes, nil).Times(1)

//	res, err := as.SearchOracleDatabaseUsedLicenses(
//		"Used",
//		true, 1, 1,
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	)

//	require.NoError(t, err)
//	assert.Equal(t, expectedRes, res)
//}

//TODO
//func TestSearchOracleDatabaseUsedLicenses_Fail(t *testing.T) {
//	mockCtrl := gomock.NewController(t)
//	defer mockCtrl.Finish()
//	db := NewMockMongoDatabaseInterface(mockCtrl)
//	as := APIService{
//		Database: db,
//	}

//	db.EXPECT().SearchOracleDatabaseUsedLicenses(
//		"Used",
//		true, 1, 1,
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	).Return(nil, aerrMock).Times(1)

//	res, err := as.SearchOracleDatabaseUsedLicenses(
//		"Used",
//		true, 1, 1,
//		"Italy", "PROD", utils.P("2019-12-05T14:02:03Z"),
//	)

//	require.Nil(t, res)
//	assert.Equal(t, aerrMock, err)
//}
