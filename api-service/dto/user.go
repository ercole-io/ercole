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

package dto

import (
	"github.com/ercole-io/ercole/v2/model"
)

type User struct {
	Username  string   `json:"username"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Groups    []string `json:"groups"`
}

type Users []User

func ToUser(userModel *model.User) User {
	if userModel != nil {
		if userModel.Groups == nil {
			userModel.Groups = []string{}
		}

		return User{
			Username:  userModel.Username,
			FirstName: userModel.FirstName,
			LastName:  userModel.LastName,
			Groups:    userModel.Groups,
		}
	}

	return User{}
}

func ToUsers(usersModel []model.User) Users {
	result := make([]User, 0, len(usersModel))

	for _, userModel := range usersModel {
		if userModel.Groups == nil {
			userModel.Groups = []string{}
		}

		result = append(result, ToUser(&userModel))
	}

	return result
}
