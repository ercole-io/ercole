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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *ThunderService) AddOciProfile(profile model.OciProfile) (*model.OciProfile, error) {
	profile.ID = as.NewObjectID()
	err := as.Database.AddOciProfile(profile)
	if err != nil {
		return nil, err
	}
	return &profile, nil
}

func (as *ThunderService) UpdateOciProfile(profile model.OciProfile) (*model.OciProfile, error) {
	if err := as.Database.UpdateOciProfile(profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
func (as *ThunderService) GetOciProfiles() ([]model.OciProfile, error) {
	oracle_cloud_profile, err := as.Database.GetOciProfiles(true)
	if err != nil {
		return nil, err
	}
	return oracle_cloud_profile, nil
}

func (as *ThunderService) GetMapOciProfiles() (map[primitive.ObjectID]model.OciProfile, error) {
	oracle_cloud_profile_with_id, err := as.Database.GetMapOciProfiles()
	if err != nil {
		return nil, err
	}
	return oracle_cloud_profile_with_id, nil
}

func (as *ThunderService) DeleteOciProfile(id primitive.ObjectID) error {
	if err := as.Database.DeleteOciProfile(id); err != nil {
		return err
	}
	return nil
}
