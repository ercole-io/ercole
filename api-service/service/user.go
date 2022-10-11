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
	"encoding/base64"
	"encoding/json"
	"errors"

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

	user.Groups = append(user.Groups, model.GroupLimited)

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

func (as *APIService) NewPassword(username string) (string, error) {
	saltByte, err := cr.GenerateRandomBytes()
	if err != nil {
		return "", err
	}

	suggestedPassword := cr.SuggestPassword()

	hashPwd, salt := cr.GenerateHashAndSalt(suggestedPassword, saltByte)

	if err := as.Database.UpdatePassword(username, hashPwd, salt); err != nil {
		return "", err
	}

	return suggestedPassword, nil
}

func (as *APIService) MatchPassword(user *model.User, password string) bool {
	if user == nil {
		return false
	}

	salt, err := base64.RawStdEncoding.DecodeString(user.Salt)
	if err != nil {
		return false
	}

	pwd, _ := cr.GenerateHashAndSalt(password, salt)

	return pwd == user.Password
}

func (as *APIService) UpdatePassword(username string, oldPass string, newPass string) error {
	user, err := as.GetUser(username)
	if err != nil {
		return err
	}

	if ok := as.MatchPassword(user, oldPass); !ok {
		return errors.New("Invalid password")
	}

	saltByte, err := cr.GenerateRandomBytes()
	if err != nil {
		return err
	}

	hashPwd, salt := cr.GenerateHashAndSalt(newPass, saltByte)

	if err := as.Database.UpdatePassword(username, hashPwd, salt); err != nil {
		return err
	}

	return nil
}
