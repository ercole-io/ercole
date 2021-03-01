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

func TestDatabasesArrayAsMap(t *testing.T) {
	dbs := []OracleDatabase{
		{
			Name:     "superdb",
			Platform: "HPP 1.2.3",
		},
		{
			Name:     "superdb2",
			Platform: "PPH 1.2.3",
		},
	}

	dbsMap := DatabasesArrayAsMap(dbs)
	assert.Len(t, dbsMap, 2)
	assert.Equal(t, OracleDatabase{
		Name:     "superdb",
		Platform: "HPP 1.2.3",
	}, dbsMap["superdb"])
	assert.Equal(t, OracleDatabase{
		Name:     "superdb2",
		Platform: "PPH 1.2.3",
	}, dbsMap["superdb2"])
}
