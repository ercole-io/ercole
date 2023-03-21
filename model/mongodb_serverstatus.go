// Copyright (c) 2023 Sorint.lab S.p.A.
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

type ServerStatusConnection struct {
	Current                 int32 `json:"current" bson:"current"`
	Available               int32 `json:"available" bson:"available"`
	TotalCreated            int32 `json:"totalCreated" bson:"totalCreated"`
	Active                  int32 `json:"active" bson:"active"`
	Threaded                int32 `json:"threaded" bson:"threaded"`
	ExhaustIsMaster         int32 `json:"exhaustIsMaster" bson:"exhaustIsMaster"`
	ExhaustHello            int32 `json:"exhaustHello" bson:"exhaustHello"`
	AwaitingTopologyChanges int32 `json:"awaitingTopologyChanges" bson:"awaitingTopologyChanges"`
	LoadBalanced            int32 `json:"loadBalanced" bson:"loadBalanced"`
}
