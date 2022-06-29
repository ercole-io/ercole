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

func (as *ThunderService) AddAwsProfile(profile model.AwsProfile) (*model.AwsProfile, error) {
	profile.ID = as.NewObjectID()

	err := as.Database.AddAwsProfile(profile)
	if err != nil {
		return nil, err
	}

	return &profile, nil
}

func (as *ThunderService) UpdateAwsProfile(profile model.AwsProfile) (*model.AwsProfile, error) {
	if err := as.Database.UpdateAwsProfile(profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
func (as *ThunderService) GetAwsProfiles() ([]model.AwsProfile, error) {
	aws_profile, err := as.Database.GetAwsProfiles(true)
	if err != nil {
		return nil, err
	}

	return aws_profile, nil
}

func (as *ThunderService) GetMapAwsProfiles() (map[primitive.ObjectID]model.AwsProfile, error) {
	aws_profile_with_id, err := as.Database.GetMapAwsProfiles()
	if err != nil {
		return nil, err
	}

	return aws_profile_with_id, nil
}

func (as *ThunderService) DeleteAwsProfile(id primitive.ObjectID) error {
	if err := as.Database.DeleteAwsProfile(id); err != nil {
		return err
	}

	return nil
}

func (as *ThunderService) SelectAwsProfile(profileId string, selected bool) error {
	if err := as.Database.SelectAwsProfile(profileId, selected); err != nil {
		return err
	}

	return nil
}
