// Copyright (c) 2023 Sorint.lab S.p.A.
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

package service

import (
	"reflect"
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/imdario/mergo"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hds *HostDataService) SaveExadata(exadata *model.OracleExadataInstance) error {
	existingExadata, err := hds.Database.FindExadataByRackID(exadata.RackID)
	if err != nil {
		if err != mongo.ErrNoDocuments {
			return err
		}
	}

	if existingExadata != nil && !reflect.DeepEqual(existingExadata, exadata) {
		if err := mergo.Merge(existingExadata, exadata); err != nil {
			return err
		}

		return hds.Database.UpdateExadata(*exadata)
	}

	exadata.CreatedAt = &time.Time{}

	return hds.Database.AddExadata(*exadata)
}
