// Copyright (c) 2020 Sorint.lab S.p.A.
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

// Package service is a package that contains varios file serving services
package service

import "sync"

// RepoServiceInterface is a interface that wrap functions for starting a service
type RepoServiceInterface interface {
	// Init initialize the service
	Init(wg *sync.WaitGroup)
}

// SubRepoServiceInterface is a interface that wrap functions for starting a subservice
type SubRepoServiceInterface interface {
	// Init initialize the service
	Init(wg *sync.WaitGroup)
}
