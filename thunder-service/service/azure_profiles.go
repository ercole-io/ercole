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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *ThunderService) AddAzureProfile(profile model.AzureProfile) (*model.AzureProfile, error) {
	profile.ID = as.NewObjectID()

	err := as.Database.AddAzureProfile(profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}
func (as *ThunderService) UpdateAzureProfile(profile model.AzureProfile) (*model.AzureProfile, error) {
	if err := as.Database.UpdateAzureProfile(profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
func (as *ThunderService) GetAzureProfiles() ([]model.AzureProfile, error) {
	azure_profile, err := as.Database.GetAzureProfiles(true)
	if err != nil {
		return nil, err
	}

	return azure_profile, nil
}

func (as *ThunderService) GetMapAzureProfiles() (map[primitive.ObjectID]model.AzureProfile, error) {
	azure_profile_with_id, err := as.Database.GetMapAzureProfiles()
	if err != nil {
		return nil, err
	}

	return azure_profile_with_id, nil
}

func (as *ThunderService) DeleteAzureProfile(id primitive.ObjectID) error {
	if err := as.Database.DeleteAzureProfile(id); err != nil {
		return err
	}

	return nil
}

func (as *ThunderService) SelectAzureProfile(profileId string, selected bool) error {
	if err := as.Database.SelectAzureProfile(profileId, selected); err != nil {
		return err
	}

	return nil
}
