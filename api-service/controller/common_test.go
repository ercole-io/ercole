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

package controller

import (
	"errors"
	"testing"

	"github.com/amreo/ercole-services/utils"
	"github.com/plandem/xlsx"
	"github.com/stretchr/testify/assert"
)

//go:generate mockgen -source ../service/service.go -destination=fake_service.go -package=controller

//Common data
var errMock error = errors.New("MockError")
var aerrMock utils.AdvancedErrorInterface = utils.NewAdvancedErrorPtr(errMock, "mock")

func AssertXLSXFloat(t *testing.T, expected float64, cell *xlsx.Cell) {
	actual, err := cell.Float()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func AssertXLSXInt(t *testing.T, expected int, cell *xlsx.Cell) {
	actual, err := cell.Int()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func AssertXLSXBool(t *testing.T, expected bool, cell *xlsx.Cell) {
	actual, err := cell.Bool()
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

type FailingReader struct{}

func (fr *FailingReader) Read(p []byte) (n int, err error) {
	return 0, aerrMock
}
