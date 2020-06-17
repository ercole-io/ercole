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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffLicenses(t *testing.T) {
	differences := DiffLicenses([]OracleDatabaseLicense{
		{Name: "Spatial queryes", Count: 10},
		{Name: "High heavy database", Count: 0},
		{Name: "Power saving", Count: 0},
		{Name: "Crackked", Count: 12},
		{Name: "Star wars support", Count: 6},
		{Name: "Star wars support SP3", Count: 0},
	}, []OracleDatabaseLicense{
		{Name: "Spatial queryes", Count: 5},
		{Name: "High heavy database", Count: 7},
		{Name: "Power saving", Count: 0},
		{Name: "Crackked", Count: 0},
		{Name: "Dark Engine", Count: 13},
	})

	assert.Equal(t, differences["Spatial queryes"], DiffFeatureActive)
	assert.Equal(t, differences["High heavy database"], DiffFeatureActivated)
	assert.Equal(t, differences["Power saving"], DiffFeatureInactive)
	assert.Equal(t, differences["Crackked"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Dark Engine"], DiffFeatureActivated)
	assert.Equal(t, differences["Star wars support"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Star wars support SP3"], DiffFeatureInactive)
}
