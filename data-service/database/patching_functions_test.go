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

package database

import (
	"context"
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestFindPatchingFunction() {
	defer m.db.Client.Database(m.dbname).Collection("patching_functions").DeleteMany(context.TODO(), bson.M{})
	id := utils.Str2oid("5ef9e436b48eac8a91f81dc5")
	pf1 := model.PatchingFunction{
		ID:        &id,
		Code:      "dfssdfsdf",
		CreatedAt: utils.P("2020-05-20T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars:      map[string]interface{}{"bar": int32(10)},
	}
	m.InsertPatchingFunction(pf1)

	m.T().Run("should_find_pf", func(t *testing.T) {
		pf, err := m.db.FindPatchingFunction("foobar")
		require.NoError(m.T(), err)
		assert.Equal(m.T(), pf1, pf)
	})

	m.T().Run("should_not_find_pf", func(t *testing.T) {
		pf, err := m.db.FindPatchingFunction("foobar2")
		require.NoError(m.T(), err)
		assert.Equal(m.T(), model.PatchingFunction{}, pf)
	})
}
