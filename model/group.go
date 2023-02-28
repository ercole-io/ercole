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

import "github.com/ercole-io/ercole/v2/utils"

const (
	// limited can only change the password
	GroupLimited = "limited"
	GroupAdmin   = "admin"
)

type Group struct {
	Name        string   `json:"name" bson:"name"`
	Description string   `json:"description" bson:"description"`
	Roles       []string `json:"roles" bson:"roles"`
	Tags        []string `json:"tags" bson:"tags"`
}

func (g *Group) IsRole(role string) bool {
	return utils.Contains(g.Roles, role)
}

func (g *Group) IsTag(tag string) bool {
	return utils.Contains(g.Tags, tag)
}
