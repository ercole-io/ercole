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

import java.util.List;

import org.springframework.data.jpa.repository.Query;
import org.springframework.data.repository.PagingAndSortingRepository;
import org.springframework.data.repository.query.Param;

import io.ercole.model.ClusterInfo;
import io.ercole.model.VMInfo;

/**
 * Repository for clusters.
 */
public interface ClusterRepository extends PagingAndSortingRepository<ClusterInfo, Long> {

    /**
     * Remove the cluster by name.
     *
     * @param name name of the cluster
     */
    void deleteByName(String name);

    /**
     * Find a cluster by name.
     *
     * @param name name of the cluster
     * @return the cluster info
     */
    ClusterInfo findByName(String name);

    /**
     * Find a VM containing the hostname.
     * @param name name
     * @return the vminfo
     */
    @Query("select vi from VMInfo vi where vi.hostName = LOWER(:#{#name})")
    VMInfo findOneVMInfoByHostnameIgnoreCase(@Param("name") String name);


    /**
     * Find a VM containing the hostname.
     * @param name name
     * @return the vminfo
     */
    @Query(value = "SELECT "
        + "	vi \n"
        + "FROM VMInfo vi \n"
        + "WHERE LOWER(CASE locate('.', vi.hostName) "
        + "    WHEN 0 THEN vi.hostName  " 
        + "    ELSE substring(vi.hostName, 1, locate('.', vi.hostName)-1)  "
        + "  END) = LOWER(:#{#name}) OR vi.hostName = :#{#name}\n")
    VMInfo findOneVMInfoByHostnameIgnoreCaseTrimDomain(@Param("name") String name);


    /**
     * Return all cluster that match the filter.
     * @param filter filter
     * @return all cluster that match the filter
     */
    @Query("SELECT cl from ClusterInfo cl WHERE LOWER(cl.name) LIKE LOWER (CONCAT('%',:filter,'%'))")
    List<ClusterInfo> getClusters(String filter);
}
