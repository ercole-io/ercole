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

package io.ercole.model;

import java.util.Date;

import javax.persistence.Entity;
import javax.validation.constraints.NotEmpty;

/**
 * Object mapped with JPA in database.
 */

@Entity
public class CurrentHost extends Host {

	/**
	 * Instantiates a new current host.
	 */
	public CurrentHost() {
		super();
	}
	
	/**
	 * @param id 
	 * @param hostname 
	 * @param environment 
	 * @param location 
	 * @param version the version
	 * @param serverVersion the serverVersion
	 * @param hostType The HostType
	 * @param databases 
	 * @param schemas 
	 * @param extraInfo 
	 * @param associatedClusterName associated cluster name
	 * @param associatedHypervisorHostname associated hypervisor hostname
	 * @param hostInfo 
	 * @param updated 
	 */
	public CurrentHost(final Long id, final @NotEmpty String hostname, final String environment,
			final String location, final String version, final String serverVersion, final String hostType,
			final String databases, final String schemas, final String extraInfo,
			final String associatedClusterName, final String associatedHypervisorHostname, final String hostInfo, final Date updated) {
		super(id, hostname, environment, location, version, serverVersion, hostType, databases, schemas,
				extraInfo, associatedClusterName, associatedHypervisorHostname, hostInfo, updated);
	}
}
