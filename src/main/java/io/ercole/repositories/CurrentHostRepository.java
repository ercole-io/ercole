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

package io.ercole.repositories;

import java.math.BigDecimal;
import java.math.BigInteger;
import java.util.Date;
import java.util.List;
import java.util.Map;
import java.util.stream.Stream;

import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.PagingAndSortingRepository;
import org.springframework.data.repository.query.Param;

import io.ercole.model.CurrentHost;

/**
 * The repository for CurrentHosts.
 */
public interface CurrentHostRepository extends PagingAndSortingRepository<CurrentHost, Long> {

	/**
	 * Find all by order by hostname asc iterable.
	 *
	 * @return the iterable
	 */
	Iterable<CurrentHost>  findAllByOrderByHostnameAsc();
	
	/**
	 * Saves in database the current host.
	 * @param host to save
	 * @return CurrentHost saved object
	 **/
	@SuppressWarnings("unchecked")
	CurrentHost save(CurrentHost host);

	/**
	 * Find by hostname current host.
	 *
	 * @param hostname to search
	 * @return CurrentHost object
	 */
	CurrentHost findByHostnameIgnoreCase(@Param("hostname")String hostname);

	/**
	 * Find by hostname current host.
	 *
	 * @param hostname to search
	 * @return CurrentHost object
	 */
	CurrentHost findByHostname(@Param("hostname")String hostname);


	/**
	 * Find all hosts stream.
	 *
	 * @return the list of hosts
	 */
	@Query("SELECT m FROM CurrentHost m")
	Stream<CurrentHost> findAllHosts();

	/**
	 * Find by db.
	 *
	 * @param db the db
	 * @param c  the c
	 * @return the page
	 */
	@Query("SELECT m FROM CurrentHost m WHERE  (m.hostType IS NULL OR m.hostType = 'oracledb') "
			+ "AND LOWER(m.databases) LIKE LOWER (CONCAT('%',:db,'%'))")
	Page<CurrentHost> findByDb(@Param("db") String db, Pageable c);

	/**
	 * Find databases by hostname current host.
	 *
	 * @param hostname to search
	 * @return CurrentHost object
	 */
	CurrentHost findDatabasesByHostname(String hostname);

	/**
	 * Find by db or by hostname.
	 *
	 * @param ricerca the ricerca
	 * @param c       the c
	 * @return the page
	 */
	@Query("SELECT m FROM CurrentHost m WHERE ((m.hostType IS NULL OR m.hostType = 'oracledb') "
			+ "AND LOWER(m.databases) LIKE LOWER (CONCAT('%',:ricerca,'%')))"
			+ " OR LOWER(m.hostname) LIKE LOWER (CONCAT('%',:ricerca,'%'))")
	Page<CurrentHost> findByDbOrByHostname(@Param("ricerca") String ricerca, Pageable c);

	/**
	 * Find by schema.
	 *
	 * @param schema the schema
	 * @param c      the c
	 * @return the page
	 */
	@Query("SELECT m FROM CurrentHost m WHERE (m.hostType IS NULL OR m.hostType = 'oracledb') "
			+ "AND LOWER(m.extraInfo) LIKE LOWER (CONCAT('%',:schema,'%'))")
	Page<CurrentHost> findBySchema(@Param("schema") String schema, Pageable c);

	/**
	 * Find by db or by hostname or by schema.
	 *
	 * @param ricerca the ricerca
	 * @param date    the date
	 * @param c       the c
	 * @return the page
	 */
	@Query(nativeQuery = true, value =
			"SELECT m.id, m.databases,m.environment,m.host_info,m.hostname,m.location,m.schemas,m.updated,"
		+ " m.host_type, m.associated_cluster_name, m.associated_hypervisor_hostname, NULL as extra_info, m.version, m.server_version from current_host m WHERE "
		+ " m.updated <= :date AND "
		+ " ("
		+ "   ( (m.host_type IS NULL OR m.host_type = 'oracledb') "
		+ "AND LOWER(m.databases) LIKE LOWER (CONCAT('%',:ricerca,'%')))"
		+ "   OR ( LOWER(m.hostname) LIKE LOWER (CONCAT('%',:ricerca,'%')))"
		+ "   OR ( (m.host_type IS NULL OR m.host_type = 'oracledb') "
		+ "AND (m.schemas) LIKE (CONCAT('%',:ricerca,'%')))"
		+ ")")
	Page<CurrentHost> findByDbOrByHostnameOrBySchema(
			@Param("ricerca") String ricerca,
			@Param("date") Date date, Pageable c);


	/**
	 * Find all not updated list.
	 *
	 * @param date the date
	 * @return List<CurrentHost>  with CurrentHost that have "updated" property lower than date
	 */
	@Query("SELECT m FROM CurrentHost m WHERE m.updated <= :date")
	List<CurrentHost> findAllNotUpdated(@Param("date") Date date);


	/**
	 * Gets server type count.
	 *
	 * @param location the location filter
	 * @return List of maps where key = environment & value = count
	 */
	@Query(nativeQuery = true, value = "with vista AS "
			+ "(select environment, 1 as contatore from current_host "
			+ "where ('*' = :location or location = :location)) "
			+ "select environment as label, count(*) as data from vista "
			+ "group by environment")
	List<Map<String, Object>> getServerTypeCount(@Param("location") String location);


	/**
	 * Gets locations.
	 *
	 * @return different types of active server locations
	 */
	@Query("SELECT DISTINCT m.location from CurrentHost m")
	List<String> getLocations();

	/**
	 * Gets db envs.
	 *
	 * @param location the location filter
	 * @return a List of count of DBs for each Environment
	 */
	@Query(nativeQuery = true, value = "with vista AS (select environment, "
			+ "databases from current_host "
			+ "where ((host_type IS NULL or host_type = 'oracledb') "
			+ "and ('*' = :location or location = :location))), "
			+ "vista2 AS (select regexp_split_to_table(databases, E'\\\\s+'), "
			+ "environment from vista) "
			+ "select count(*), environment from vista2 group by environment")
	List<Map<String, Object>> getDbEnvs(@Param("location") String location);


	/**
	 * Gets db features count.
	 *
	 * @param location the location filter
	 * @return a list of Features with status True from all DBs
	 */
	@Query(nativeQuery = true, value = "With reports as ( "
			+ "select a.p as interno from current_host ch, "
			+ "jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') a (p) "
			+ "where ((host_type IS NULL or host_type = 'oracledb') "
			+ "and ('*' = :location or location = :location))), "
			+ "vista as ( "
			+ "SELECT value "
			+ "FROM reports r, jsonb_array_elements(r.interno#>'{Features}') obj "
			+ "WHERE CAST((obj->>'Status') AS boolean) is true) "
			+ "select CAST(p.v AS text) as value "
			+ "from vista cross join lateral jsonb_each(value) p(k,v) "
			+ "where p.v <> 'true'")
	List<String> getDbFeaturesCount(@Param("location") String location);


	/**
	 * Gets db versions count.
	 *
	 * @param location the location filter
	 * @return a result set of DB versions and counter of occurrences
	 */
	@Query(nativeQuery = true, value = "With reports as ( "
			+ "select a.p as interno from current_host ch, "
			+ "jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') a (p) "
			+ "where ((host_type IS NULL or host_type = 'oracledb')"
			+ " and ('*' = :location or location = :location))) "
			+ "SELECT count(*), regexp_replace(CAST((r.interno#>'{Version}') as text), "
			+ "'\"','','g') as version "
			+ "FROM reports r "
			+ "group by version")
	List<Map<String, Object>> getDbVersionsCount(@Param("location") String location);


	/**
	 * Gets host type count.
	 *
	 * @param location the location filter
	 * @return a result set of host types (physical, virtual..) and counter of occurences
	 */
	@Query(nativeQuery = true, value = "with vista as "
			+ "(SELECT CAST(host_info as json)->'Type' "
			+ "as type, 1 as counter FROM current_host "
			+ "where ('*' = :location or location = :location)) "
			+ "select regexp_replace(CAST(type as text),'\"','','g') as tipo, count(*) "
			+ "from vista "
			+ "group by tipo")
	List<Map<String, Object>> getHostTypeCount(@Param("location") String location);
	
	
	/*NB: ATTENTION!!! the maps returned contain a mixed-type pairs. They are a String-BigInteger/String pair*/

	/**
	 * Gets os type count.
	 *
	 * @param location the location filter
	 * @return a result set of host OS types and counter of occurences
	 */
	@Query(nativeQuery = true, value = "with vista as "
			+ "(SELECT CAST(host_info as json)->'OS' as os, "
			+ "1 as counter FROM current_host "
			+ "where ('*' = :location or location = :location)) "
			+ "select regexp_replace(CAST(os as text),'\"','','g') as sistemi, count(*) "
			+ "from vista "
			+ "group by sistemi")
	List<Map<String, Object>> getOsTypeCount(@Param("location") String location);


	/**
	 * This is '/getresize' endpoint.
	 *
	 * @param location the location filter
	 * @return a result set of hostnames and unused CPUs
	 */
	@Query(nativeQuery = true, value = "With reports as ( "
			+ "select a.p as interno, hh->'CPUThreads' as cores, ch.hostname " 
			+ "from current_host ch, " 
			+ "jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') a (p), " 
			+ "CAST(host_info as jsonb) hh " 
			+ "where ((host_type IS NULL or host_type = 'oracledb')"
			+ " and ('*' = :location or location = :location))), "
			+ "vista as ( " 
			+ "SELECT hostname, CAST(CAST(cores as text) as smallint), " 
			+ "CAST(regexp_replace(CAST((r.interno->'Work') as text),'\"','','g') as smallint) as work, " 
			+ "CAST((r.interno->'Name') as text) as name " 
			+ "FROM reports r " 
			+ "WHERE (r.interno->>'Work') ~ '\\d+') " 
			+ "SELECT hostname, cores - sum(work) as resize " 
			+ "from vista " 
			+ "group by hostname, cores "
			+ "order by resize DESC "
			+ "LIMIT 10")
	List<Map<String, Object>> getResizeByHosts(@Param("location") String location);


	/**
	 * Gets work by dbs.
	 *
	 * @param location the location filter
	 * @return Top Work with DB and Hostname
	 */
	@Query(nativeQuery = true, value = "With reports as ( "
			+ "select a.p as interno, hh->'CPUCores' as cores, ch.hostname " 
			+ "from current_host ch, " 
			+ "jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') a (p), " 
			+ "CAST(host_info as jsonb) hh " 
			+ "where ((host_type IS NULL or host_type = 'oracledb') and "
			+ "('*' = :location or location = :location))) "
			+ "SELECT hostname, " 
			+ "CAST(regexp_replace(CAST((r.interno->'Work') as text),'\"','','g') as smallint) as work, " 
			+ "regexp_replace(CAST((r.interno->'Name') as text),'\"','','g') as database " 
			+ "FROM reports r " 
			+ "WHERE (r.interno->>'Work') ~ '\\d+' " 
			+ "order by work DESC "
			+ "LIMIT 10")
	List<Map<String, Object>> getWorkByDbs(@Param("location") String location);


	/**
	 * Gets licenses count.
	 *
	 * @param location the location filter
	 * @return a result set of Enterprise & Standard
	 */
	@Query(nativeQuery = true, value = "WITH host_info AS (SELECT ch.hostname, "
		+	"ch.associated_cluster_name IS NOT NULL AS virtual, "
		+ "dbs->'Licenses' AS licenses, ch.associated_cluster_name AS cluster_name FROM current_host ch, "
		+ "jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS dbs WHERE (ch.host_type "
		+ "IS NULL OR ch.host_type = 'oracledb') AND ('*' = :location or ch.location = :location)), " 
		+ "featured_host_info AS (SELECT phi.hostname, phi.virtual, phi.cluster_name, lic->'Name' AS "
		+ "license_name, CAST(CAST(lic->'Count' AS TEXT) AS REAL) AS license_count FROM host_info phi, "
		+ "jsonb_array_elements(phi.licenses) AS lic), aggregated_featured_host_info AS (SELECT " 
		+ "fhi.hostname, fhi.virtual, fhi.cluster_name, fhi.license_name, max(fhi.license_count) AS " 
		+ "license_count FROM featured_host_info fhi GROUP BY fhi.hostname, "
		+	"fhi.license_name, fhi.cluster_name, "
		+ "fhi.virtual), summed_phy_features AS (SELECT afhi.license_name, sum(afhi.license_count) AS " 
		+ "license_count FROM aggregated_featured_host_info afhi WHERE virtual = false GROUP BY license_name "
		+ "), virtual_featured_host_info AS (SELECT afhi.cluster_name, ci.cpu, afhi.license_name, "
		+ "max(afhi.license_count) AS license_count FROM aggregated_featured_host_info afhi LEFT JOIN "
		+ "cluster_info ci ON ci.name = afhi.cluster_name WHERE "
		+ "afhi.virtual = true GROUP BY afhi.cluster_name, "
		+ "afhi.virtual, afhi.license_name, ci.cpu), summed_virtual_features AS (SELECT vfhi.license_name, "
		+ "sum( CASE WHEN vfhi.license_count > 0 THEN vfhi.cpu*50 ELSE 0 END ) AS license_count FROM "
		+ "virtual_featured_host_info vfhi GROUP BY license_name), partial_summed_featured AS (SELECT * "
		+ "FROM summed_phy_features UNION ALL SELECT * FROM summed_virtual_features), summed_featured AS "
		+ "(SELECT regexp_replace(CAST(license_name AS TEXT),'\"','','g') AS name, sum(license_count) AS sum "
		+ "FROM partial_summed_featured GROUP BY license_name), partial_oracle_licenses_count AS ( SELECT "
		+ "( CASE WHEN name = 'Oracle STD' THEN 'Standard' ELSE 'Enterprise' END) AS type, sf.sum AS counter "
		+ "FROM summed_featured sf WHERE name LIKE 'Oracle%'), oracle_licenses_count AS (SELECT type, "
		+ "sum(counter) as counter FROM partial_oracle_licenses_count polc GROUP BY type) SELECT * FROM " 
		+ "oracle_licenses_count")
	List<Map<BigDecimal, String>> getLicensesCount(@Param("location") String location);


	/**
	 * Gets licenses count no values.
	 *
	 * @return a result set of Enterprise & Standard licenses and 0 counter of occurrences
	 */
	@Query(nativeQuery = true, value = "SELECT * "
			+ "FROM (VALUES (0, 'Enterprise'), (0, 'Standard')) as t(counter, type)")
	List<Map<BigDecimal, String>> getLicensesCountNoValues();


	/**
	 * Gets all licenses count.
	 *
	 * @param location the location filter
	 * @return a ResultSet of all kind of Licenses and relative counter
	 */
	@Query(nativeQuery = true, value = "WITH host_info AS (SELECT ch.hostname, "
		+	"ch.associated_cluster_name IS NOT NULL AS virtual, "
		+ "dbs->'Licenses' AS licenses,  ch.associated_cluster_name AS cluster_name FROM current_host ch, "
		+ "jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS dbs WHERE (ch.host_type "
		+ "IS NULL OR ch.host_type = 'oracledb') AND ('*' = :location or ch.location = :location)), " 
		+ "featured_host_info AS (SELECT phi.hostname, phi.virtual, phi.cluster_name, lic->'Name' AS " 
		+ "license_name, CAST(CAST(lic->'Count' AS TEXT) AS REAL) AS license_count FROM host_info phi, "
		+ "jsonb_array_elements(phi.licenses) AS lic), aggregated_featured_host_info AS (SELECT " 
		+ "fhi.hostname, fhi.virtual, fhi.cluster_name, fhi.license_name, max(fhi.license_count) AS " 
		+ "license_count FROM featured_host_info fhi GROUP BY fhi.hostname, fhi.license_name, " 
		+ "fhi.cluster_name, fhi.virtual), summed_phy_features AS (SELECT afhi.license_name, " 
		+ "sum(afhi.license_count) AS license_count FROM aggregated_featured_host_info afhi WHERE "
		+ "virtual = false GROUP BY license_name), virtual_featured_host_info AS (SELECT " 
		+ "afhi.cluster_name, ci.cpu, afhi.license_name, max(afhi.license_count) AS license_count "
		+ "FROM aggregated_featured_host_info afhi LEFT JOIN cluster_info ci ON ci.name = afhi.cluster_name "
		+ "WHERE afhi.virtual = true GROUP BY afhi.cluster_name, afhi.virtual, afhi.license_name, "
		+ "ci.cpu), summed_virtual_features AS (SELECT vfhi.license_name, sum( CASE WHEN "
		+ "vfhi.license_count > 0 THEN vfhi.cpu*0.50 ELSE 0 END ) AS license_count FROM " 
		+ "virtual_featured_host_info vfhi GROUP BY license_name), partial_summed_featured AS (SELECT "
		+ "* FROM summed_phy_features UNION ALL SELECT * FROM summed_virtual_features), summed_featured "
		+ "AS (SELECT regexp_replace(CAST(license_name AS TEXT),'\"','','g') AS name, ceil(sum(license_count)) "
		+ "AS sum FROM partial_summed_featured GROUP BY license_name) SELECT * FROM summed_featured")
	List<Map<String, Object>> getAllLicensesCount(@Param("location") String location);

	/**
	 * Gets all hosts using license.
	 *
	 * @param license the license
	 * @return the list of hosts using the licensegetAllLicensesCount
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH dbs AS ( "
		+ "	SELECT "
		+ "		hostname, " 
		+ "		db->>'Name' AS dbname, "
		+ "		db->'Licenses' as lics "
		+ "	FROM "
		+ "		current_host, "
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db "
		+ "), licenses AS ( "
		+ "	SELECT "
		+ "		hostname, "
		+ "		dbname "
		+ "	FROM "
		+ "		dbs, "
		+ "		jsonb_array_elements(lics) AS lic "
		+ "	WHERE "
		+ "		CAST((lic->>'Count') AS numeric) > 0 AND "
		+ "		lic->>'Name' = :license "
		+ ") SELECT "
		+ "	hostname, " 
		+ "	string_agg(dbname, ' ') AS dbs "
		+ "FROM licenses "
		+ "GROUP BY hostname;")
	List<Map<String, Object>> getAllHostsUsingLicense(@Param("license") String license);

	/**
	 * Gets compliance.
	 *
	 * @return status of Licenses: true if you've got more Licenses than used, false otherwise
	 */
	@Query(nativeQuery = true, value = "WITH host_info AS (SELECT ch.hostname,"
		+	" ch.associated_cluster_name IS NOT NULL AS virtual, "
		+ "dbs->'Licenses' AS licenses, ch.associated_cluster_name AS cluster_name FROM current_host ch, "
		+ "jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS dbs WHERE (ch.host_type IS "
		+ "NULL OR ch.host_type = 'oracledb') AND ('*' = '*' or ch.location = '*')), featured_host_info "
		+ "AS (SELECT phi.hostname, phi.virtual, phi.cluster_name, lic->'Name' AS license_name, " 
		+ "CAST(CAST(lic->'Count' AS TEXT) AS REAL) AS license_count FROM host_info phi, "
		+ "jsonb_array_elements(phi.licenses) AS lic), aggregated_featured_host_info AS (SELECT "
		+ "fhi.hostname, fhi.virtual, fhi.cluster_name, fhi.license_name, max(fhi.license_count) AS " 
		+ "license_count FROM featured_host_info fhi GROUP BY fhi.hostname, fhi.license_name, " 
		+ "fhi.cluster_name, fhi.virtual), summed_phy_features AS (SELECT afhi.license_name, " 
		+ "sum(afhi.license_count) AS license_count FROM aggregated_featured_host_info afhi WHERE " 
		+ "virtual = false GROUP BY license_name), virtual_featured_host_info AS (SELECT afhi.cluster_name, "
		+ "ci.cpu, afhi.license_name, max(afhi.license_count) AS license_count FROM " 
		+ "aggregated_featured_host_info afhi LEFT JOIN cluster_info ci ON ci.name = afhi.cluster_name "
		+ "WHERE afhi.virtual = true GROUP BY afhi.cluster_name, afhi.virtual, afhi.license_name, ci.cpu "
		+ "), summed_virtual_features AS (SELECT vfhi.license_name, sum(CASE WHEN vfhi.license_count > 0 " 
		+ "THEN vfhi.cpu*0.50 ELSE 0 END ) AS license_count FROM virtual_featured_host_info vfhi GROUP BY "
		+ "license_name), partial_summed_featured AS (SELECT * FROM summed_phy_features UNION ALL SELECT "
		+ "* FROM summed_virtual_features), summed_featured AS (SELECT regexp_replace(CAST(license_name "
		+ "AS TEXT),'\"','','g') AS name, sum(license_count) AS sum FROM partial_summed_featured GROUP BY "
		+ "license_name), checked_feature AS (SELECT sf.name, lic.license_count >= sf.sum AS result FROM "
		+ "summed_featured sf LEFT JOIN license lic ON lic.id = sf.name) SELECT result FROM checked_feature")
	List<Boolean> getCompliance();


	/**
	 * Gets total licenses for compliance.
	 *
	 * @return count of total charged(free) licenses and used licenses
	 */
	@Query(nativeQuery = true, value = "WITH host_info AS (SELECT ch.hostname, "
			+ "ch.associated_cluster_name IS NOT NULL AS virtual, "
		+ "dbs->'Licenses' AS licenses, ch.associated_cluster_name AS cluster_name FROM current_host ch, "
		+ "jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS dbs WHERE (ch.host_type IS " 
		+ "NULL OR ch.host_type = 'oracledb') AND ('*' = '*' or ch.location = '*')), featured_host_info " 
		+ "AS (SELECT phi.hostname, phi.virtual, phi.cluster_name, lic->'Name' AS license_name, " 
		+ "CAST(CAST(lic->'Count' AS TEXT) AS REAL) AS license_count FROM  host_info phi, " 
		+ "jsonb_array_elements(phi.licenses) AS lic), aggregated_featured_host_info AS (SELECT " 
		+ "fhi.hostname, fhi.virtual, fhi.cluster_name, fhi.license_name, max(fhi.license_count) AS " 
		+ "license_count FROM featured_host_info fhi GROUP BY fhi.hostname, fhi.license_name, " 
		+ "fhi.cluster_name, fhi.virtual), summed_phy_features AS (SELECT afhi.license_name, " 
		+ "sum(afhi.license_count) AS license_count FROM aggregated_featured_host_info afhi WHERE " 
		+ "virtual = false GROUP BY license_name), virtual_featured_host_info AS (SELECT afhi.cluster_name, "
		+ "ci.cpu, afhi.license_name, max(afhi.license_count) AS license_count FROM " 
		+ "aggregated_featured_host_info afhi LEFT JOIN cluster_info ci ON ci.name = afhi.cluster_name "
		+ "WHERE afhi.virtual = true GROUP BY afhi.cluster_name, afhi.virtual, afhi.license_name, ci.cpu "
		+ "), summed_virtual_features AS (SELECT vfhi.license_name, sum(CASE WHEN vfhi.license_count > 0 " 
		+ "THEN vfhi.cpu*0.50 ELSE 0 END ) AS license_count FROM virtual_featured_host_info vfhi GROUP BY "
		+ "license_name), partial_summed_featured AS (SELECT * FROM summed_phy_features UNION ALL SELECT "
		+ "* FROM summed_virtual_features), summed_featured AS (SELECT regexp_replace(CAST(license_name "
		+ "AS TEXT),'\"','','g') AS name, sum(license_count) AS sum FROM partial_summed_featured GROUP BY "
		+ "license_name), checked_feature AS (SELECT sum(sf.sum) AS used, sum(lic.license_count) AS free "
		+ "FROM summed_featured sf LEFT JOIN license lic ON lic.id = sf.name) SELECT * FROM checked_feature")
	Map<String, BigInteger> getTotalLicensesForCompliance();

	/**
	 * Get all addms.
	 * @param env the env filter
	 * @param search the search filter
	 * @return the list of addms
	 */
	@Query(nativeQuery =  true, value = ""
		+ "WITH host_database AS ("
		+ "	SELECT "
		+ "		ch.hostname,"
		+ "		db->'ADDMs' AS addms,"
		+ "		ch.environment,"
		+ "		db->>'Name' AS dbName"
		+ "	FROM "
		+ "		current_host ch,"
		+ "		jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db"
		+ "	WHERE "
		+ "		db->'ADDMs' IS NOT NULL AND"
		+ "		db->>'ADDMs' != 'null' AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND "
		+ "		("
		+ "			LOWER(db->>'Name') LIKE LOWER(CONCAT('%',:search,'%')) OR "
		+ "			LOWER(ch.hostname) LIKE LOWER(CONCAT('%',:search,'%'))"
		+ "		) AND"
		+ "		LOWER(ch.environment) LIKE LOWER(CONCAT('%',:env,'%'))"
		+ "), addms AS ("
		+ "	SELECT"
		+ "		hb.hostname,"
		+ "		hb.environment,"
		+ "		hb.dbname,"
		+ "		addm->>'Benefit' AS benefit,"
		+ "		addm->>'Finding' AS finding,"
		+ "		addm->>'Recommendation' AS recommendation,"
		+ "		addm->>'Action'AS action"
		+ "	FROM "
		+ "		host_database hb,"
		+ "		jsonb_array_elements(addms) AS addm"		
		+ ") SELECT * FROM addms")
	List<Map<String, Object>> getADDMs(@Param("env") String env, @Param("search") String search);

	/**
	 * Get all segment advisors.
	 * @param env the env filter
	 * @param search the search filter
	 * @return the list of segment advisors
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH host_database AS ("
		+ "	SELECT"
		+ "		ch.hostname,"
		+ "		db->'SegmentAdvisors' AS segmentAdvisors,"
		+ "		ch.environment,"
		+ "		db->>'Name' AS dbName"
		+ "	FROM "
		+ "		current_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE"
		+ "		db->'SegmentAdvisors' IS NOT NULL AND"
		+ "		db->>'SegmentAdvisors' != 'null' AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		("
		+ "			LOWER(db->>'Name') LIKE LOWER(CONCAT('%',:search,'%')) OR "
		+ "			LOWER(ch.hostname) LIKE LOWER(CONCAT('%',:search,'%'))"
		+ "		) AND"
		+ "		LOWER(ch.environment) LIKE LOWER(CONCAT('%',:env,'%')) AND"
		+ "		(db->'InstanceNumber' IS NULL OR db->>'InstanceNumber' = '1')"
		+ "), addm AS ("
		+ "	SELECT"
		+ "		hb.hostname,"
		+ "		hb.environment,"
		+ "		hb.dbName,"
		+ "		CASE WHEN (segmentAdvisor->>'Reclaimable' = '<1') THEN CAST(0.5 AS real) ELSE CAST((segmentAdvisor->>'Reclaimable') AS real) END AS reclaimable,"
		+ "		segmentAdvisor->>'SegmentName' AS segmentName,"
		+ "		segmentAdvisor->>'SegmentOwner' AS segmentOwner,"
		+ "		segmentAdvisor->>'SegmentType' AS segmentType,"
		+ "		segmentAdvisor->>'PartitionName' AS partitionName,"
		+ "		segmentAdvisor->>'Recommendation' AS recommendation"
		+ "	FROM"
		+ "		host_database hb,"
		+ "		jsonb_array_elements(segmentAdvisors) AS segmentAdvisor"
		+ "	WHERE"
		+ "		segmentAdvisor->>'Reclaimable' != '-'"
		+ ") SELECT * FROM addm")
	List<Map<String, Object>> getSegmentAdvisors(@Param("env") String env, @Param("search") String search);
	
	/**
	 * Get all hosts psu status.
	 * @param windowTime window time
	 * @param status status
	 * @return	psu status of all hosts
	 */
	@Query(nativeQuery = true, value = ""
		+ "	WITH host_database AS (SELECT"
		+ "			ch.hostname,"
		+ "			db->>'Name' AS dbName,"
		+ "			db->>'Version' AS dbVer,"
		+ "			db"
		+ "		FROM" 
		+ "			current_host ch,"
		+ "			jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db"
		+ "		WHERE "
		+ "			(ch.host_type IS NULL OR ch.host_type = 'oracledb')"
		+ "	), host_database_without_psu AS (SELECT"
		+ "			hb.hostname,"
		+ "			hb.dbName,"
		+ "			hb.dbVer,"
		+ "			text('') AS psuDescription,"
		+ "			date'0001-01-01' AS psuDate,"
		+ "			text('KO') AS status"
		+ "		FROM"
		+ "			host_database hb"
		+ "		WHERE"
		+ "			hb.db->'LastPSUs' IS NULL OR" 
		+ "			hb.db->>'LastPSUs' = 'null' OR"
		+ "			jsonb_array_length(hb.db->'LastPSUs') = 0"
		+ "	), host_database_with_psu AS (SELECT"
		+ "			hd.hostname,"
		+ "			hd.dbname,"
		+ "			hd.dbver,"
		+ "			CAST(psu->>'Date' AS date) AS psuDate,"
		+ "			psu->>'Description' as psuDescription"
		+ "		FROM"
		+ "			host_database hd,"
		+ "			jsonb_array_elements(hd.db->'LastPSUs') AS psu"
		+ "		WHERE"
		+ "			hd.db->'LastPSUs' IS NOT NULL AND"
		+ "			hd.db->>'LastPSUs' != 'null' AND"
		+ "			jsonb_array_length(hd.db->'LastPSUs') > 0 AND"
		+ "			psu->>'Date' != 'N/A'"
		+ "), host_databases_last_psu_per_date AS (SELECT"
		+ "			hdwp.hostname,"
		+ "			hdwp.dbname,"
		+ "			hdwp.dbver,"
		+ "			max(hdwp.psuDate) AS psuDate"
		+ "		FROM"
		+ "			host_database_with_psu hdwp"
		+ "		GROUP BY"
		+ "			hdwp.hostname, hdwp.dbname, hdwp.dbver"
		+ "	),host_databases_with_last_psu AS (SELECT"
		+ "			hdwp.hostname,"
		+ "			hdwp.dbname,"
		+ "			hdwp.dbver,"
		+ "			hdwp.psuDescription,"
		+ "			hdwp.psuDate,"
		+ "			CASE WHEN hdwp.psuDate >= (date (:windowTime)) THEN"
		+ "				text('OK')"
		+ "			ELSE"
		+ "				text('KO')"
		+ "			END AS status"
		+ "		FROM "
		+ "			host_database_with_psu hdwp"
		+ "				INNER JOIN"
		+ "					host_databases_last_psu_per_date hdpsu" 
		+ "				ON"
		+ "					hdwp.hostname = hdpsu.hostname AND"
		+ "					hdwp.dbname = hdpsu.dbname AND"
		+ "					hdwp.dbver = hdpsu.dbver"
		+ "		WHERE"
		+ "			hdwp.psuDate = hdpsu.psuDate"
		+ ") SELECT * FROM host_databases_with_last_psu WHERE status = UPPER(:status)"
        + " OR :status = 'all' UNION ALL SELECT * FROM host_database_without_psu WHERE "
        + "status = UPPER(:status) OR :status = 'all';")
	List<Map<String, Object>> getAllHostPSUStatus(
	        @Param("windowTime") Date windowTime,
            @Param("status") String status);

	/**
	 * Get the top 15 reclaimable database.
	 * @param location location
	 * @return Top 15 reclaimable databas
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH host_database AS ("
		+ "	SELECT"
		+ "		ch.hostname,"
		+ "		db->'SegmentAdvisors' AS segmentAdvisors,"
		+ "		ch.environment,"
		+ "		db->>'Name' AS dbName"
		+ "	FROM"
		+ "		current_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		('*' = :location or ch.location = :location) AND"
		+ "		db->'SegmentAdvisors' IS NOT NULL AND"
		+ "		db->>'SegmentAdvisors' != 'null' AND"
		+ "		(db->'InstanceNumber' IS NULL OR db->>'InstanceNumber' = '1')"
		+ "), segment_advisors AS ("
		+ "	SELECT"
		+ "		hb.dbName,"
		+ "		hb.hostname,"
		+ "		( CASE WHEN (jsonb_typeof(segmentAdvisor->'Reclaimable') = 'number') OR (jsonb_typeof(segmentAdvisor->'Reclaimable') = 'string' AND segmentAdvisor->>'Reclaimable' ~ '^[+-]?([0-9]*[.])?[0-9]+$') THEN"
		+ "			CAST((segmentAdvisor->>'Reclaimable') AS real)"
		+ "	      WHEN jsonb_typeof(segmentAdvisor->'Reclaimable') = 'string' THEN"
		+ "			0.5"
		+ "		  ELSE"
		+ "			0"
		+ "		  END"
		+ "		) AS reclaimable"
		+ "	FROM"
		+ "		host_database hb,"
		+ "		jsonb_array_elements(segmentAdvisors) AS segmentAdvisor"
		+ ") SELECT CONCAT(hostname, ' ', dbname) AS dbname, SUM(reclaimable)"
        + " as sum FROM segment_advisors GROUP BY dbname,hostname ORDER BY sum desc LIMIT 15;")
	List<Map<String, Object>> getTopReclaimableDatabase(@Param("location") String location);

	/**
	 * Get patch status stats.
	 * @param location location
	 * @param windowTime windowTime
	 * @return patch status stats
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH host_database AS (SELECT"
		+ "		ch.hostname,"
		+ "		db->>'Name' AS dbName,"
		+ "		db "
		+ "	FROM"
		+ "		current_host ch,"
		+ "		jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db"
		+ "	WHERE "
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		('*' = :location or ch.location = :location)"
		+ "), host_database_without_psu AS (SELECT"
		+ "		text('KO') AS status"
		+ "	FROM"
		+ "		host_database hb"
		+ "	WHERE"
		+ "		hb.db->'LastPSUs' IS NULL OR"
		+ "		hb.db->>'LastPSUs' = 'null' OR"
		+ "		jsonb_array_length(hb.db->'LastPSUs') = 0"
		+ "), host_database_with_psu AS (SELECT"
		+ "		hd.hostname,"
		+ "		hd.dbname,"
		+ "		CAST(psu->>'Date' AS date) AS psuDate,"
		+ "		psu->>'Description' as psuDescription"
		+ "	FROM"
		+ "		host_database hd,"
		+ "		jsonb_array_elements(hd.db->'LastPSUs') AS psu"
		+ "	WHERE"
		+ "		hd.db->'LastPSUs' IS NOT NULL AND"
		+ "		hd.db->>'LastPSUs' != 'null' AND"
		+ "		jsonb_array_length(hd.db->'LastPSUs') > 0 AND"
		+ "		psu->>'Date' != 'N/A'"
		+ "	), host_databases_last_psu_per_date AS (SELECT"
		+ "		hdwp.hostname,"
		+ "		hdwp.dbname,"
		+ "		max(hdwp.psuDate) AS psuDate"
		+ "	FROM"
		+ "		host_database_with_psu hdwp"
		+ "	GROUP BY"
		+ "		hdwp.hostname, hdwp.dbname"
		+ "), host_databases_with_last_psu AS (SELECT"
		+ "		CASE WHEN hdwp.psuDate >= :windowTime THEN"
		+ "			text('OK')"
		+ "		ELSE"
		+ "			text('KO')"
		+ "		END AS status"
		+ "	FROM "
		+ "		host_database_with_psu hdwp"
		+ "		INNER JOIN"
		+ "			host_databases_last_psu_per_date hdpsu"
		+ "		ON"
		+ "			hdwp.hostname = hdpsu.hostname AND"
		+ "			hdwp.dbname = hdpsu.dbname"
		+ "	WHERE"
		+ "		hdwp.psuDate = hdpsu.psuDate"
		+ "), all_status AS (SELECT"
		+ "		status"
		+ "	FROM "
		+ "		host_databases_with_last_psu "
		+ "	UNION ALL SELECT "
		+ "		status"
		+ "	FROM "
		+ "		host_database_without_psu"
		+ ") SELECT status, count(*) FROM all_status GROUP BY status;")
	List<Map<String, Object>> getPatchStatusStats(@Param("location") String location,
		@Param("windowTime") Date windowTime);

	/**
	 * Get the environments.
	 * @return the environments.
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT"
		+ "	environment "
		+ "FROM"
		+ "	current_host "
		+ "GROUP BY"
		+ "	environment"
		)
	List<String> getEnviroments();

	/**
	 * Return the 'used' data history of the hostname/db.
	 * @param hostname hostname
	 * @param dbname dbname
	 * @return the 'used' data history of the hostname/db.
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH data AS ("
		+ "	SELECT "
		+ "		updated,"
		+ "		db->>'Used' AS used"
		+ "	FROM "
		+ "		current_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname "
		+ "	UNION ALL SELECT "
		+ "		updated,"
		+ "		db->>'Used' AS used"
		+ "	FROM "
		+ "		historical_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname"
		+ ") SELECT * FROM data ORDER BY updated ASC;")
	List<Map<String, Object>> getUsedDataHistory(@Param("hostname") final String hostname, @Param("dbname") final String dbname);

	/**
	 * Return the 'segmentsSize' data history of the hostname/db.
	 * @param hostname hostname
	 * @param dbname dbname
	 * @return the 'segmentsSize' data history of the hostname/db.
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH data AS ("
		+ "	SELECT "
		+ "		updated,"
		+ "		db->>'SegmentsSize' AS segmentsSize"
		+ "	FROM "
		+ "		current_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname AND"
		+ "		db->'SegmentsSize' IS NOT NULL AND" 
		+ "		db->>'SegmentsSize' != 'null'"
		+ "	UNION ALL SELECT "
		+ "		updated,"
		+ "		db->>'SegmentsSize' AS segmentsSize"
		+ "	FROM "
		+ "		historical_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname AND"
		+ "		db->'SegmentsSize' IS NOT NULL AND" 
		+ "		db->>'SegmentsSize' != 'null'"
		+ ") SELECT * FROM data ORDER BY updated ASC;")
	List<Map<String, Object>> getSegmentsSizeDataHistory(@Param("hostname") final String hostname, @Param("dbname") final String dbname);

/**
	 * Return the 'used' data history of the hostname/db.
	 * @param hostname hostname
	 * @param dbname dbname
	 * @return the 'used' data history of the hostname/db.
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH data AS ("
		+ "	SELECT "
		+ "		updated,"
		+ "		db->>'DailyCPUUsage' AS usage"
		+ "	FROM "
		+ "		current_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname AND db->>'DailyCPUUsage' IS NOT NULL"
		+ "	UNION ALL SELECT "
		+ "		updated,"
		+ "		db->>'DailyCPUUsage' AS usage"
		+ "	FROM "
		+ "		historical_host ch,"
		+ "		jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db"
		+ "	WHERE "
		+ "		hostname = :hostname AND"
		+ "		(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND"
		+ "		db->>'Name' = :dbname AND db->>'DailyCPUUsage' IS NOT NULL"
		+ ") SELECT * FROM data d WHERE NOT EXISTS ( "
		+ " SELECT * FROM data d2 WHERE d2.updated > d.updated AND  date(d2.updated) = date(d.updated) "
		+ ") ORDER BY updated ASC; ")
	List<Map<String, Object>> getDailyCPUUsageDataHistory(@Param("hostname") final String hostname, @Param("dbname") final String dbname);


	/**
	 * Return the list of databases.
	 * @param c pageable
	 * @param env env
	 * @return the list of databases
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	ch.hostname, "
		+ "	ch.environment, "
		+ "	ch.location, "
		+ "	db->>'Name' AS dbname, "
		+ "	db->>'UniqueName' AS unique_name, "
		+ "	db->>'Version' AS dbver, "
		+ "	db->>'Status' AS status, "
		+ "	db->>'Charset' AS charset, "
		+ "	CAST(db->>'BlockSize' AS INT) AS block_size, "
		+ "	CAST(db->>'CPUCount' AS INT) AS cpu_count, "
		+ "	db->>'Work' AS work, "
		+ "	(CAST(db->>'PGATarget' AS REAL) + "
		+ "		CAST(db->>'SGATarget' AS REAL) + "
		+ "		(CASE WHEN db->>'MemoryTarget' = '' THEN "
		+ "			0 "
		+ "		ELSE "
		+ "			CAST(db->>'MemoryTarget' AS REAL) "
		+ "		END) "
		+ "	) AS memory, "
		+ "	CAST(db->>'Used' AS real) as datafile_size, "
		+ "	CAST(db->>'SegmentsSize' AS real) as segments_size, "
		+ "	(db->>'Archivelog' = 'ARCHIVELOG') as archive_log_status, "
		+ "	CAST(db->>'Dataguard' AS boolean) AS dataguard, "
		+ "	( "
		+ "		SELECT  "
		+ "			CAST(fe->>'Status' AS boolean) "
		+ "		FROM  "
		+ "			jsonb_array_elements(db->'Features') AS fe "
		+ "		WHERE  "
		+ "			fe->>'Name' = 'Real Application Clusters' "
		+ "	) AS rac, "
		+ "	CAST(hi->>'SunCluster' AS boolean) OR "
		+ "	CAST(hi->>'VeritasCluster' as boolean) OR "
		+ "	CAST(hi->>'OracleCluster' as boolean) OR "
		+ "	(CASE WHEN hi->>'AixCluster' IS NULL THEN  "
		+ "		false "
		+ "	ELSE "
		+ "		CAST(hi->>'AixCluster' AS boolean) "
		+ "	END) AS ha "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db, "
		+ "	CAST(ch.host_info AS jsonb) AS hi "
		+ "WHERE  "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) ")
	Page<Map<String, Object>> getDatabases(Pageable c, String env);

	/**
	 * Count the databases grouped by dataguard status.
	 * @param env env
	 * @return the count of databases grouped by dataguard status
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	CAST(db->>'Dataguard' AS bool) as status,"
		+ "	COUNT(*) AS count "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db "
		+ "WHERE "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) "
		+ "GROUP BY "
		+ "	db->>'Dataguard'"
	)
	List<Map<String, Object>> countDatabaseGroupedByDataguardStatus(final String env);

	/**
	 * Count the databases grouped by real application cluster feature status.
	 * @param env env
	 * @return the count of databases grouped by real application cluster feature status
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	CAST(fe->>'Status' AS bool) AS status, "
		+ "	COUNT(*) "
		+ "FROM  "
		+ "	current_host ch, "
		+ "	jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db, "
		+ "	jsonb_array_elements(db->'Features') AS fe "
		+ "WHERE  "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) "
		+ "	AND fe->>'Name' = 'Real Application Clusters' "
		+ "GROUP BY "
		+ "	fe->>'Status'; ")
	List<Map<String, Object>> countDatabasesGroupedByRealApplicationClusterFeatureStatus(final String env);

	/**
	 * Count the databases grouped by archive log status.
	 * @param env env
	 * @return the count of databases grouped by archive log status
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	db->>'Archivelog' = 'ARCHIVELOG' AS status, "
		+ "	count(*) "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db "
		+ "WHERE  "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) "
		+ "GROUP BY "
		+ "	db->>'Archivelog' = 'ARCHIVELOG'")
	List<Map<String, Object>> countDatabasesGroupedByArchiveLogStatus(final String env);

	/**
	 * Return the sum of the segments size.
	 * @param env env
	 * @return the sum of the segments size
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	SUM(CAST(db->>'SegmentsSize' AS REAL)) AS SegmentSize "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db "
		+ "WHERE "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) "
		+ "	AND db->>'SegmentsSize' != 'N/A'")
	float getTotalSegmentsSize(String env);

	/**
	 * Return the sum of the datafile size.
	 * @param env env
	 * @return the sum of the datafile size
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	sum(CAST(db->>'Used' AS real)) as datafile_size "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db "
		+ "WHERE "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) "
		+ "	AND db->>'Used' != 'N/A'")
	float getTotalDatafileSize(final String env);

	/**
	 * Return the sum of the memory size.
	 * @param env env
	 * @return the sum of the memory size
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	SUM(CAST(db->>'PGATarget' AS REAL) + "
		+ "		CAST(db->>'SGATarget' AS REAL) + "
		+ "		(CASE WHEN db->>'MemoryTarget' = '' THEN "
		+ "			0 "
		+ "		ELSE "
		+ "			CAST(db->>'MemoryTarget' AS REAL) "
		+ "		END) "
		+ "	) AS memory "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements((CAST(extra_info AS jsonb))->'Databases') AS db "
		+ "WHERE  "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') "
		+ " AND (:env = '' OR ch.environment = :env) ")
	float getTotalMemorySize(String env);

	/**
	 * Return the sum of the database work.
	 * @param env env
	 * @return the sum of the database work
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	SUM(CAST(db->>'Work' AS REAL)) AS Work "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements(CAST(extra_info AS jsonb)->'Databases') AS db "
		+ "WHERE "
		+ "	(ch.host_type IS NULL OR ch.host_type = 'oracledb') AND "
		+ " (:env = '' OR ch.environment = :env) AND"
		+ "	db->>'Work' != 'N/A' ")
	float getTotalDatabaseWork(final String env);


	/**
	 * Return the exadata devices.
	 * @return the exadata devices
	 */
	@Query(nativeQuery = true, value = ""
		+ "SELECT "
		+ "	ch.hostname, "
		+ "	exadev->>'Hostname' AS dev_hostname, "
		+ "	exadev->>'ServerType' AS dev_type, "
		+ "	exadev->>'Model' AS dev_model, "
		+ "	exadev->>'CPUEnabled' AS dev_cpu, "
		+ "	exadev->>'Memory' AS dev_memory, "
		+ "	exadev->>'ExaSwVersion' AS dev_version, "
		+ "	exadev->>'PowerCount' AS dev_power_count, "
		+ "	exadev->>'TempActual' AS dev_temp_actual "
		+ "FROM "
		+ "	current_host ch, "
		+ "	jsonb_array_elements((CAST(extra_info AS jsonb))->'Exadata'->'Devices') AS exadev "
		+ "WHERE "
		+ "	ch.host_type = 'exadata'")
	List<Map<String, String>> findExadataDevices();

	/**
	 * Get exadata stats.
	 * @return exadata stats
	 */
	@Query(nativeQuery = true, value = ""
		+ "WITH parsed_data AS ( "
		+ "	SELECT "
		+ "		(CASE WHEN (exadev->>'CPUEnabled' = '-') THEN "
		+ "			0 "
		+ "		ELSE "
		+ "			CAST(regexp_replace(exadev->>'CPUEnabled', '^(\\d+)\\/(\\d+)$', '\\1') AS REAL) "
		+ "		END) AS dev_used_cpu, "
		+ "		(CASE WHEN (exadev->>'CPUEnabled' = '-') THEN "
		+ "			0 "
		+ "		ELSE "
		+ "			CAST(regexp_replace(exadev->>'CPUEnabled', '^(\\d+)\\/(\\d+)$', '\\2') AS REAL) "
		+ "		END) AS dev_total_cpu, "
		+ "		(CASE WHEN (exadev->>'Memory' = '-') THEN "
		+ "			0 "
		+ "		ELSE "
		+ "			CAST(regexp_replace(exadev->>'Memory', '^(\\d+)GB$', '\\1') AS REAL) "
		+ "		END) AS dev_used_memory, "
		+ "		(CASE WHEN (exadev->>'ExaSwVersion' = '-') THEN "
		+ "			current_date "
		+ "		ELSE "
		+ "			to_date(regexp_replace(exadev->>'ExaSwVersion', '^\\d+\\.\\d+\\.\\d+(\\.\\d+\\.\\d+)?(-\\d+)?\\.(\\d+)$', '\\3'), 'YYMMDD') "
		+ "		END) AS version_date "
		+ "	FROM "
		+ "		current_host ch, "
		+ "		jsonb_array_elements((CAST(extra_info AS jsonb))->'Exadata'->'Devices') AS exadev "
		+ "	WHERE "
		+ "		ch.host_type = 'exadata' "
		+ ") SELECT  "
		+ "	count(*) AS count, "
		+ "	sum(pd.dev_used_cpu) AS used_cpu, "
		+ "	sum(pd.dev_total_cpu) AS total_cpu, "
		+ "	sum(pd.dev_used_memory) AS memory, "
		+ "	sum(CASE WHEN (pd.version_date >= current_date - interval '6 months') THEN "
		+ "		1 "
		+ "	ELSE  "
		+ "		0 "
		+ "	END) AS count_patched_after_six_month, "
		+ "	sum(CASE WHEN (pd.version_date >= current_date - interval '12 months') THEN "
		+ "		1 "
		+ "	ELSE  "
		+ "		0 "
		+ "	END) AS count_patched_after_twelve_month "
		+ "FROM  "
		+ "	parsed_data pd; ")
	Map<String, Object> getExadataStats();

	@Query(nativeQuery = true, value = ""
		+ "WITH parsed_data AS ( "
		+ "	SELECT "
		+ "		COUNT(*) AS disks_count, "
		+ "		AVG(CAST(cd->>'UsedPerc' AS REAL)) AS disks_used_perc, "
		+ "		SUM((CASE WHEN (CAST(cd->>'ErrCount' AS REAL) > 0) THEN "
		+ "			1 "
		+ "		ELSE "
		+ "			0 "
		+ "		END)) AS failed_disks "
		+ "	FROM "
		+ "		current_host ch, "
		+ "		jsonb_array_elements((CAST(extra_info AS jsonb))->'Exadata'->'Devices') AS exadev, "
		+ "		jsonb_array_elements(exadev->'CellDisks') AS cd "
		+ "	WHERE "
		+ "		ch.host_type = 'exadata' AND exadev->>'ServerType' = 'StorageServer' "
		+ ") SELECT * FROM parsed_data")
	Map<String, Object> getExadataDisksStats();
}
