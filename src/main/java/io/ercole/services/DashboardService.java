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

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.ArrayList;
import java.util.Calendar;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;
import java.util.regex.Pattern;

import javax.transaction.Transactional;

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import io.ercole.repositories.AlertRepository;
import io.ercole.repositories.CurrentHostRepository;

/**
 * The Component for presenting to user interface the needed data in JSON format.
 */
@Service
public class DashboardService {

	private Logger logger = LoggerFactory.getLogger(HostService.class);


	@Autowired
	private CurrentHostRepository currentRepo;

	@Autowired
	private AlertRepository alertRepo;

	@Value(value = "#{'${dashboard.aggregate.os.rules.regex}'.split('\n')}")
	private List<String> osAggregateRulesRegex;
	@Value(value = "#{'${dashboard.aggregate.os.rules.group}'.split('\n')}")
	private List<String> osAggregateRulesGroup;

	/**
	 * @param idArray
	 *            is an array of Alert IDs to aknowledge
	 * @return true if all IDs have been set to ACK
	 */
	public boolean acknowledgeAlert(final Long[] idArray) {
		int i = idArray.length;

		for (Long id : idArray) {
			if (alertRepo.setFromNewToAck(id) == 1) {
				i--;
			}
		}

		return i == 0;
	}

	/**
	 * @param location the location filter
	 * @return List of couples <ServerEnvironment, CounterInteger>
	 */
	public List<Map<String, Object>> getServerEnv(final String location) {
		return currentRepo.getServerTypeCount(location);
	}

	/**
	 * @return List of distinct server locations
	 */
	public List<String> getServerLocation() {
		return currentRepo.getLocations();
	}

	/**
	 * @param location the location filter
	 * @return count of DBs for each Environment
	 */
	public List<Map<String, Object>> getDbEnv(final String location) {
		return currentRepo.getDbEnvs(location);
	}

	/**
	 * @param location the location filter
	 * @return Feature and Counter
	 */
	public List<Map<String, Object>> getDbFeaturesCount(final String location) {
		List<String> queryVal = currentRepo.getDbFeaturesCount(location);
		List<Map<String, Object>> retVal = new ArrayList<>();

		Map<String, Integer> mapCounter = new HashMap<>();
		for (String feature : queryVal) {
			feature = feature.replaceAll("\"", "");
			if (mapCounter.get(feature) == null) {
				mapCounter.put(feature, 1);
			} else {
				int counter = mapCounter.get(feature);
				counter++;
				mapCounter.put(feature, counter);
			}
		}

		for (Map.Entry<String, Integer> entry : mapCounter.entrySet()) {
			Map<String, Object> map = new HashMap<>();
			map.put("license", entry.getKey());
			map.put("value", entry.getValue());
			retVal.add(map);
		}

		return retVal;
	}

	/**
	 * @param location the location filter
	 * @return a List of DB Versions + Counter
	 */
	public List<Map<String, Object>> getDbVersionsCount(final String location) {
		return currentRepo.getDbVersionsCount(location);
	}

	/**
	 * @param location the location filter
	 * @return a List of Host Types + Counter
	 */
	public List<Map<String, Object>> getHostTypeCount(final String location) {
		return currentRepo.getHostTypeCount(location);
	}

	/**
	 * @param location the location filter
	 * @return a List of OS Types + Counter
	 */
	public List<Map<String, Object>> getOsTypeCount(final String location) {
		Pattern[] patterns = buildOsTypePatterns();

		List<Map<String, Object>> values = currentRepo.getOsTypeCount(location);
		Map<String, Integer> aggregatedValues = new HashMap<>();
		List<Map<String, Object>> modifiedValues = new ArrayList<>();
		for (Map<String, Object> value : values) {
			String groupedName = getGroupedNameOfOsType((String) value.get("sistemi"), patterns);
			Integer count = ((BigInteger) value.get("count")).intValue();

			aggregatedValues.put(groupedName,
					(aggregatedValues.containsKey(groupedName)
							? aggregatedValues.get(groupedName) : 0) + count);
		}

		for (Map.Entry<String, Integer> value : aggregatedValues.entrySet()) {
			Map<String, Object> item = new HashMap<>();
			item.put("sistemi", value.getKey());
			item.put("count", value.getValue());

			modifiedValues.add(item);
		}

		return modifiedValues;
	}

	private Pattern[] buildOsTypePatterns() {
		Pattern[] patterns;
		if (osAggregateRulesGroup == null || osAggregateRulesRegex == null) {
			patterns = new Pattern[0];
		} else if (osAggregateRulesGroup.size() != osAggregateRulesRegex.size()) {
			logger.error("Error! the length of dashboard.aggregate.os.rules.regex "
					+ "and dashboard.aggregate.os.rules.group have not the same");
			patterns = new Pattern[0];
		} else {
			patterns = new Pattern[osAggregateRulesGroup.size()];
			for (int i = 0; i < osAggregateRulesRegex.size(); i++) {
				patterns[i] = Pattern.compile(osAggregateRulesRegex.get(i));
			}
		}

		return patterns;
	}

	/**
	 * @param originalName full nmae
	 * @param patterns prebuild patterns
	 * @return
	 */
	private String getGroupedNameOfOsType(final String originalName, final Pattern[] patterns) {
		for (int i = 0; i < patterns.length; i++) {
			if (patterns[i].matcher(originalName).find()) {
				return osAggregateRulesGroup.get(i);
			}
		}
		return originalName;
	}

	/**
	 * @param location the location filter
	 * @return a JSON Array of Hostnames with unused CPUs
	 */
	public List<Map<String, Object>> getResize(final String location) {
		return currentRepo.getResizeByHosts(location);
	}


	/**
	 * @param location the location filter
	 * @return a JSON Array of Hostnames, DBs and Work
	 */
	public List<Map<String, Object>> getWorkByDbs(final String location) {
		return currentRepo.getWorkByDbs(location);
	}

	/**
	 * @param location the location filter
	 * @return a JSON Array of Enterprise&Standard Licenses + Counter
	 */
	public List<Map<BigDecimal, String>> getLicensesCount(final String location) {
		List<Map<BigDecimal, String>> retVal = currentRepo.getLicensesCount(location);

		if (retVal.isEmpty()) {
			retVal = currentRepo.getLicensesCountNoValues();
		}

		return retVal;

	}

	/**
	 * @param location the location filter
	 * @return a JSON Array of All Licenses + Counter
	 */
	public List<Map<String, Object>> getAllLicensesCount(final String location) {
		return currentRepo.getAllLicensesCount(location);
	}

	/**
	 * @param license the license filter
	 * @return a JSON Array of All hosts
	 */
	@Transactional
	public List<String> getAllHostUsingLicense(final String license) {
		return currentRepo.getAllHostsUsingLicense(license);
	}

	/**
	 * @return the result: ("Compliance"=true) if there are more Licenses than used,
	 * ("Compliance"=false) if there are more used Licenses than charged
	 */
	@Transactional
	public Map<String, Object> getCompliance() {
		List<Boolean> queryRes = currentRepo.getCompliance();
		Map<String, BigInteger> totalLicenses = currentRepo.getTotalLicensesForCompliance();
		Map<String, Object> retVal = new HashMap<>();

		if (queryRes.contains(Boolean.valueOf("false"))) {
			retVal.put("Compliance", false);
		} else {
			retVal.put("Compliance", true);
		}

		retVal.put("Licenses", totalLicenses);

		return retVal;
	}

	/**
	 * Get all addms.
	 * @param env the env filter
	 * @param search the search filter
	 * @return the list of addms
	 */
	@Transactional
	public List<Map<String, Object>> getADDMs(final String env, final String search) {
		return currentRepo.getADDMs(env, search);
	}
	
	/**
	 * Get all segment advisors.
	 * @param env the env filter
	 * @param search the search filter
	 * @return the list of segment advisors
	 */
	@Transactional
	public List<Map<String, Object>> getSegmentAdvisors(final String env, final String search) {
		return currentRepo.getSegmentAdvisors(env, search);
	}

	/**
	 * Get all patch advisors.
	 * @param status the env filter
	 * @param windowTime window time
	 * @return the list of patch advisors
	 */
	@Transactional
	public List<Map<String, Object>> getPatchAdvisors(final String status, final int windowTime) {
		Calendar calendar = Calendar.getInstance();
		calendar.set(Calendar.HOUR_OF_DAY, 0);
		calendar.set(Calendar.MINUTE, 0);
		calendar.set(Calendar.DATE, 1);
		calendar.add(Calendar.MONTH, -windowTime);
		Date time = calendar.getTime();
		return currentRepo.getAllHostPSUStatus(time, status);
	}

	/**
	 * Get the environments.
	 * @return the environments.
	 */
	@Transactional
	public List<String> getEnvironments() {
		return currentRepo.getEnviroments();
	}

	/**
	 * Get the top 15 reclaimable database.
	 * @param location location
	 * @return Top 15 reclaimable databas
	 */
	@Transactional
	public List<Map<String, Object>> getTopReclaimableDatabase(final String location) {
		return currentRepo.getTopReclaimableDatabase(location);
	}

	/**
	 * Get patch status stats.
	 * @param location location
	 * @param windowTime windowTime
	 * @return patch status stats
	 */
	@Transactional
	public List<Map<String, Object>> getPatchStatusStats(final String location, final int windowTime)  {
		Calendar calendar = Calendar.getInstance();
		calendar.set(Calendar.HOUR_OF_DAY, 0);
		calendar.set(Calendar.MINUTE, 0);
		calendar.set(Calendar.DATE, 1);
		calendar.add(Calendar.MONTH, -windowTime);
		Date time = calendar.getTime();
		return currentRepo.getPatchStatusStats(location, time);
	}

}
