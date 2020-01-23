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

import java.text.ParseException;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;

/**
 * Business rules for generating Alerts.
 */
public class AlertFactory {
	
	private static final String DATABASE = "Database";
	private static final String SERVER = "server";
	private static final String HAS_ENABLED = " has enabled ";
	
	/**
	 * Creates a basic generator.
	 */
	public AlertFactory() {
		// used for generating Alert objects based on business implementations
	}
	
	
	
	/**
	 * @param newDatabases from the incoming agent update
	 * @param hostname where newDatabases are hosted
	 * @return Alert, with Severy NOTICE, Code NEW_DATABASE,
	 * server hostname for "having created a new DB on the server"
	 */
	public Alert fireNewDatabaseAlert(final List<String> newDatabases, final String hostname) {
		Alert alert;
		if (!newDatabases.isEmpty() && newDatabases.size() == 1) {
			alert = new Alert(hostname, AlertCode.NEW_DATABASE, 
					"Server: " + hostname 
					+ " new database created: " + newDatabases.get(0), AlertSeverity.NOTICE);
		} else {
			alert = new Alert(hostname, AlertCode.NEW_DATABASE, "Server " + hostname 
					+ " new databases created: " + newDatabases.toString(), 
					AlertSeverity.NOTICE);
		}
		return alert;
	}
	
	
	/**
	 * @param hostname is the Hostname
	 * @return an Alert for having added a new Server
	 */
	public Alert fireNewServerAlert(final String hostname) {
		String description = "Server " + hostname + " added to Ercole";
		return new Alert(hostname, AlertCode.NEW_SERVER, description, 
				AlertSeverity.NOTICE);
	}
	
	/**
	 * @param hostname is the Hostname
	 * @return an Alert for having added a new Server
	 */
	public Alert fireMissingHostAlert(final String hostname) {
		String description = "Server " + hostname + " is missing";
		return new Alert(hostname, AlertCode.MISSING_HOST, description, 
				AlertSeverity.NOTICE);
	}

	/**
	 * @param newDbs map of Databases incoming from agent
	 * @param oldActiveFeatures map of allready activated Features on server
	 * @param oldDbs map of Databases from old log
	 * @param hostname to use for logging the new Alert
	 * @return Alert, with Severy NOTICE, Code NEW_OPTION,
	 * server hostname for "having activated a feature on a DB,
	 * but the same feature is allready activated in the same
	 * hostname"
	 */
	public List<Alert> getAlertForDuplicatedActiveFeature(final Map<String, Map<String, Boolean>> newDbs,
		final Map<String, Boolean> oldActiveFeatures, 
		final Map<String, Map<String, Boolean>> oldDbs, final String hostname) {
		

		List<Alert> retVal = new ArrayList<>();
		
		for (Map.Entry<String, Map<String, Boolean>> entry : newDbs.entrySet()) {
			Map<String, Boolean> newFeatures = entry.getValue();
			List<String> featuresAttivate = new ArrayList<>();
			StringBuilder builder = new StringBuilder();
			
			for (String feature : oldActiveFeatures.keySet()) {
				
				// se l'update contiene lo stesso db, MA con features precedentemente 
				// disattivate e ora attivate, e se altri db che avevano già la feature
				// attiva
				if (oldDbs.get(entry.getKey()) != null 
						&& oldDbs.get(entry.getKey()).get(feature) == null
						&& newFeatures.get(feature) != null) {
					featuresAttivate.add(feature);	
				}
				
				// se invece si tratta di un nuovo db
				if (oldDbs.get(entry.getKey()) == null 
						&& newFeatures.get(feature) != null) {
					featuresAttivate.add(feature);
				}
				
				
			}
			
			if (!featuresAttivate.isEmpty()) {
				builder.append(DATABASE + " " + entry.getKey() + HAS_ENABLED
						+ " these options (already present on other databases) "
						+ SERVER + ": " + featuresAttivate);
				retVal.add(new Alert(hostname, AlertCode.NEW_OPTION, 
						builder.toString(), AlertSeverity.NOTICE));
			} 
		}
		
		return retVal;
	}
	
	
	/**
	 * @param newDbs map of Databases incoming from agent
	 * @param oldDeactivatedFeatures map of deactivated Features on server
	 * @param hostname to use for logging the new Alert
	 * @return Alert, with Severy CRITICAL, Code NEW_OPTION,
	 * server hostname for "having activated a feature on a DB,
	 * that hasn't been activated before on the hostname"
	 * @throws ParseException if there are problems in 
	 * 		Feature description among all DBs: one could have more/less chars
	 * 		than the first one from DB array
	 */
	public List<Alert> getAlertforNewActiveFeature(final Map<String, Map<String, Boolean>> newDbs,
			final Map<String, Boolean> oldDeactivatedFeatures, 
			final String hostname) throws ParseException {
	
		List<Alert> retVal = new ArrayList<>();
		
		for (Map.Entry<String, Map<String, Boolean>> entry : newDbs.entrySet()) {
			StringBuilder builder = new StringBuilder();

			Map<String, Boolean> newFeatures = entry.getValue();
			String db = entry.getKey();
			List<String> featuresAttivate = new ArrayList<>();

			for (String feature : oldDeactivatedFeatures.keySet()) {		
				
				if (newFeatures.get(feature) == null) {
					throw new ParseException(feature + " on " + hostname, 0);
				}
				// se la feature disattivata si trova ATTIVA tra le nuove features
				// aggiungi alla lista dei risultati per il db
				if (newFeatures.get(feature) != null && newFeatures.get(feature)) {
					featuresAttivate.add(feature);	
				}				
			}
			
			if (!featuresAttivate.isEmpty()) {
				builder.append(DATABASE + " " + db + HAS_ENABLED
						+ "new options on server"
						+ SERVER + ": " + featuresAttivate);
				 retVal.add(new Alert(hostname, AlertCode.NEW_OPTION, 
						builder.toString(), AlertSeverity.CRITICAL));
			}
		}
		return retVal;
	}



	/**
	 * @param elencoNewDbs is a map of <DB, Map<Feature,True/False>>
	 * @param hostname is the hostname on which the DBs are present
	 * @return a List of Alerts of new Features
	 */
	public List<Alert> getAlertforNewActiveFeature(final Map<String, Map<String, Boolean>> elencoNewDbs, 
			final String hostname) {
		
		List<Alert> retVal = new ArrayList<>();
		StringBuilder builder;
		
		for (Map.Entry<String, Map<String, Boolean>> entry : elencoNewDbs.entrySet()) {
			List<String> featuresAttivate = new ArrayList<>();
			Map<String, Boolean> newFeatures = entry.getValue();
			String db = entry.getKey();
			builder = new StringBuilder();
			
			
			for (Map.Entry<String, Boolean> entry2 : newFeatures.entrySet()) {
				
				// se la feature disattivata si trova ATTIVA tra le nuove features
				// aggiungi alla lista dei risultati per il db
				if (newFeatures.get(entry2.getKey())) {
					featuresAttivate.add(entry2.getKey());	
				}	
			}
			
			if (!featuresAttivate.isEmpty()) {
				builder.append(DATABASE + " " + db + HAS_ENABLED
						+ "new options on server "
						+ SERVER + ": " + featuresAttivate);
				 retVal.add(new Alert(hostname, AlertCode.NEW_OPTION, 
						builder.toString(), AlertSeverity.CRITICAL));
			}			
			
		}
		
		return retVal;
	}
	
	/**
	 * @param hostname to get the Alert for
	 * @return Alert for having activated Enterprise License
	 * 		(CPU Cores have grown or Enterprise License activated)
	 */
	public Alert getAlertForEnterpriseLicenseActivated(final String hostname) {
		StringBuilder builder = new StringBuilder();
		builder.append("A new Enterprise license has been enabled");
		return new Alert(hostname, AlertCode.NEW_LICENSE, builder.toString(), AlertSeverity.CRITICAL);
	}
}
