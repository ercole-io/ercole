// Copyright (c) 2024 Sorint.lab S.p.A.
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
package dto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLicensesCount_standard(t *testing.T) {
	edition := "standard"
	vcpus := 9

	actual := getLicensesCount(edition, vcpus)

	assert.Equal(t, 2, actual)
}

func TestGetLicensesCount_enterprise(t *testing.T) {
	edition := "enterprise"
	vcpus := 9

	actual := getLicensesCount(edition, vcpus)

	assert.Equal(t, 4, actual)
}

func TestGetLicensesCount_zero(t *testing.T) {
	edition := "enterprise"
	vcpus := 0

	actual := getLicensesCount(edition, vcpus)

	assert.Equal(t, 1, actual)
}
