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
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package service

import (
	"testing"
	time "time"

	"github.com/stretchr/testify/assert"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
)

func TestGetLicenseComplianceHistory(t *testing.T) {
	// TODO
}

func TestSortAndKeepOnlyLastEntryOfEachDay(t *testing.T) {
	testCases := []struct {
		input    []dto.LicenseComplianceHistoricValue
		expected []dto.LicenseComplianceHistoricValue
	}{
		{
			input:    []dto.LicenseComplianceHistoricValue{},
			expected: []dto.LicenseComplianceHistoricValue{},
		},
		// Sort
		{
			input: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  0,
					Covered:   0,
					Purchased: 0,
				},
				{
					Date:      time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed:  0,
					Covered:   0,
					Purchased: 0,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed:  0,
					Covered:   0,
					Purchased: 0,
				},
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  0,
					Covered:   0,
					Purchased: 0,
				},
			},
		},
		// Keep only one for day #1
		{
			input: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 14, 3, 0, 0, 0, time.Local),
					Consumed:  1,
					Covered:   2,
					Purchased: 5,
				},
				{
					Date:      time.Date(2021, 6, 14, 12, 0, 0, 0, time.Local),
					Consumed:  3,
					Covered:   4,
					Purchased: 5,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed:  3,
					Covered:   4,
					Purchased: 5,
				},
			},
		},
		// Keep only one for day #2
		{
			input: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 14, 3, 0, 0, 0, time.Local),
					Consumed:  1,
					Covered:   2,
					Purchased: 5,
				},
				{
					Date:      time.Date(2021, 6, 14, 12, 0, 0, 0, time.Local),
					Consumed:  3,
					Covered:   4,
					Purchased: 5,
				},
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  5,
					Covered:   6,
					Purchased: 10,
				},
				{
					Date:      time.Date(2021, 6, 15, 1, 0, 0, 0, time.Local),
					Consumed:  7,
					Covered:   8,
					Purchased: 10,
				},
				{
					Date:      time.Date(2021, 6, 15, 1, 5, 0, 0, time.Local),
					Consumed:  9,
					Covered:   10,
					Purchased: 11,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 14, 00, 0, 0, 0, time.Local),
					Consumed:  3,
					Covered:   4,
					Purchased: 5,
				},
				{
					Date:      time.Date(2021, 6, 15, 00, 0, 0, 0, time.Local),
					Consumed:  9,
					Covered:   10,
					Purchased: 11,
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual := sortAndKeepOnlyLastEntryOfEachDay(testCase.input)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestMergeMySqlLicensesCompliance(t *testing.T) {
	testCases := []struct {
		input    []dto.LicenseComplianceHistory
		expected []dto.LicenseComplianceHistory
	}{
		{
			input:    []dto.LicenseComplianceHistory{},
			expected: []dto.LicenseComplianceHistory{},
		},
		{
			input: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
			expected: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
		},
		{
			input: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "MySQL Enterprise per cluster",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
			expected: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "MySQL Enterprise",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
		},
		{
			input: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "MySQL Enterprise per cluster",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "MySQL Enterprise per host",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
			expected: []dto.LicenseComplianceHistory{
				{
					LicenseTypeID:   "A00001",
					ItemDescription: "Something",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "pippo",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
				{
					LicenseTypeID:   "",
					ItemDescription: "MySQL Enterprise",
					Metric:          "",
					History:         []dto.LicenseComplianceHistoricValue{},
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual := mergeMySqlLicensesCompliance(testCase.input)
		assert.Equal(t, testCase.expected, actual)
	}
}

func TestMergeLicenseComplianceHistoricValues(t *testing.T) {
	testCases := []struct {
		a        []dto.LicenseComplianceHistoricValue
		b        []dto.LicenseComplianceHistoricValue
		expected []dto.LicenseComplianceHistoricValue
	}{
		{
			a:        []dto.LicenseComplianceHistoricValue{},
			b:        []dto.LicenseComplianceHistoricValue{},
			expected: []dto.LicenseComplianceHistoricValue{},
		},

		{
			a: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  21,
					Covered:   22,
					Purchased: 25,
				},
			},
			b: []dto.LicenseComplianceHistoricValue{},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  21,
					Covered:   22,
					Purchased: 25,
				},
			},
		},
		{
			a: []dto.LicenseComplianceHistoricValue{},
			b: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  42,
					Covered:   43,
					Purchased: 25,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  42,
					Covered:   43,
					Purchased: 25,
				},
			},
		},
		{
			a: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  10,
					Covered:   10,
					Purchased: 10,
				},
				{
					Date:      time.Date(2021, 6, 18, 0, 0, 0, 0, time.Local),
					Consumed:  10,
					Covered:   10,
					Purchased: 10,
				},
				{
					Date:      time.Date(2021, 6, 21, 0, 0, 0, 0, time.Local),
					Consumed:  10,
					Covered:   10,
					Purchased: 10,
				},
			},
			b: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  20,
					Covered:   20,
					Purchased: 20,
				},
				{
					Date:      time.Date(2021, 6, 17, 0, 0, 0, 0, time.Local),
					Consumed:  20,
					Covered:   20,
					Purchased: 20,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:      time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed:  30,
					Covered:   30,
					Purchased: 30,
				},
				{
					Date:      time.Date(2021, 6, 17, 0, 0, 0, 0, time.Local),
					Consumed:  20,
					Covered:   20,
					Purchased: 20,
				},
				{
					Date:      time.Date(2021, 6, 18, 0, 0, 0, 0, time.Local),
					Consumed:  10,
					Covered:   10,
					Purchased: 10,
				},
				{
					Date:      time.Date(2021, 6, 21, 0, 0, 0, 0, time.Local),
					Consumed:  10,
					Covered:   10,
					Purchased: 10,
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual := mergeLicenseComplianceHistoricValues(testCase.a, testCase.b)
		assert.Equal(t, testCase.expected, actual)
	}
}
