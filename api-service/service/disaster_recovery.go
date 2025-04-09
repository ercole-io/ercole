// Copyright (c) 2025 Sorint.lab S.p.A.
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

import "fmt"

func (as *APIService) CreateDR(hostname string, clusterVeritasHostnames []string) (string, error) {
	for i := 0; i < len(clusterVeritasHostnames); i++ {
		clusterVeritasHostnames[i] = fmt.Sprintf("%s_DR", clusterVeritasHostnames[i])
	}

	return as.Database.CreateDR(hostname, clusterVeritasHostnames)
}
