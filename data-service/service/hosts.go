// Copyright (c) 2021 Sorint.lab S.p.A.
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
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	patch_lib "github.com/ercole-io/ercole/v2/utils/patch"
)

// UpdateHostInfo saves the hostdata
func (hds *HostDataService) InsertHostData(hostdata model.HostDataBE) error {
	var err error

	hostdata.ServerVersion = hds.ServerVersion
	hostdata.Archived = false
	hostdata.CreatedAt = hds.TimeNow()
	hostdata.ServerSchemaVersion = model.SchemaVersion
	hostdata.ID = primitive.NewObjectIDFromTimestamp(hds.TimeNow())

	if hds.Config.DataService.EnablePatching {
		hostdata, err = hds.patchHostData(hostdata)
		if err != nil {
			return err
		}
	}

	previousHostdata, err := hds.Database.FindMostRecentHostDataOlderThan(hostdata.Hostname, hostdata.CreatedAt)
	if err != nil {
		hds.Log.Error(err)
		return err
	}

	if previousHostdata != nil && previousHostdata.Archived {
		// if the last one is archived, it was dismissed
		// we want to behave like there isn't any previous one
		previousHostdata = nil
	}

	if previousHostdata == nil {
		if err := hds.throwNewServerAlert(hostdata.Hostname); err != nil {
			return err
		}
	}

	if hostdata.Features.Oracle != nil {
		hds.oracleDatabasesChecks(previousHostdata, &hostdata)
	}

	if hostdata.Clusters != nil {
		hds.clusterInfoChecks(hostdata.Clusters)
	}

	_, err = hds.Database.DismissHost(hostdata.Hostname)
	if err != nil {
		return err
	}

	if hds.Config.DataService.LogInsertingHostdata {
		hds.Log.Info(utils.ToJSON(hostdata))
	}

	_, err = hds.Database.InsertHostData(hostdata)
	if err != nil {
		return err
	}

	if err := hds.Database.DeleteNoDataAlertByHost(hostdata.Hostname); err != nil {
		hds.Log.Error(err)
	}

	if hostdata.Errors != nil && len(hostdata.Errors) > 0 {
		if err := hds.throwAgentErrorsAlert(hostdata.Hostname, hostdata.Errors); err != nil {
			hds.Log.Error(err)
		}
	}

	return nil
}

// patchHostData patch the hostdata using the pf stored in the db
func (hds *HostDataService) patchHostData(hostdata model.HostDataBE) (model.HostDataBE, error) {
	patch, err := hds.Database.FindPatchingFunction(hostdata.Hostname)
	if err != nil {
		return model.HostDataBE{}, err
	}

	if patch.Hostname == hostdata.Hostname && patch.Code != "" {
		if hds.Config.DataService.LogDataPatching {
			hds.Log.Infof("Patching %s hostdata with the patch %s\n", patch.Hostname, patch.ID)
		}

		return patch_lib.PatchHostdata(patch, hostdata)
	}

	return hostdata, nil
}

func (hds *HostDataService) AlertInvalidHostData(validationErr error, hostdata *model.HostDataBE) {
	errs := make([]model.AgentError, 0)
	errs = append(errs, model.NewAgentError(validationErr))

	hostname := ""
	if hostdata != nil && hostdata.Errors != nil && len(hostdata.Errors) > 0 {
		hostname = hostdata.Hostname
		errs = append(errs, hostdata.Errors...)
	}

	if err := hds.throwAgentErrorsAlert(hostname, errs); err != nil {
		hds.Log.Error(err)
	}
}
