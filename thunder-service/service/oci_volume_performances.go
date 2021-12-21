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

func (as *ThunderService) AddOciVolumePerformance(volumePerformance model.OciVolumePerformance) (*model.OciVolumePerformance, error) {
	volumePerformance.ID = as.NewObjectID()
	err := as.Database.AddOciVolumePerformance(volumePerformance)
	if err != nil {
		return nil, err
	}
	return &volumePerformance, nil
}

func (as *ThunderService) UpdateOciVolumePerformance(profile model.OciVolumePerformance) (*model.OciVolumePerformance, error) {
	if err := as.Database.UpdateOciVolumePerformance(profile); err != nil {
		return nil, err
	}

	return &profile, nil
}
func (as *ThunderService) GetOciVolumePerformances() ([]model.OciVolumePerformance, error) {
	oracle_cloud_profile, err := as.Database.GetOciVolumePerformances()
	if err != nil {
		return nil, err
	}
	return oracle_cloud_profile, nil
}

func (as *ThunderService) getOciVolumePerformance(vpu int, size int) (*model.OciVolumePerformance, error) {
	volPerf, err := as.Database.GetOciVolumePerformance(vpu, size)
	if err != nil {
		return nil, err
	}
	return volPerf, nil
}

func (as *ThunderService) DeleteOciVolumePerformance(id primitive.ObjectID) error {
	if err := as.Database.DeleteOciVolumePerformance(id); err != nil {
		return err
	}
	return nil
}
