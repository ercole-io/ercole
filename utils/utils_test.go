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

package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHideMongoDBPassword_Success(t *testing.T) {
	x := HideMongoDBPassword("mongodb://localhost:27017/")
	assert.Equal(t, x, "mongodb://localhost:27017/")
}

func TestHideMongoDBPassword_Success2(t *testing.T) {
	x := HideMongoDBPassword("mongodb://sysop:moon@localhost:27017/")
	assert.Equal(t, x, "mongodb://***:***@localhost:27017/")
}

func TestHideMongoDBPassword_Fail(t *testing.T) {
	x := HideMongoDBPassword("mongodb://sysop:moon@localhost:27017/")
	assert.NotEqual(t, x, "mongodb://sysop:moon@localhost:27017/")
}
