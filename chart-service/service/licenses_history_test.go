package service

import (
	"testing"
	time "time"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/stretchr/testify/assert"
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
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 0,
					Covered:  0,
				},
				{
					Date:     time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed: 0,
					Covered:  0,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed: 0,
					Covered:  0,
				},
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 0,
					Covered:  0,
				},
			},
		},
		// Keep only one for day #1
		{
			input: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 14, 3, 0, 0, 0, time.Local),
					Consumed: 1,
					Covered:  2,
				},
				{
					Date:     time.Date(2021, 6, 14, 12, 0, 0, 0, time.Local),
					Consumed: 3,
					Covered:  4,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 14, 0, 0, 0, 0, time.Local),
					Consumed: 3,
					Covered:  4,
				},
			},
		},
		// Keep only one for day #2
		{
			input: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 14, 3, 0, 0, 0, time.Local),
					Consumed: 1,
					Covered:  2,
				},
				{
					Date:     time.Date(2021, 6, 14, 12, 0, 0, 0, time.Local),
					Consumed: 3,
					Covered:  4,
				},
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 5,
					Covered:  6,
				},
				{
					Date:     time.Date(2021, 6, 15, 1, 0, 0, 0, time.Local),
					Consumed: 7,
					Covered:  8,
				},
				{
					Date:     time.Date(2021, 6, 15, 1, 5, 0, 0, time.Local),
					Consumed: 9,
					Covered:  10,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 14, 00, 0, 0, 0, time.Local),
					Consumed: 3,
					Covered:  4,
				},
				{
					Date:     time.Date(2021, 6, 15, 00, 0, 0, 0, time.Local),
					Consumed: 9,
					Covered:  10,
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
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 21,
					Covered:  22,
				},
			},
			b: []dto.LicenseComplianceHistoricValue{},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 21,
					Covered:  22,
				},
			},
		},
		{
			a: []dto.LicenseComplianceHistoricValue{},
			b: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 42,
					Covered:  43,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 42,
					Covered:  43,
				},
			},
		},
		{
			a: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 10,
					Covered:  10,
				},
				{
					Date:     time.Date(2021, 6, 18, 0, 0, 0, 0, time.Local),
					Consumed: 10,
					Covered:  10,
				},
				{
					Date:     time.Date(2021, 6, 21, 0, 0, 0, 0, time.Local),
					Consumed: 10,
					Covered:  10,
				},
			},
			b: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 20,
					Covered:  20,
				},
				{
					Date:     time.Date(2021, 6, 17, 0, 0, 0, 0, time.Local),
					Consumed: 20,
					Covered:  20,
				},
			},
			expected: []dto.LicenseComplianceHistoricValue{
				{
					Date:     time.Date(2021, 6, 15, 0, 0, 0, 0, time.Local),
					Consumed: 30,
					Covered:  30,
				},
				{
					Date:     time.Date(2021, 6, 17, 0, 0, 0, 0, time.Local),
					Consumed: 20,
					Covered:  20,
				},
				{
					Date:     time.Date(2021, 6, 18, 0, 0, 0, 0, time.Local),
					Consumed: 10,
					Covered:  10,
				},
				{
					Date:     time.Date(2021, 6, 21, 0, 0, 0, 0, time.Local),
					Consumed: 10,
					Covered:  10,
				},
			},
		},
	}

	for _, testCase := range testCases {
		actual := mergeLicenseComplianceHistoricValues(testCase.a, testCase.b)
		assert.Equal(t, testCase.expected, actual)
	}
}
