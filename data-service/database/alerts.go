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
	"go.mongodb.org/mongo-driver/bson"
)

// DeleteAllNoDataAlerts delete all alerts with code NO_DATA
func (md *MongoDatabase) DeleteAllNoDataAlerts() utils.AdvancedErrorInterface {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("alerts").
		DeleteMany(context.TODO(), bson.M{"alertCode": model.AlertCodeNoData})

	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return nil
}

// DeleteNoDataAlertByHost delete NO_DATA alert by hostname
func (md *MongoDatabase) DeleteNoDataAlertByHost(hostname string) utils.AdvancedErrorInterface {
	_, err := md.Client.Database(md.Config.Mongodb.DBName).
		Collection("alerts").
		DeleteOne(context.TODO(),
			bson.M{
				"alertCode":          model.AlertCodeNoData,
				"otherInfo.hostname": hostname,
			})

	if err != nil {
		return utils.NewAdvancedErrorPtr(err, "DB ERROR")
	}

	return nil
}
