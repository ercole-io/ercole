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

package service

import (
	"strings"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (hds *HostDataService) clusterInfoChecks(clusters []model.ClusterInfo) {
	hds.assignKnownHostnames(clusters)
}

func (hds *HostDataService) assignKnownHostnames(clusters []model.ClusterInfo) {
	var hostnames map[string]string
	{
		knownHostnames, err := hds.Database.GetHostnames()
		if err != nil {
			hds.Log.Error(utils.NewError(err, "Can't retrieve hostnames"))
			return
		}

		hostnames = make(map[string]string, len(knownHostnames))
		for _, knownHostname := range knownHostnames {
			alternatives := []string{
				strings.ToLower(knownHostname),
				strings.ToLower(strings.Split(knownHostname, ".")[0]),
			}

			for _, a := range alternatives {
				hostnames[a] = knownHostname
			}
		}
	}

	for i := range clusters {
		cluster := &clusters[i]
	vms:
		for j := range cluster.VMs {
			vm := &cluster.VMs[j]

			alternatives := []string{
				strings.ToLower(vm.Hostname),
				strings.ToLower(strings.Split(vm.Hostname, ".")[0]),
			}
			for _, a := range alternatives {
				if knownHostname, ok := hostnames[a]; ok {
					vm.Hostname = knownHostname
					continue vms
				}
			}
		}
	}
}
