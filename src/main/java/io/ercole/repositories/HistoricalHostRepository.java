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

import java.util.Date;
import java.util.List;

import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.CrudRepository;
import org.springframework.data.repository.query.Param;

import io.ercole.model.HistoricalHost;

/**
 * The Interface MacchinaRepository.
 */
public interface HistoricalHostRepository extends CrudRepository<HistoricalHost, Long> {
		
	/**
	 * @param hostname 
	 * @param data 
	 * @param pageable 
	 * @return List<HistoricalHost> with 0 or more HistoricalHost objects
	 */
	@Query("SELECT m FROM HistoricalHost m WHERE m.hostname = :hostname AND m.archived <= :date "
			+ "ORDER BY m.archived DESC")
	List<HistoricalHost> findByHostnameOrderByArchived(@Param("hostname") String hostname, 
			@Param("date") Date data, Pageable pageable);
	
	
	
	/**
	 * @param hostname 
	 * @param data 
	 * @return filtered List<HistoricalHost> with only 1 HistoricalHost
	 */
	default List<HistoricalHost> findFirstHostnameByArchivedDesc(String hostname, Date data) {
		return findByHostnameOrderByArchived(hostname, data, 
				PageRequest.of(0, 1, new Sort(Sort.Direction.DESC, "archived")));
	}

	/**
	 * @param lowerLimit 
	 * @return List<HistoricalHost> with "archived" attribute greater than lowerLimit
	 */
	@Query("SELECT m FROM HistoricalHost m WHERE m.archived <= :date")
	List<HistoricalHost> findAllOlderThan(@Param("date")Date lowerLimit);
}
