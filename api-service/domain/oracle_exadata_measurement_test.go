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
package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	oem = OracleExadataMeasurement{
		unparsedValue: "3.1415926535T",
		Symbol:        "TIB",
		Quantity:      3.1415926535,
	}
)

func TestString(t *testing.T) {
	expected := "3.141593 TIB"
	actual := oem.String()

	assert.Equal(t, expected, actual)

}

func TestHuman(t *testing.T) {
	expected := "3216.99 GIB"
	actual, err := oem.Human("GIB")

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestConvert(t *testing.T) {
	expected := &OracleExadataMeasurement{
		unparsedValue: "3216.991232 GIB",
		Symbol:        "GIB",
		Quantity:      3216.991232,
	}

	actual, err := Convert(oem, "GIB")

	require.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestAdd(t *testing.T) {
	original := &oem
	expected := &OracleExadataMeasurement{
		unparsedValue: "3.143546 TIB",
		Symbol:        "TIB",
		Quantity:      3.143546,
	}

	adding := struct {
		qty    float64
		symbol string
	}{2, "GIB"}

	original.Add(adding.qty, adding.symbol)

	assert.Equal(t, expected, original)
}

func TestSub(t *testing.T) {
	original := &OracleExadataMeasurement{
		unparsedValue: "10 TIB",
		Symbol:        "TIB",
		Quantity:      10,
	}

	expected := &OracleExadataMeasurement{
		unparsedValue: "9.998047 TIB",
		Symbol:        "TIB",
		Quantity:      9.998047,
	}

	subtracting := OracleExadataMeasurement{
		unparsedValue: "2 GIB",
		Symbol:        "GIB",
		Quantity:      2,
	}

	original.Sub(subtracting)

	assert.Equal(t, expected, original)
}

func TestNewOracleExadataMeasurement(t *testing.T) {
	expected := &OracleExadataMeasurement{
		Symbol: "MIB",
	}

	actual := NewOracleExadataMeasurement()

	assert.Equal(t, expected, actual)
}

func TestNewUnknownOracleExadataMeasurement(t *testing.T) {
	expected := &OracleExadataMeasurement{
		unparsedValue: "UNKNOWN",
	}

	actual := NewUnknownOracleExadataMeasurement()

	assert.Equal(t, expected, actual)
}

func TestIntToOracleExadataMeasurement(t *testing.T) {
	expected := &OracleExadataMeasurement{
		unparsedValue: "2 GIB",
		Symbol:        "GIB",
		Quantity:      2,
	}

	actual, err := IntToOracleExadataMeasurement(2, "GIB")

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestStringToOracleExadataMeasurement(t *testing.T) {
	expected := &OracleExadataMeasurement{
		unparsedValue: "2 GIB",
		Symbol:        "GIB",
		Quantity:      2,
	}

	actual, err := StringToOracleExadataMeasurement("2 GIB")

	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestGetPercentage(t *testing.T) {
	total := OracleExadataMeasurement{
		unparsedValue: "10T",
		Symbol:        "TIB",
		Quantity:      10.0,
	}

	expected := "31.44%"

	actual := GetPercentage(oem, total)

	assert.Equal(t, expected, actual)
}
