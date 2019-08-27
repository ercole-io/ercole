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

import org.springframework.data.rest.core.config.Projection;

/**
 * Projection for not visualizing extraInfo attribute of CurrentHost objects.
 *
 */
@Projection(name = "noExtraInfo", types = { CurrentHost.class })
public interface NoExtraInfo {
	
	/**
	 * @return hostname
	 */
	String getHostname();
	/**
	 * @return the environment
	 */
	String getEnvironment();
	/**
	 * @return the location
	 */
	String getLocation();
	/**
	 * @return the version
	 */
	String getVersion();

		/**
	 * @return the serverVersion
	 */
	String getServerVersion();

	/**
	 * @return the hostType
	 */
	String getHostType();

	/**
	 * @return the databases
	 */
	String getDatabases();
	
	/**
	 * @return the schemas
	 */
	String getSchemas();

	/**
	 * @return JSON infos regarding host OS and HW
	 */
	String getHostInfo();

}
