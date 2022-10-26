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
	"encoding/json"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *APIService) GetRole(name string) (*model.Role, error) {
	role, err := as.Database.GetRole(name)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (as *APIService) GetRoles() ([]model.Role, error) {
	roles, err := as.Database.GetRoles()
	if err != nil {
		return nil, err
	}

	return roles, nil
}

func (as *APIService) AddRole(role model.Role) error {
	locations, err := as.ListLocations("", "", utils.MAX_TIME)
	if err != nil {
		return err
	}

	if !utils.Contains(locations, role.Location) {
		return utils.ErrInvalidLocation
	}

	raw, err := json.Marshal(role)
	if err != nil {
		return err
	}

	if err := schema.ValidateRole(raw); err != nil {
		return err
	}

	return as.Database.AddRole(role)
}

func (as *APIService) UpdateRole(role model.Role) error {
	locations, err := as.ListLocations("", "", utils.MAX_TIME)
	if err != nil {
		return err
	}

	if !utils.Contains(locations, role.Location) {
		return utils.ErrInvalidLocation
	}

	raw, err := json.Marshal(role)
	if err != nil {
		return err
	}

	if err := schema.ValidateRole(raw); err != nil {
		return err
	}

	documents := bson.D{
		primitive.E{Key: "description", Value: role.Description},
		primitive.E{Key: "location", Value: role.Location},
		primitive.E{Key: "permission", Value: role.Permission},
	}

	if err := as.Database.UpdateRole(role.Name, documents); err != nil {
		return err
	}

	return nil
}

func (as *APIService) RemoveRole(roleName string) error {
	return as.Database.RemoveRole(roleName)
}
