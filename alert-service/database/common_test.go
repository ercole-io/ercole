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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

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

var alert2 model.Alert = model.Alert{
	AlertAffectedTechnology: nil,
	AlertCategory:           model.AlertCategoryAgent,
	AlertCode:               model.AlertCodeNoData,
	AlertSeverity:           model.AlertSeverityCritical,
	AlertStatus:             model.AlertStatusNew,
	Date:                    utils.P("2019-11-05T18:02:03Z"),
	Description:             "No data received from the host myhost in the last 90 days",
	OtherInfo: map[string]interface{}{
		"hostname": "myhost",
	},
	ID: utils.Str2oid("5dd4113f0085a6fac03c4fed"),
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
