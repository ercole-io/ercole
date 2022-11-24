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

import "github.com/ercole-io/ercole/v2/model"

type Node struct {
	Name   string `json:"name"`
	Parent string `json:"parent"`
}

func ToNode(nodeModel *model.Node) Node {
	if nodeModel != nil {
		return Node{
			Name:   nodeModel.Name,
			Parent: nodeModel.Parent,
		}
	}
	
	return Node{}
}

func ToNodes(nodesModel []model.Node) []Node {
	result := make([]Node, 0, len(nodesModel))

	for _, nodeModel := range nodesModel {
		result = append(result, ToNode(&nodeModel))
	}

	return result
}
