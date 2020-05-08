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

package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDiffFeature(t *testing.T) {
	differences := DiffFeature([]Feature{
		{Name: "Spatial queryes", Status: true},
		{Name: "High heavy database", Status: false},
		{Name: "Power saving", Status: false},
		{Name: "Crackked", Status: true},
		{Name: "Star wars support", Status: true},
		{Name: "Star wars support SP3", Status: false},
	}, []Feature{
		{Name: "Spatial queryes", Status: true},
		{Name: "High heavy database", Status: true},
		{Name: "Power saving", Status: false},
		{Name: "Crackked", Status: false},
		{Name: "Dark Engine", Status: true},
	})

	assert.Equal(t, differences["Spatial queryes"], DiffFeatureActive)
	assert.Equal(t, differences["High heavy database"], DiffFeatureActivated)
	assert.Equal(t, differences["Power saving"], DiffFeatureInactive)
	assert.Equal(t, differences["Crackked"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Dark Engine"], DiffFeatureActivated)
	assert.Equal(t, differences["Star wars support"], DiffFeatureDeactivated)
	assert.Equal(t, differences["Star wars support SP3"], DiffFeatureInactive)
}
