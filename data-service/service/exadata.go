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
	"errors"

	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/mongo"
)

func (hds *HostDataService) SaveExadata(exadata *model.OracleExadataInstance) error {
	existingExadata, err := hds.Database.FindExadataByRackID(exadata.RackID)
	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	if existingExadata != nil {
		return hds.updateExistingExadata(existingExadata, exadata)
	}

	return hds.addNewExadata(exadata)
}

func (hds *HostDataService) updateExistingExadata(existingExadata, newExadata *model.OracleExadataInstance) error {
	if existingExadata.Hostname != newExadata.Hostname {
		if err := hds.Database.UpdateExadataHostname(existingExadata.RackID, newExadata.Hostname); err != nil {
			return err
		}
	}

	existingComponents, err := hds.getExistingExadataComponent(existingExadata, newExadata)
	if err != nil {
		return err
	}

	for _, ec := range existingComponents {
		if err := hds.Database.SetExadataComponent(ec.RackID, ec); err != nil {
			return err
		}
	}

	newComponents, err := hds.getNewExadataComponent(existingExadata, newExadata)
	if err != nil {
		return err
	}

	for _, component := range newComponents {
		if err := hds.Database.PushComponentToExadataInstance(newExadata.RackID, component); err != nil {
			return err
		}
	}

	return nil
}

func (hds *HostDataService) addNewExadata(exadata *model.OracleExadataInstance) error {
	exadata.CreatedAt, exadata.UpdatedAt = hds.TimeNow(), hds.TimeNow()
	return hds.Database.AddExadata(*exadata)
}

func (hds *HostDataService) getNewExadataComponent(old, new *model.OracleExadataInstance) ([]model.OracleExadataComponent, error) {
	if old == nil {
		return nil, errors.New("old exadata instance cannot be nil during comparision")
	}

	if new == nil {
		return nil, errors.New("new exadata instance cannot be nil during comparision")
	}

	if old.RackID != new.RackID {
		return nil, errors.New("cannot compare different exadata instances")
	}

	newInstances := make([]model.OracleExadataComponent, 0)

	hostnames := make(map[string]bool)

	for _, oc := range old.Components {
		hostnames[oc.Hostname] = true
	}

	for _, nc := range new.Components {
		if _, ok := hostnames[nc.Hostname]; !ok {
			newInstances = append(newInstances, nc)
		}
	}

	return newInstances, nil
}
func (hds *HostDataService) getExistingExadataComponent(old, new *model.OracleExadataInstance) ([]model.OracleExadataComponent, error) {
	if old == nil {
		return nil, errors.New("old exadata instance cannot be nil during comparision")
	}

	if new == nil {
		return nil, errors.New("new exadata instance cannot be nil during comparision")
	}

	if old.RackID != new.RackID {
		return nil, errors.New("cannot compare different exadata instances")
	}

	existingComponents := make([]model.OracleExadataComponent, 0)

	hostnames := make(map[string]bool)

	for _, oc := range old.Components {
		hostnames[oc.Hostname] = true
	}

	for _, nc := range new.Components {
		if _, ok := hostnames[nc.Hostname]; ok {
			existingComponents = append(existingComponents, nc)
		}
	}

	return existingComponents, nil
}
