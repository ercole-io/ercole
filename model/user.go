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

package model

import (
	"time"

	"github.com/ercole-io/ercole/v2/utils"
)

const SuperUser = "ercole"

type User struct {
	Username  string     `json:"username" bson:"username"`
	Password  string     `json:"password,omitempty" bson:"password"`
	Salt      string     `json:"-" bson:"salt"`
	LastLogin *time.Time `json:"lastLogin,omitempty" bson:"lastLogin"`
	FirstName string     `json:"firstName,omitempty" bson:"firstName"`
	LastName  string     `json:"lastName,omitempty" bson:"lastName"`
	Groups    []string   `json:"groups" bson:"groups"`
}

func (u *User) IsGroup(group string) bool {
	return utils.Contains(u.Groups, group)
}

func (u *User) IsAdmin() bool {
	return u.IsGroup(GroupAdmin)
}
