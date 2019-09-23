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

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

/**
 * It contain info about a VM.
 */
@Entity
public class VMInfo {
    
    @Id
	@GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;
    
    //@Column(unique = true)
    private String name;
    private String clusterName;
    //@Column(unique = true)
    private String hostName;
    private String physicalHost;

    /**
     * @return the id
     */
    public Long getId() {
        return this.id;
    }

    /**
     * @param id the id to be set
     */
    public void setId(final Long id) {
        this.id = id;
    }

    /**
     * @return the name
     */
    public String getName() {
        return this.name;
    }
    /**
     * @param name the name
     */
    public void setName(final String name) {
        this.name = name;
    }
    /**
     * @return the cluster name
     */
    public String getClusterName() {
        return this.clusterName;
    }
    /**
     * @param clusterName the cluster name
     */
    public void setClusterName(final String clusterName) {
        this.clusterName = clusterName;
    }
    /**
     * @return the hostname
     */
    public String getHostName() {
        return this.hostName;
    }
    /**
     * @param hostName the hostname
     */
    public void setHostName(final String hostName) {
        this.hostName = hostName;
    }
    /**
     * @return the physicalHost
     */
    public String getPhysicalHost() {
        return this.physicalHost;
    }
    /**
     * @param physicalHost the physicalHost
     */
    public void setPhysicalHost(final String physicalHost) {
        this.physicalHost = physicalHost;
    }
    /**
     * Create a new VMInfo.
     * @param id the id
     * @param name the name
     * @param clusterName the cluster name
     * @param hostName the hostname
     * @param physicalHost the physical host
     */
    public VMInfo(final Long id, final String name, final String clusterName, final String hostName, final String physicalHost) {
        this.id = id;
        this.name = name;
        this.clusterName = clusterName;
        this.hostName = hostName;
        this.physicalHost = physicalHost;
    }
    /**
     * Create a new VMInfo.
     */
    public VMInfo() {
    }
}
