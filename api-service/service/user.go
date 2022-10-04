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

package service

import (
	"encoding/json"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/schema"
	cr "github.com/ercole-io/ercole/v2/utils/crypto"
)

func (as *APIService) ListUsers() ([]model.User, error) {
	return as.Database.ListUsers()
}

func (as *APIService) GetUser(username string) (*model.User, error) {
	return as.Database.GetUser(username)
}

func (as *APIService) AddUser(user model.User) error {
	salt, err := cr.GenerateRandomBytes()
	if err != nil {
		return err
	}

	user.Password, user.Salt = cr.GenerateHashAndSalt(user.Password, salt)

	raw, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := schema.ValidateUser(raw); err != nil {
		return err
	}

	return as.Database.AddUser(user)
}

func (as *APIService) UpdateUserGroups(updatedUser model.User) error {
	return as.Database.UpdateUserGroups(updatedUser)
}

func (as *APIService) RemoveUser(username string) error {
	return as.Database.RemoveUser(username)
}
