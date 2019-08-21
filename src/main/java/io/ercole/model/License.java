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

import javax.persistence.Entity;
import javax.persistence.Id;

/**
 * The Class Alert.
 */
@Entity
public class License {


	/** The id. */
	@Id
	private String id;

	/** The number of owned licenses. */
	private Long licenseCount;

	/**
	 * Instantiates a new license.
	 */
	public License() {
		// used for testing and JPA purposes
	}

	/**
	 * Gets the id.
	 *
	 * @return the id
	 */
	public String getId() {
		return id;
	}

	/**
	 * Sets the id.
	 *
	 * @param id the new id
	 */
	public void setId(final String id) {
		this.id = id;
	}

	/**
	 * Gets the number of owned licenses.
	 *
	 * @return the license count
	 */
	public Long getLicenseCount() {
		return licenseCount;
	}

	/**
	 * Sets the number of owned licenses.
	 *
	 * @param licenseCount the new license count
	 */
	public void setLicenseCount(final Long licenseCount) {
		this.licenseCount = licenseCount;
	}
	
	

}
