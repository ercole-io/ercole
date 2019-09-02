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

import java.io.IOException;
import java.math.BigDecimal;
import java.util.ArrayList;
import java.util.Collection;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.http.ResponseEntity;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.context.SecurityContextHolder;
import org.springframework.security.core.userdetails.UserDetails;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.PathVariable;
import org.springframework.web.bind.annotation.PutMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.bind.annotation.RestController;

import io.ercole.model.ClusterInfo;
import io.ercole.services.DashboardService;
import io.ercole.services.GenerateExcelAddmService;
import io.ercole.services.GenerateExcelPatchService;
import io.ercole.services.GenerateExcelSegmentService;
import io.ercole.services.GenerateExcelService;
import io.ercole.services.HostService;

/**
 * Controller for user actions.
 */
@RestController
public class UserController {

	@Autowired
	private DashboardService dashService;

	@Autowired
	private GenerateExcelService genService;

	@Autowired
	private GenerateExcelPatchService generateExcelPatchService;

	@Autowired
	private GenerateExcelSegmentService generateExcelSegmentService;
	@Autowired
	private GenerateExcelAddmService generateExcelAddmService;

	@Autowired
	private HostService hostService;


	/**
	 * Acknowledge alerts.
	 *
	 * @param array of Alert IDs to aknowledge
	 */
	@PutMapping(value = "/acknowledge", consumes = "application/json")
	public void acknowledgeAlerts(@RequestBody final Long[] array) {
		dashService.acknowledgeAlert(array);
	}

	/**
	 * Generate excel response entity.
	 *
	 * @return the response entity
	 * @throws IOException the io exception
	 */
	@GetMapping(path = "/generateexcel")
	public ResponseEntity<byte[]> generateExcel() throws IOException {
		return genService.initExcel();
	}


	/**
	 * Gets server env.
	 *
	 * @param location the location filter
	 * @return JSONArray of couples <ServerEnvironment, CounterInteger>
	 */
	@GetMapping(value = "/getserverenv")
	public List<Map<String, Object>> getServerEnv(@RequestParam final String location) {
		return dashService.getServerEnv(location);
	}

	/**
	 * Gets server location.
	 *
	 * @return JSONArray of distinct server locations
	 */
	@GetMapping(value = "/getserverslocations")
	public List<String> getServerLocation() {
		return dashService.getServerLocation();
	}

	/**
	 * Gets db feature count.
	 *
	 * @param location the location filter
	 * @return JSONArray of distinct features and counter for usage
	 */
	@GetMapping(value = "/getdbfeatures")
	public List<Map<String, Object>> getDbFeatureCount(@RequestParam final String location) {
		return dashService.getDbFeaturesCount(location);
	}

	/**
	 * Gets db env.
	 *
	 * @param location the location filter
	 * @return JSONArray of "count of DBs for each Environment"
	 */
	@GetMapping(value = "/getdbenv")
	public List<Map<String, Object>> getDbEnv(@RequestParam final String location) {
		return dashService.getDbEnv(location);
	}

	/**
	 * Gets db versions count.
	 *
	 * @param location the location filter
	 * @return a JSON Array of DB Versions+Counter
	 */
	@GetMapping(value = "/getdbversions")
	public List<Map<String, Object>> getDbVersionsCount(@RequestParam final String location) {
		return dashService.getDbVersionsCount(location);
	}

	/**
	 * Gets host type count.
	 *
	 * @param location the location filter
	 * @return a JSON Array of HostTypes+Counter
	 */
	@GetMapping(value = "/gethosttypes")
	public List<Map<String, Object>> getHostTypeCount(@RequestParam final String location) {
		return dashService.getHostTypeCount(location);
	}

	/**
	 * Gets os type count.
	 *
	 * @param location the location filter
	 * @return a JSON Array of OS Types+Counter
	 */
	@GetMapping(value = "/getostypes")
	public List<Map<String, Object>> getOsTypeCount(@RequestParam final String location) {
		return dashService.getOsTypeCount(location);
	}


	/**
	 * Gets resize.
	 *
	 * @param location the location filter
	 * @return a JSON Array of Top unused CPUs with Hostnames
	 */
	@GetMapping(value = "/getresize")
	public List<Map<String, Object>> getResize(@RequestParam final String location) {
		return dashService.getResize(location);
	}

	/**
	 * Gets work by dbs.
	 *
	 * @param location the location filter
	 * @return a JSON Array of Top Work by DBs
	 */
	@GetMapping(value = "/getwork")
	public List<Map<String, Object>> getWorkByDbs(@RequestParam final String location) {
		return dashService.getWorkByDbs(location);
	}

	/**
	 * Gets ent std licenses count.
	 *
	 * @param location the location filter
	 * @return a JSON Array of counted Standard & Enterprise Licenses
	 */
	@GetMapping(value = "/getlicenses")
	List<Map<BigDecimal, String>> getEntStdLicensesCount(@RequestParam final String location) {
		return dashService.getLicensesCount(location);
	}

	/**
	 * Gets all licenses count.
	 *
	 * @param location the location filter
	 * @return a JSON Array of counted All kind of Licenses
	 */
	@GetMapping(value = "/getalllicenses")
	List<Map<String, Object>> getAllLicensesCount(@RequestParam final String location) {
		return dashService.getAllLicensesCount(location);
	}


	/**
	 * Gets all host using license.
	 *
	 * @param license the license
	 * @return a JSON Array of counted All kind of Licenses
	 */
	@GetMapping(value = "/getallhostusinglicense")
	List<Map<String, Object>> getAllHostUsingLicense(@RequestParam final String license) {
		return dashService.getAllHostUsingLicense(license);
	}


	/**
	 * Gets compliance.
	 *
	 * @return (" Compliance " = true) if there are more Licenses than used, ("Compliance"=false) if there are more used Licenses than charged
	 */
	@GetMapping(value = "/getcompliance")
	Map<String, Object> getCompliance() {
		return dashService.getCompliance();
	}

	/**
	 * Get user informations.
	 *
	 * @return an object containing user informations
	 */
	@GetMapping("/whoami")
	public final Map<String, Object> getUser() {
		UserDetails principal = //
			(UserDetails) SecurityContextHolder.getContext().getAuthentication().getPrincipal();
		
		Map<String, Object> userInformations = new HashMap<>();
		userInformations.put("username", principal.getUsername());
		// fetch autorities
		Collection<GrantedAuthority> authorities = (Collection<GrantedAuthority>) principal.getAuthorities();
		List<String> auths = new ArrayList<>();
		for (GrantedAuthority s : authorities) {
			auths.add(s.getAuthority());
		}
		
		userInformations.put("authorities", auths);
		userInformations.put("fullName", "");
		
		return userInformations;
	}

	/**
	 * Get all addms.
	 *
	 * @param env    the env filter
	 * @param search the search filter
	 * @return the list of addms
	 */
	@GetMapping("/getalladdms")
	public final List<Map<String, Object>> getADDMs(
			@RequestParam final String env,
			@RequestParam final String search) {
		return dashService.getADDMs(env, search);
	}

	/**
	 * Generate excel addm response entity.
	 *
	 * @param env    the env
	 * @param search the search
	 * @return the response entity
	 * @throws IOException the io exception
	 */
	@GetMapping(path = "/generate-addm-excel")
	public ResponseEntity<byte[]> generateExcelAddm(
			@RequestParam final String env,
			@RequestParam final String search) throws IOException {
		return generateExcelAddmService.initExcel(env, search);
	}

	/**
	 * Get all segment advisors.
	 *
	 * @param env    the env filter
	 * @param search the search filter
	 * @return the list of segment advisors
	 */
	@GetMapping("/getallsegmentadvisors")
	public final List<Map<String, Object>> getSegmentAdvisors(
			@RequestParam final String env,
			@RequestParam final String search) {
		return dashService.getSegmentAdvisors(env, search);
	}


	/**
	 * Get the environments.
	 *
	 * @return the environments.
	 */
	@GetMapping("/environments")
	public List<String> getEnvironments() {
		return dashService.getEnvironments();
	}


	/**
	 * Gets patch advisors.
	 *
	 * @param status     the status
	 * @param windowTime the window time
	 * @return the patch advisors
	 * @throws IOException the io exception
	 */
	@GetMapping("/getallpatchadvisors")
	public final List<Map<String, Object>> getPatchAdvisors(
			@RequestParam final String status,
			@RequestParam final int windowTime) {
		return dashService.getPatchAdvisors(status, windowTime);
	}


	/**
	 * Generate excel patch response entity.
	 *
	 * @param windowTime the window time
	 * @param status     the status
	 * @return the response entity
	 * @throws IOException the io exception
	 */
	@GetMapping(path = "/generate-patch-excel")
	public ResponseEntity<byte[]> generateExcelPatch(
			@RequestParam final int windowTime,
			@RequestParam final String status) throws IOException {

		return generateExcelPatchService.initExcel(windowTime, status);
	}

	/**
	 * Generate excel segment response entity.
	 *
	 * @param env    the env
	 * @param search the search
	 * @return the response entity
	 * @throws IOException the io exception
	 */
	@GetMapping(path = "/generate-segment-excel")
	public ResponseEntity<byte[]> generateExcelSegment(
			@RequestParam final String env,
			@RequestParam final String search) throws IOException {
			return generateExcelSegmentService.initExcel(env, search);
	}

	/**
	 * Get the top 15 reclaimable database.
	 *
	 * @param location location
	 * @return Top 15 reclaimable databas
	 */
	@GetMapping("/gettopreclaimabledatabase")
	public final List<Map<String, Object>> getTopReclaimableDatabase(@RequestParam final String location) {
		return dashService.getTopReclaimableDatabase(location);
	}

	/**
	 * Get patch status stats.
	 * @param location location
	 * @param windowTime windowTime
	 * @return patch status stats
	 */
	@GetMapping("/getpatchstatusstats")
	public final List<Map<String, Object>> getPatchStatusStats(@RequestParam final String location, @RequestParam final int windowTime) {
		return dashService.getPatchStatusStats(location, windowTime);
	}

	/**
     * Return all cluster that match the filter.
     * @param filter filter
     * @return all cluster that match the filter
     */	
	@GetMapping("/getclusters")
	public final List<ClusterInfo> getCluster(@RequestParam final String filter) {
		return hostService.getClusters(filter);		
	}

	/**
	 * Return the used data history of all databases of host.
	 * @param hostname hostname
	 * @return the used data history of all databases of host
	 */
	@GetMapping("/hosts/{hostname}/useddatahistory")
	public final Map<String, Object> getUsedDataHistory(@PathVariable final String hostname) {
		return hostService.getUsedDataHistory(hostname);
	}

	/**
	 * Return the segmentsSize data history of all databases of host.
	 * @param hostname hostname
	 * @return the segmentsSize data history of all databases of host
	 */
	@GetMapping("/hosts/{hostname}/segmentssizedatahistory")
	public final Map<String, Object> getSegmentsSizeDataHistory(@PathVariable final String hostname) {
		return hostService.getSegmentsSizeDataHistory(hostname);
	}
}
