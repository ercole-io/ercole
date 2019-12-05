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

package io.ercole.utilities;

import java.util.ArrayList;
import java.util.Arrays;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

import org.json.JSONArray;
import org.json.JSONObject;

import io.ercole.model.ClusterInfo;
import io.ercole.model.CurrentHost;
import io.ercole.model.VMInfo;

/**
 * Utility for filtering JSON objects and arrays.
 */
public final class JsonFilter {
	
	private static final String FEATURES = "Features";
	private static final String STATUS = "Status";
	private static final String NAME = "Name";
	
	private JsonFilter() {
		
	}
	
	/**
	 * @param array database JSONArray to manipulate
	 * @return Map with the following structure: 
	 * <JSONArrayObjectAttribute <key,value>>
	 */
	public static Map<String, Map<String, Boolean>> getFeaturesMapping(final JSONArray array) {
		Map<String, Map<String, Boolean>> retMap = new HashMap<>();
		
		for (int x = 0; x < array.length(); x++) {
			JSONObject database = array.getJSONObject(x);
			JSONArray features = database.getJSONArray(FEATURES);
			Map<String, Boolean> elencoFeatures = new HashMap<>();
			
			for (int y = 0; y < features.length(); y++) {
				String key = (String) features.getJSONObject(y).get(NAME);
				Boolean value = (Boolean) features.getJSONObject(y).get(STATUS);
				
				elencoFeatures.put(key, value);
			}
			
			retMap.put(database.getString("Name"), elencoFeatures);
		}
		
		return retMap;
	}
	
	
	/**
	 * @param dbArray JSONArray of JSONObjects that have
	 * a JSONArray as property, called "Features"; this method has to be used only on
	 * Database JSONArrays
	 * @return a Map with key="name of feature" and value=True
	 */
	public static Map<String, Boolean> getTrueFeaturesFromDbArray(final JSONArray dbArray) {
		Map<String, Boolean> retMap = new HashMap<>();
		
		for (int x = 0; x < dbArray.length(); x++) {
			JSONObject database = dbArray.getJSONObject(x);
			JSONArray oldFeatures = database.getJSONArray(FEATURES);
			
			for (int y = 0; y < oldFeatures.length(); y++) {
				String key = (String) oldFeatures.getJSONObject(y).get(NAME);
				Boolean value = (Boolean) oldFeatures.getJSONObject(y).get(STATUS);
				
				if (value) {
					retMap.put(key, value);
				}
			}
		}
		
		return retMap;
	}
	
	
	/**
	 * @param oldDbArray JSONArray of JSONObjects that have
	 * a JSONArray as property, called "Features"; this method has to be used only on
	 * Database JSONArrays
	 * @return a Map with key="name of feature" and value=False
	 */
	public static Map<String, Boolean> getFalseFeaturesFromDbArray(final JSONArray oldDbArray) {
		Map<String, Boolean> retMap = setAllFeaturesToFalse(oldDbArray);
		JSONArray oldFeatures = oldDbArray.getJSONObject(0).getJSONArray(FEATURES);
		
		for (int x = 0; x < oldFeatures.length(); x++) {
			for (int y = 0; y < oldDbArray.length(); y++) {

				String key = (String) oldDbArray.getJSONObject(y).getJSONArray(FEATURES)
						.getJSONObject(x).get(NAME);
				Boolean value = (Boolean) oldDbArray.getJSONObject(y).getJSONArray(FEATURES)
						.getJSONObject(x).get(STATUS);
				
				if (value) {
					retMap.remove(key);
				}
			}
		}
		return retMap;
	}
	
	
	/**
	 * @param dbArray is the JSONArray
	 * @return a Map of Fetures set to false
	 */
	public static Map<String, Boolean> setAllFeaturesToFalse(final JSONArray dbArray) {
		Map<String, Boolean> retMap = new HashMap<>();
		
		JSONObject database = dbArray.getJSONObject(0);
		JSONArray oldFeatures = database.getJSONArray(FEATURES);
		
		for (int x = 0; x < oldFeatures.length(); x++) {	
			String key = (String) oldFeatures.getJSONObject(x).get(NAME);
			retMap.put(key, false);
		}
		
		return retMap;
	}
	
	
	/**
	 * @param object the raw JSONObject incoming from the agent for updates
	 * @return CurrentHost object builded from JSON properties
	 */
	public static CurrentHost buildCurrentHostFromJSON(final JSONObject object) {
		JSONObject hostObj = object.getJSONObject("Info");
		JSONObject extraObj = object.getJSONObject("Extra");
		
		CurrentHost host = new CurrentHost();

		host.setDatabases(object.getString("Databases"));
		host.setHostname(object.getString("Hostname"));
		host.setEnvironment(object.getString("Environment"));
		host.setLocation(object.getString("Location"));
		if (object.has("Version")) {
			host.setVersion(object.getString("Version"));
		} else {
			host.setVersion("unknown");
		}
		host.setServerVersion("<ERROR>");

		if (object.has("HostType")) {
			host.setHostType(object.getString("HostType"));
		} else {
			host.setHostType("oracledb");
		}
		host.setSchemas(object.getString("Schemas"));
		host.setHostInfo(hostObj.toString());
		host.setExtraInfo(extraObj.toString());

		return host;
	}
	
	
	/**
	 * @param newHost the incoming Host from agent update
	 * @param oldHost the old registered Host that's being subsituted by
	 * newHost
	 * @return List of newly added databases
	 */
	public static List<String> getNewDatabases(final CurrentHost newHost, final CurrentHost oldHost) {
		final String[] valuesOld = oldHost.getDatabases().trim().split(" ");
		final String[] valuesNew = newHost.getDatabases().trim().split(" ");

		List<String> retVal = new ArrayList<>();

		if (valuesNew.length > valuesOld.length) {
			retVal = filterMoreNewDBLessOldDB(valuesNew, valuesOld);
		} else {
			List<String> newList = new ArrayList<>();
			
			for (String value : valuesNew) {
				newList.add(value);
			}
			
			for (String value : valuesOld) {
				if (newList.contains(value)) {
					newList.remove(value);
				}
			}

			if (!newList.isEmpty() && !newList.contains("")) {
				retVal = newList;
			}
		}
		return retVal;
	}
	
	
	/**
	 * @param newHost is the New CurrentHost from agent update
	 * @return  a List of new databases if any
	 */
	public static List<String> getDatabases(final CurrentHost newHost) {
		final String[] valuesNew = newHost.getDatabases().trim().split(" ");
		
		return Arrays.asList(valuesNew);
	}
	
	
	// used when incoming more NEW databases than OLD ones
	private static List<String> filterMoreNewDBLessOldDB(final String[] valuesNew, final String[] valuesOld) {
		List<String> retVal = new ArrayList<>();
		
		for (int i = 0; i < valuesNew.length; i++) {
			int counter = 0;
			
			for (int j = 0; j < valuesOld.length; j++) {
				if (valuesNew[i].equals(valuesOld[j])) {
					break;
				}
				
				counter++;
				if (counter == valuesOld.length) {
					retVal.add(valuesNew[i]);
				}
			}
		}
		
		return retVal;
	}
	
	
	/**
	 * @param newHost the incoming Host from agent update
	 * @param oldHost the old registered Host that's being subsituted by
	 * newHost
	 * @return List of same databases between newHost & oldHost
	 */
	public static List<String> getSameDatabases(final CurrentHost newHost, final CurrentHost oldHost) {
		final String[] valuesOld = oldHost.getDatabases().trim().split(" ");
		final String[] valuesNew = newHost.getDatabases().trim().split(" ");
		
		List<String> oldList = Arrays.asList(valuesOld);
		List<String> newList = Arrays.asList(valuesNew);
		
		List<String> retVal = new ArrayList<>();

		if (oldList.isEmpty() || newList.isEmpty()) {
			return retVal;
		} else {
			int i = 0;
			while (i < oldList.size()) {
				if (newList.contains(oldList.get(i))) {
					retVal.add(oldList.get(i));
				}
				i++;
			}
			return retVal;
		}
	}
	
	
	
	/**
	 * @param oldDbs is a JSONArray from old CurrentHost
	 * 		to search for Enterprise Licenses
	 * @param newDbs is a JSONArray from new CurrentHost 
	 * 		to search for Enterprise Licenses
	 * @return true if newDbs has Oracle ENT or Oracle EXT licenses
	 * 		and oldDbs has none of them
	 */
	public static boolean hasNewEnterpriseLicenses(final JSONArray oldDbs, final JSONArray newDbs) {
		return !hasEnterpriseLicenses(oldDbs) && hasEnterpriseLicenses(newDbs);
	}
	
	/**
	 * @param dbArray is the array of DBs incoming from the agent update
	 * @return true if the update has enterprise licenses
	 */
	public static boolean hasEnterpriseLicenses(final JSONArray dbArray) {
		for (int i = 0; i < dbArray.length(); i++) {
			JSONArray licenses = dbArray.getJSONObject(i).getJSONArray("Licenses");
			
			for (int j = 0; j < licenses.length(); j++) {
				if (((String) licenses.getJSONObject(j).get("Name")).matches("[oO]racle ENT") 
					|| ((String) licenses.getJSONObject(j).get("Name")).matches("[oO]racle EXT")
					&& (((Integer) licenses.getJSONObject(j)
							.get("Count")) != 0)) {				
					return true;
				}
			}
		}
		
		return false;
	}
	
	
	
	/**
	 * @param oldHost is the old CurrentHost
	 * @param newHost is the new CurrentHost from agent update
	 * @return true if the new CurrentHost has more CPUs than before
	 */
	public static boolean hasMoreCPUCores(final CurrentHost oldHost, final CurrentHost newHost) {
		Integer oldCPUs = (Integer) new JSONObject(oldHost.getHostInfo()).get("CPUCores");
		Integer newCPUs = (Integer) new JSONObject(newHost.getHostInfo()).get("CPUCores");
		
		return newCPUs > oldCPUs;
	}

	/**
	 * Build a list of cluster infos from Json.
	 * @param array source json
	 * @return the list of cluster infos
	 */
	public static List<ClusterInfo> buildClusterInfosFromJson(final JSONArray array) {
		List<ClusterInfo> infos = new ArrayList<>();
		for (int i = 0; i < array.length(); i++) {
			infos.add(buildClusterInfoFromJsonObject(array.getJSONObject(i)));
		}
		
		return infos;
	}
	/**
	 * Build a cluster from JSONObject.
	 * @param obj source obj
	 * @return a cluster info
	 */
	public static ClusterInfo buildClusterInfoFromJsonObject(final JSONObject obj) {
		ClusterInfo ci = new ClusterInfo();
		ci.setName(obj.getString("Name"));
		ci.setType(obj.getString("Type"));
		ci.setCpu(obj.getInt("CPU"));
		ci.setSockets(obj.getInt("Sockets"));
		ci.setVms(buildVMInfosFromJsonArray(obj.getJSONArray("VMs")));
		return ci;
	}
	/**
	 * Build the list of vm infos from JSONArray.
	 * @param vmsJson the source JSONArray
	 * @return the list of vm infos
	 */
	public static List<VMInfo> buildVMInfosFromJsonArray(final JSONArray vmsJson) {
		List<VMInfo> vms = new ArrayList<>();
		for (int i = 0; i < vmsJson.length(); i++) {
			vms.add(buildVMInfoFromJsonObject(vmsJson.getJSONObject(i)));
		}

		return vms;
	}
	/**
	 * Build the vm info from JSONObject.
	 * @param obj the source JSONObject
	 * @return the vm info 
	 */
	public static VMInfo buildVMInfoFromJsonObject(final JSONObject obj) {
		VMInfo info = new VMInfo();
		info.setClusterName(obj.getString("ClusterName"));
		info.setHostName(obj.getString("Hostname"));
		info.setName(obj.getString("Name"));
		info.setPhysicalHost(obj.getString("PhysicalHost"));
		info.setName(obj.getString("Name"));
		
		if (info.getHostName() == null || info.getHostName().equals("")) {
			info.setHostName(info.getName());
		}
		return info;
	}

	/**
	 * Gets features.
	 *
	 * @param feat the features
	 * @return the features
	 */
	public static String getTrueFeatures(final JSONArray feat) {
		String features = "";
		for (int i = 0; i < feat.length(); i++) {
			if (feat.get(i).toString().contains("true")) {
				JSONObject jsonObject = (JSONObject) feat.get(i);
				features = features.concat(jsonObject.getString("Name") + ",");
			}
		}
		return features;
	}
}
