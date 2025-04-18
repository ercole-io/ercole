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
package service

import (
	"fmt"
	"strings"
)

func (hds *HostDataService) getVeritasHostsFqdn(hostname string, clusterHosts []string) []string {
	parts := strings.Split(hostname, ".")

	if len(parts) <= 1 {
		return clusterHosts
	}

	fqdn := strings.TrimPrefix(hostname, parts[0])
	res := make([]string, 0, len(clusterHosts))

	for _, clusterHost := range clusterHosts {
		if !strings.HasSuffix(clusterHost, fqdn) {
			clusterHost = fmt.Sprintf("%s%s", clusterHost, fqdn)
		}

		res = append(res, clusterHost)
	}

	return res
}
