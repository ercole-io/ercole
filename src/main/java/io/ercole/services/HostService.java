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

import java.text.ParseException;
import java.time.ZonedDateTime;
import java.util.Date;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import javax.transaction.Transactional;

import org.apache.commons.lang3.time.DateUtils;
import org.json.JSONArray;
import org.json.JSONObject;
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Service;

import io.ercole.model.Alert;
import io.ercole.model.AlertFactory;
import io.ercole.model.ClusterInfo;
import io.ercole.model.CurrentHost;
import io.ercole.model.HistoricalHost;
import io.ercole.model.VMInfo;
import io.ercole.repositories.AlertRepository;
import io.ercole.repositories.ClusterRepository;
import io.ercole.repositories.CurrentHostRepository;
import io.ercole.repositories.HistoricalHostRepository;
import io.ercole.utilities.DateUtility;
import io.ercole.utilities.JsonFilter;

/**
 * Service component for Host.
 */
@Service
public class HostService {

	private Logger logger = LoggerFactory.getLogger(HostService.class);
	
	private int offset = ZonedDateTime.now().getOffset().getTotalSeconds();
	private static final String UPDATED = "updated";
	private static final String SKIPPED = "skipped";
	private static final String INSERTED = "inserted";
	private static final String ERROR = "error";
	private static final String DATABASES = "Databases";
	private static final String HOSTNAME = "Hostname";
	private static String orclDb = "oracledb";

	@Value("${agent.update.rate}")
	private int updateRate;

	@Value("${application.version}")
	private String version;

	@Autowired
	private CurrentHostRepository currentRepo;

	@Autowired
	private HistoricalHostRepository historicalRepo;

	@Autowired
	private AlertRepository alertRepo;
	
	@Autowired
	private ClusterRepository clusterRepo;
	
	@Autowired
	private MailService mailService;

	/**
	 * @param object
	 *            JSON object to deserialize
	 * @param hostType The hostType
	 * @return String of custom status code
	 * @throws ParseException if there are problems with feature description
	 * 		in all DBs
	 */
	@Transactional
	public String updateWithAgent(final JSONObject object, final String hostType) throws ParseException {
		
		CurrentHost host = JsonFilter.buildCurrentHostFromJSON(object);
		fixAssociatedClusterName(host);
		//Retrocompatibility check!
		if (host.getHostType() == null) {
			host.setHostType(orclDb);
			object.put("HostType", orclDb);
		}

		String hostname = object.getString(HOSTNAME);
		CurrentHost oldCurrent = currentRepo.findByHostname(hostname);
		if (oldCurrent != null && (DateUtility.isValidUpdateRange(oldCurrent.getUpdated(), 
				updateRate))) {		
			
			//Compatibility check!
			if (oldCurrent.getHostType() == null) {
				oldCurrent.setHostType(hostType);
			}
			processUpdate(object, oldCurrent);
			
			return UPDATED;
		} else if (oldCurrent == null) {
			processInsert(object);
			
			return INSERTED;
		} else {
			return ERROR;
		}
	}
	
	private boolean processUpdate(final JSONObject object, final CurrentHost oldCurrent) throws ParseException {
		String hostname = object.getString(HOSTNAME);
		CurrentHost newHost = JsonFilter.buildCurrentHostFromJSON(object);
		newHost.setServerVersion(version);
		fixAssociatedClusterName(newHost);
		moveFromCurrentToHistorical(newHost, oldCurrent);
		AlertFactory generator = new AlertFactory();

		assert (newHost.getHostType() != null);
		if (newHost.getHostType().equals(orclDb)) {
			List<String> newDatabases = JsonFilter.getNewDatabases(newHost, oldCurrent);

			if (!newDatabases.isEmpty()) {
				Alert alert = alertRepo.save(generator.fireNewDatabaseAlert(newDatabases, hostname));
				mailService.send(alert);
			}	

			JSONArray oldDbArray = new JSONObject(oldCurrent.getExtraInfo()).getJSONArray(DATABASES);
			JSONArray newDbArray = new JSONObject(newHost.getExtraInfo()).getJSONArray(DATABASES);
	
			if (JsonFilter.hasMoreCPUCores(oldCurrent, newHost) 
					|| JsonFilter.hasNewEnterpriseLicenses(oldDbArray, newDbArray)) {
				Alert alert = alertRepo.save(generator.getAlertForEnterpriseLicenseActivated(hostname));
				mailService.send(alert);
			}

			Map<String, Map<String, Boolean>> elencoNewDbs = JsonFilter.getFeaturesMapping(newDbArray);
			Map<String, Map<String, Boolean>> elencoOldDbs = JsonFilter.getFeaturesMapping(oldDbArray);
	
			Map<String, Boolean> elencoOldActivatedFeatures;
			Map<String, Boolean> elencoOldDeactivatedFeatures;
	
			if (!elencoNewDbs.isEmpty()) {
				elencoOldActivatedFeatures = JsonFilter.getTrueFeaturesFromDbArray(oldDbArray);
	
				for (Alert alert : generator.getAlertForDuplicatedActiveFeature(elencoNewDbs,
						elencoOldActivatedFeatures, elencoOldDbs, hostname)) {
					Alert a = alertRepo.save(alert);
					mailService.send(a);
				}
				
				if (oldDbArray.length() != 0) {
					elencoOldDeactivatedFeatures =
							JsonFilter.getFalseFeaturesFromDbArray(oldDbArray);
				} else {
					elencoOldDeactivatedFeatures =
							JsonFilter.setAllFeaturesToFalse(newDbArray);
				}
				
				List<Alert> featureAlerts;
				
				featureAlerts = generator.getAlertforNewActiveFeature(elencoNewDbs, 
							elencoOldDeactivatedFeatures, hostname);
				if (featureAlerts != null) {
					for (Alert alert : featureAlerts) {
						Alert a = alertRepo.save(alert);
						mailService.send(a);
					}
				}
			}
		} else if (newHost.getHostType().equals("virtualization")) {
			JSONObject oldClustersJSON = new JSONObject(oldCurrent.getExtraInfo());
			JSONObject newClustersJSON = new JSONObject(newHost.getExtraInfo());
			List<ClusterInfo> oldClusters = JsonFilter.buildClusterInfosFromJson(oldClustersJSON.getJSONArray("Clusters"));
			List<ClusterInfo> newClusters = JsonFilter.buildClusterInfosFromJson(newClustersJSON.getJSONArray("Clusters"));

			//Clean old clusters
			oldClusters.forEach(cl -> {
				//Search a new cluster with the same name
				for (ClusterInfo newCl : newClusters) {
					if (newCl.getName() == cl.getName()) {
						return;
					}
				}
				clusterRepo.delete(cl);
			});
			updateOrInsertClustersInfo(
					JsonFilter.buildClusterInfosFromJson(newClustersJSON.getJSONArray("Clusters")));
			fixAssociatedClusterName(newHost);
		}
		
		return true;
	}

	private void moveFromCurrentToHistorical(final CurrentHost currentToInsert, 
			final CurrentHost currentToHistorize) {
		archiveHost(currentToHistorize);
		currentRepo.delete(currentToHistorize);
		currentToInsert.setUpdated(new Date());
		currentRepo.save(currentToInsert);
	}

	private void archiveHost(final CurrentHost current) {
		HistoricalHost historical = new HistoricalHost(current.getId(), 
				current.getHostname(), current.getEnvironment(),
				current.getLocation(), current.getVersion(), 
				current.getServerVersion(), current.getHostType(), 
				current.getDatabases(), current.getSchemas(), 
				current.getExtraInfo(), current.getAssociatedClusterName(), current.getAssociatedHypervisorHostname(), current.getHostInfo(), 
				current.getUpdated());
		historicalRepo.save(historical);
	}
	
	
	
	private boolean processInsert(final JSONObject object) {
		String hostname = object.getString(HOSTNAME);
		CurrentHost host = JsonFilter.buildCurrentHostFromJSON(object);
		host.setServerVersion(version);
		fixAssociatedClusterName(host);

		//Security check!
		assert (host.getHostType() != null);
		
		host.setUpdated(new Date());
		currentRepo.save(host);
		
		AlertFactory generator = new AlertFactory();

		if (host.getHostType().equals(orclDb)) {
			JSONArray newDbArray = new JSONObject(host.getExtraInfo()).getJSONArray(DATABASES);
		

			if (JsonFilter.hasEnterpriseLicenses(newDbArray)) {
				Alert alert = alertRepo.save(generator.getAlertForEnterpriseLicenseActivated(hostname));
				mailService.send(alert);
			}
			
			Map<String, Map<String, Boolean>> newFeaturesByDb = JsonFilter.getFeaturesMapping(newDbArray);
	
			if (!newFeaturesByDb.isEmpty()) {
				for (Alert alert : generator.getAlertforNewActiveFeature(newFeaturesByDb, 
						host.getHostname())) {
					Alert a = alertRepo.save(alert);
					mailService.send(a);
				}
			}
	
			List<String> dbs = JsonFilter.getDatabases(host);
	
			if (!dbs.isEmpty()) {
				Alert alert = alertRepo.save(generator.fireNewDatabaseAlert(dbs, host.getHostname()));
				mailService.send(alert);
			}	
		} else if (host.getHostType().equals("virtualization")) {
			JSONObject clustersJSON = new JSONObject(host.getExtraInfo());
			updateOrInsertClustersInfo(
					JsonFilter.buildClusterInfosFromJson(clustersJSON.getJSONArray("Clusters")));
			fixAssociatedClusterName(host);
		}
		

		Alert alert = alertRepo.save(generator.fireNewServerAlert(hostname));
		mailService.send(alert);

		return true;
	}
	

	/**
	 * @param hostname
	 *            to search in HistoricalHost entities
	 * @param date
	 *            date queried
	 * @return HistoricalHost object found
	 */
	public HistoricalHost getHistoricalLogs(final String hostname, final Date date) {
		HistoricalHost historical = null;

		if (!historicalRepo.findFirstHostnameByArchivedDesc(hostname, date).isEmpty()) {
			historical = historicalRepo.findFirstHostnameByArchivedDesc(hostname, date).get(0);
			historical.setUpdated(DateUtils.addSeconds(historical.getUpdated(), offset));
			historical.setArchived(DateUtils.addSeconds(historical.getArchived(), offset));
		}

		return historical;
	}

	/**
	 * @param updateRate
	 *            used by tests
	 */
	public void setUpdateRate(final int updateRate) {
		this.updateRate = updateRate;
	}

	/**
	 * Update or insert all information about the clusters and update all clusters host.
	 * @param infos informations
	 */
	public void updateOrInsertClustersInfo(final List<ClusterInfo> infos) {
		//Foreach clusters: add cluster (and vms), and update all vms data
		infos.forEach(cluster -> {
			//Find old cluster with the same name
			ClusterInfo oldCluster = clusterRepo.findByName(cluster.getName());
			//Update the associated cluster name of the vms of the old cluster 
			if (oldCluster != null) {
				oldCluster.getVms().forEach(vm -> {
					CurrentHost foundHost = currentRepo.findByHostname(vm.getHostName());
					if (foundHost != null) {
						foundHost.setAssociatedClusterName(null);
						foundHost.setAssociatedHypervisorHostname(null);
						currentRepo.save(foundHost);
					}
				});
				clusterRepo.delete(oldCluster);	
			}
			//Save the current cluster
			cluster.setUpdated(new Date());
			cluster = clusterRepo.save(cluster);
			//Update all relative VM
			cluster.getVms().forEach(vm -> {
				CurrentHost foundHost = currentRepo.findByHostname(vm.getHostName());
				if (foundHost != null) {
					foundHost.setAssociatedClusterName(vm.getClusterName());
					foundHost.setAssociatedHypervisorHostname(vm.getPhysicalHost());
					currentRepo.save(foundHost);
				}
			});
		});
	}

	/**
	 * Fix the associated cluster name.
	 * @param current current
	 */
	public void fixAssociatedClusterName(final CurrentHost current) {
		VMInfo info = clusterRepo.findOneVMInfoByHostname(current.getHostname());
		if (info != null) {
			current.setAssociatedClusterName(info.getClusterName());
			current.setAssociatedHypervisorHostname(info.getPhysicalHost());
		}
	}

	/**
	 * Archive the host.
	 * @param hostname hostname of the host
	 */
	public void archiveHost(final String hostname) {
		CurrentHost host = currentRepo.findByHostname(hostname);
		System.out.println(host);
		if (host != null) {
			archiveHost(host);
			currentRepo.delete(host);
		}
	}

    /**
     * Return all cluster that match the filter.
     * @param filter filter
     * @return all cluster that match the filter
     */	
	public List<ClusterInfo> getClusters(final String filter) {
		return clusterRepo.getClusters(filter);
	}

	/**
	 * Return the used data history of all databases of host.
	 * @param hostname hostname
	 * @return the used data history of all databases of host
	 */
	public Map<String, Object> getUsedDataHistory(final String hostname) {
		CurrentHost current = currentRepo.findByHostname(hostname);
		HashMap<String, Object> result = new HashMap<>();
		if (current == null) {
			return result;
		}
		String[] dbs = current.getDatabases().split(" "); 
		for (String db : dbs) {
			result.put(db, currentRepo.getUsedDataHistory(hostname, db));
		}
		return result;
	}

	/**
	 * Return the segmentsSize data history of all databases of host.
	 * @param hostname hostname
	 * @return the segmentsSize data history of all databases of host
	 */
	public Map<String, Object> getSegmentsSizeDataHistory(final String hostname) {
		CurrentHost current = currentRepo.findByHostname(hostname);
		HashMap<String, Object> result = new HashMap<>();
		if (current == null) {
			return result;
		}
		String[] dbs = current.getDatabases().split(" "); 
		for (String db : dbs) {
			result.put(db, currentRepo.getSegmentsSizeDataHistory(hostname, db));
		}
		return result;
	}
}
