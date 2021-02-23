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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestDeleteAllNoDataAlerts_Success() {
	var alert1 model.Alert = model.Alert{
		AlertCode:               model.AlertCodeNewServer,
		AlertSeverity:           model.AlertSeverityInfo,
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertStatus:             model.AlertStatusNew,
		Date:                    utils.P("2019-11-05T18:02:03Z"),
		Description:             "pippo",
		OtherInfo: map[string]interface{}{
			"hostname": "myhost",
		},
		ID: utils.Str2oid("5dd40bfb12f54dfda7b1c291"),
	}

	var alert3 model.Alert = model.Alert{
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertStatus:             model.AlertStatusNew,
		Date:                    utils.P("2019-11-05T18:02:03Z"),
		Description:             "test desc pippo",
		OtherInfo: map[string]interface{}{
			"hostname": "pippo-host",
		},
		ID: utils.Str2oid("5dd40bfb12f54dfda7b1c292"),
	}

	var alert4 model.Alert = model.Alert{
		AlertCode:               model.AlertCodeNoData,
		AlertSeverity:           model.AlertSeverityCritical,
		AlertAffectedTechnology: nil,
		AlertCategory:           model.AlertCategoryEngine,
		AlertStatus:             model.AlertStatusNew,
		Date:                    utils.P("2019-12-25T18:02:03Z"),
		Description:             "test desc pluto",
		OtherInfo: map[string]interface{}{
			"hostname": "pluto-host",
		},
		ID: utils.Str2oid("5dd40bfb12f54dfda7b1c293"),
	}

	defer m.db.Client.Database(m.dbname).Collection("alerts").DeleteMany(context.TODO(), bson.M{})

	_, err := m.db.Client.Database(m.dbname).Collection("alerts").InsertOne(context.TODO(), alert1)
	require.NoError(m.T(), err)

	_, err = m.db.Client.Database(m.dbname).Collection("alerts").InsertOne(context.TODO(), alert3)
	require.NoError(m.T(), err)
	_, err = m.db.Client.Database(m.dbname).Collection("alerts").InsertOne(context.TODO(), alert4)
	require.NoError(m.T(), err)

	err = m.db.DeleteAllNoDataAlerts()
	require.NoError(m.T(), err)

	// Check that there are no more AlertCodeNoData alerts
	val, erro := m.db.Client.Database(m.dbname).Collection("alerts").
		Find(context.TODO(), bson.M{"alertCode": model.AlertCodeNoData})
	require.NoError(m.T(), erro)

	res := make([]model.Alert, 0)
	erro = val.All(context.TODO(), &res)
	require.NoError(m.T(), erro)
	require.Equal(m.T(), 0, len(res))

	// Check that there's still alert1
	val, erro = m.db.Client.Database(m.dbname).Collection("alerts").
		Find(context.TODO(), bson.M{})
	require.NoError(m.T(), erro)

	res = make([]model.Alert, 0)
	erro = val.All(context.TODO(), &res)
	require.NoError(m.T(), erro)
	require.Equal(m.T(), 1, len(res))
	require.Equal(m.T(), alert1, (res)[0])
}
