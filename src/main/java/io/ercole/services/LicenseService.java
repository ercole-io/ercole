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

package io.ercole.services;

import java.util.List;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.stereotype.Service;

import io.ercole.model.License;
import io.ercole.repositories.LicenseRepository;

/**
 * Service component for License.
 */
@Service
public class LicenseService {

	@Autowired
	private LicenseRepository licenseRepo;

	/**
	 * Update licenses.
	 * 
	 * @param licenses
	 *            the list of licenses to update
	 * 
	 * @return a list of updated licenses
	 */
	public Iterable<License> updateLicenses(final List<License> licenses) {
		Iterable<License> repoLicenses = licenseRepo.findAll();
		for (License repoLicense : repoLicenses) {
			for (License license : licenses) {
				if (license.getId().equals(repoLicense.getId())) {
					repoLicense.setLicenseCount(license.getLicenseCount());
				}
			}
		}
		return licenseRepo.saveAll(repoLicenses);
	}

}
