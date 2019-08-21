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

import java.util.Date;
import java.util.List;

import javax.persistence.CascadeType;
import javax.persistence.ElementCollection;
import javax.persistence.Entity;
import javax.persistence.FetchType;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.OneToMany;

/**
 * It contain info about a cluster.
 */
@Entity
public class ClusterInfo {

    /** The id. */
	@Id
	@GeneratedValue(strategy = GenerationType.SEQUENCE)
    private long id;
    //@Column(unique = true)
    private String name;
    private int cpu;
    private int sockets;   
    @ElementCollection(fetch = FetchType.LAZY)
    @OneToMany(cascade = CascadeType.ALL, orphanRemoval = true)
    private List<VMInfo> vms;

    private Date updated;

    /**
     * @return the id
     */
    public long getId() {
        return this.id;
    }
    /**
     * Set the id.
     * @param id id
     */
    public void setId(final long id) {
        this.id = id;
    }
    /**
     * @return the name
     */
    public String getName() {
        return this.name;
    }
    /**
     * @param name name
     */
    public void setName(final String name) {
        this.name = name;
    }
    /**
     * @return the cpu
     */
    public int getCpu() {
        return this.cpu;
    }
    /**
     * Set the cpu.
     * @param cpu the cpu
     */
    public void setCpu(final int cpu) {
        this.cpu = cpu;
    }
    /**
     * @return sockets
     */
    public int getSockets() {
        return this.sockets;
    }
    /**
     * Set the .
     * @param sockets sockets
     */
    public void setSockets(final int sockets) {
        this.sockets = sockets;
    }
    /**
     * @return the vms
     */
    public List<VMInfo> getVms() {
        return this.vms;
    }
    /***
     * @param vms vms
     */
    public void setVms(final List<VMInfo> vms) {
        this.vms = vms;
    }
    /**
     * 
     * @param id id
     * @param name name
     * @param cpu cpu
     * @param sockets sockets
     * @param vms vms
     * @param updated updated
     */
    public ClusterInfo(final long id, final String name, final int cpu,
                       final int sockets, final List<VMInfo> vms, final Date updated) {
        this.id = id;
        this.name = name;
        this.cpu = cpu;
        this.sockets = sockets;
        this.vms = vms;
        this.updated = updated;
    }

    /**
     * 
     */
    public ClusterInfo() {
    }

    /**
     * @return the updated
     */
    public Date getUpdated() {
        return updated;
    }

    /**
     * @param updated the updated to set
     */
    public void setUpdated(final Date updated) {
        this.updated = updated;
    }
}
