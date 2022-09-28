// Copyright (c) 2022 Sorint.lab S.p.A.
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

func (as *APIService) InsertGroup(group model.GroupType) (*model.GroupType, error) {
	group.ID = as.NewObjectID()

	err := as.Database.InsertGroup(group)
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (as *APIService) UpdateGroup(group model.GroupType) (*model.GroupType, error) {
	if err := as.Database.UpdateGroup(group); err != nil {
		return nil, err
	}

	return &group, nil
}

func (as *APIService) GetGroup(id primitive.ObjectID) (*model.GroupType, error) {
	group, err := as.Database.GetGroup(id)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (as *APIService) GetGroups() ([]model.GroupType, error) {
	groups, err := as.Database.GetGroups()
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (as *APIService) DeleteGroup(id primitive.ObjectID) error {
	if err := as.Database.DeleteGroup(id); err != nil {
		return err
	}

	return nil
}
