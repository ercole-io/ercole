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
)

func (as *APIService) GetNodes(groups []string) ([]model.Node, error) {
	userRoles := make([]string, 0)

	for _, groupName := range groups {
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

	if as.Config.APIService.EnableOciMenu || as.Config.APIService.EnableAwsMenu || as.Config.APIService.EnableGcpMenu {
		nodes = append(nodes, model.Node{
			Name: "Cloud Advisors",
			Roles: []string{
				"admin",
				"read_cloud",
			},
			Parent: "",
		})
	}

	if as.Config.APIService.EnableOciMenu {
		nodes = append(nodes,
			model.Node{
				Name: "OCI",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "Cloud Advisors",
			},
			model.Node{
				Name: "Profile Configurations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "OCI",
			},
			model.Node{
				Name: "Recommendations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "OCI",
			},
		)
	}

	if as.Config.APIService.EnableAwsMenu {
		nodes = append(nodes,
			model.Node{
				Name: "AWS",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "Cloud Advisors",
			},
			model.Node{
				Name: "Profile Configurations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "AWS",
			},
			model.Node{
				Name: "Recommendations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "AWS",
			},
		)
	}

	if as.Config.APIService.EnableGcpMenu {
		nodes = append(nodes,
			model.Node{
				Name: "GCP",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "Cloud Advisors",
			},
			model.Node{
				Name: "Profile Configurations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "GCP",
			},
			model.Node{
				Name: "Recommendations",
				Roles: []string{
					"admin",
					"read_cloud",
				},
				Parent: "GCP",
			},
		)
	}

	return nodes, nil
}

func (as *APIService) GetNode(name string) (*model.Node, error) {
	node, err := as.Database.GetNodeByName(name)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (as *APIService) AddNode(node model.Node) error {
	return as.Database.AddNode(node)
}

func (as *APIService) UpdateNode(node model.Node) error {
	return as.Database.UpdateNode(node)
}

func (as *APIService) RemoveNode(name string) error {
	return as.Database.RemoveNode(name)
}
