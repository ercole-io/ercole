// Copyright (c) 2019 Sorint.lab S.p.A.
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

import (
	"sync"

	"github.com/ercole-io/ercole/config"
)

// RepoService is a concrete implementation of RepoServiceInterface
type RepoService struct {
	// Config contains the reposervice global configuration
	Config config.Configuration
	// SubServices contains all subservices
	SubServices []SubRepoServiceInterface
}

// Init initialize all services
func (rs *RepoService) Init(wg *sync.WaitGroup) {
	for _, s := range rs.SubServices {
		s.Init(wg)
	}
}
