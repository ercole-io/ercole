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

import java.util.Calendar;
import java.util.Date;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.repository.query.Param;
import org.springframework.web.bind.annotation.DeleteMapping;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.RestController;

import io.ercole.model.CurrentHost;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.services.HostService;

/**
 * RestController for CRUD actions.
 */
@RestController
public class HostController {

	@Autowired
	private CurrentHostRepository repository;

	@Autowired
	private HostService hostService;
	
	/**
	 * 
	 * @param ricerca
	 *            a
	 * @param days
	 *            a
	 * @param c
	 *            a
	 * @return adsada
	 */
	@GetMapping(produces = "application/json", value = "/hosts")
	public Page<CurrentHost> getHosts(
			@Param("ricerca") final String ricerca, 
			@Param("days") final Integer days, 
			final Pageable c) {

		int d = 0;
		if (days != null) {
			d = days;
		}

		Calendar calendar = Calendar.getInstance();
		calendar.set(Calendar.HOUR_OF_DAY, 23);
		calendar.set(Calendar.MINUTE, 59);
		calendar.add(Calendar.DATE, -d);
		Date time = calendar.getTime();
		return repository.findByDbOrByHostnameOrBySchema(ricerca, time, c);

	}

	/**
	 * Delete a host.
	 * @param hostname host to delete
	 */
	@DeleteMapping(value = "/hosts/{hostname}")
	public void archiveHost(@PathVariable final String hostname) {
		hostService.archiveHost(hostname);
	}
}
