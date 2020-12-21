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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"bytes"
	"net/http"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UpdateHostInfo saves the hostdata
func (hds *HostDataService) InsertHostData(hostdata model.HostDataBE) (interface{}, utils.AdvancedErrorInterface) {
	var aerr utils.AdvancedErrorInterface

	hostdata.ServerVersion = hds.ServerVersion
	hostdata.Archived = false
	hostdata.CreatedAt = hds.TimeNow()
	hostdata.ServerSchemaVersion = model.SchemaVersion
	hostdata.ID = primitive.NewObjectIDFromTimestamp(hds.TimeNow())

	if hds.Config.DataService.EnablePatching {
		hostdata, aerr = hds.PatchHostData(hostdata)
		if aerr != nil {
			return nil, aerr
		}
	}

	if hostdata.Features.Oracle != nil {
		hds.oracleDatabasesChecks(hostdata.Info, hostdata.Features.Oracle)
	}

	_, aerr = hds.Database.ArchiveHost(hostdata.Hostname)
	if aerr != nil {
		return nil, aerr
	}

	if hds.Config.DataService.LogInsertingHostdata {
		hds.Log.Info(utils.ToJSON(hostdata))
	}
	res, aerr := hds.Database.InsertHostData(hostdata)
	if aerr != nil {
		return nil, aerr
	}

	alertHostDataInsertionURL := utils.NewAPIUrlNoParams(
		hds.Config.AlertService.RemoteEndpoint,
		hds.Config.AlertService.PublisherUsername,
		hds.Config.AlertService.PublisherPassword,
		"/queue/host-data-insertion/"+res.InsertedID.(primitive.ObjectID).Hex()).String()

	if resp, err := http.Post(alertHostDataInsertionURL, "application/json", bytes.NewReader([]byte{})); err != nil {
		return nil, utils.NewAdvancedErrorPtr(err, "EVENT ENQUEUE")

	} else if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, utils.NewAdvancedErrorPtr(utils.ErrEventEnqueue, "EVENT ENQUEUE")
	}

	return res.InsertedID, nil
}

// PatchHostData patch the hostdata using the pf stored in the db
func (hds *HostDataService) PatchHostData(hostdata model.HostDataBE) (model.HostDataBE, utils.AdvancedErrorInterface) {
	patch, err := hds.Database.FindPatchingFunction(hostdata.Hostname)
	if err != nil {
		return model.HostDataBE{}, err
	}

	if patch.Hostname == hostdata.Hostname && patch.Code != "" {
		if hds.Config.DataService.LogDataPatching {
			hds.Log.Printf("Patching %s hostdata with the patch %s\n", patch.Hostname, patch.ID)
		}

		return utils.PatchHostdata(patch, hostdata)
	}

	return hostdata, nil
}
