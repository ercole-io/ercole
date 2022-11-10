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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (as *APIService) GetNodes(user *model.User) ([]model.Node, error) {
	if user == nil {
		return nil, utils.ErrInvalidUser
	}

	userRoles := make([]string, 0)

	for _, groupName := range user.Groups {
		group, err := as.Database.GetGroup(groupName)
		if err != nil {
			return nil, err
		}

		userRoles = append(userRoles, group.Roles...)
	}

	nodes, err := as.Database.GetNodesByRoles(userRoles)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}
