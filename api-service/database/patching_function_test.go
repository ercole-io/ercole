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

package database

import (
	"context"

	"github.com/amreo/ercole-services/model"
	"github.com/amreo/ercole-services/utils"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestSaveAndFindPatchingFunction() {
	defer m.db.Client.Database(m.dbname).Collection("patching_functions").DeleteMany(context.TODO(), bson.M{})

	pf, err := m.db.FindPatchingFunction("foobar")

	m.Require().NoError(err)
	m.Assert().Equal(model.PatchingFunction{}, pf)

	err = m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        nil,
		Code:      "dfssdfsdf",
		CreatedAt: utils.P("2020-05-20T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars:      map[string]interface{}{"bar": 10},
	})

	m.Require().NoError(err)

	pf, err = m.db.FindPatchingFunction("foobar")

	m.Assert().Equal("dfssdfsdf", pf.Code)
	m.Assert().Equal(utils.P("2020-05-20T09:53:34+00:00").UTC(), pf.CreatedAt)
	m.Assert().Equal("foobar", pf.Hostname)
	m.Assert().Equal(map[string]interface{}{"bar": int32(10)}, pf.Vars)

	err = m.db.SavePatchingFunction(model.PatchingFunction{
		ID:        pf.ID,
		Code:      "ffff",
		CreatedAt: utils.P("2020-05-21T09:53:34+00:00").UTC(),
		Hostname:  "foobar",
		Vars:      map[string]interface{}{"bar": 2},
	})

	m.Require().NoError(err)

	pf, err = m.db.FindPatchingFunction("foobar")

	m.Assert().Equal("ffff", pf.Code)
	m.Assert().Equal(utils.P("2020-05-21T09:53:34+00:00").UTC(), pf.CreatedAt)
	m.Assert().Equal("foobar", pf.Hostname)
	m.Assert().Equal(map[string]interface{}{"bar": int32(2)}, pf.Vars)

}
