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

package io.ercole.controller;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RestController;

import io.ercole.model.License;
import io.ercole.services.LicenseService;

/**
 * Controller for License.
 */
@RestController
public class LicenseController {

	@Autowired
	private LicenseService licenseService;
	
	/**
	 * Update licenses.
	 * 
	 * @param licenses list of licenses to update.
	 * 
	 * @return list of updated Licenses
	 * 
	 */
	@PutMapping(value = "/updatelicenses", consumes = "application/json")
	public Iterable<License> updateLicenses(@RequestBody final List<License> licenses) {
		return licenseService.updateLicenses(licenses);
	}
	

}
