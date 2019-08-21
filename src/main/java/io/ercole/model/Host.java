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

import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Inheritance;
import javax.persistence.InheritanceType;
import javax.persistence.Lob;
import javax.persistence.MappedSuperclass;
import javax.validation.constraints.NotEmpty;

import org.hibernate.annotations.Type;

import com.fasterxml.jackson.annotation.JsonRawValue;

/**
 * Superclass of all Host tables.
 */
@MappedSuperclass
@Inheritance(strategy = InheritanceType.TABLE_PER_CLASS)
public abstract class Host {

	
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	private Long id;
	
	@NotEmpty
	private String hostname;
	
	private String environment;
	private String location;
	private String hostType;
	
	@Lob
	@Type(type = "org.hibernate.type.TextType")
	private String databases;
	
	@Lob
	@Type(type = "org.hibernate.type.TextType")
	private String schemas;
	
	@Lob
	@Type(type = "org.hibernate.type.TextType")
	@JsonRawValue
	private String extraInfo;
	
	@Lob
	@Type(type = "org.hibernate.type.TextType")
	@JsonRawValue
	private String hostInfo;
	
	private Date updated;
	
	/*IT CAN BE NULL! If the host isn't in a cluster, this should be null*/
	private String associatedClusterName;

	/**
	 * Instantiates a new host.
	 */
	protected Host() {
		
	}
		
	/**
	 * Instantiates a new host.
	 *
	 * @param ident the id
	 * @param name the hostname
	 * @param envir the environment
	 * @param loc the location
	 * @param hostType the type of use (oracle, virtualization,...)
	 * @param data the databases
	 * @param schem the schemas
	 * @param extra the extra info
	 * @param associatedClusterName associatedClusterName
	 * @param host the host info
	 * @param updat the updated
	 */
	public Host(final Long ident, final @NotEmpty String name, final String envir, 
			final String loc, final String hostType,
			final String data, final String schem, final String extra, 
			final String associatedClusterName, final String host, final Date updat) {
		super();
		this.id = ident;
		this.hostname = name;
		this.environment = envir;
		this.location = loc;
		this.hostType = hostType;
		this.databases = data;
		this.schemas = schem;
		this.extraInfo = extra;
		this.associatedClusterName = associatedClusterName;
		this.hostInfo = host;
		this.updated = updat;
	}
	
	

	/**
	 * Gets the id.
	 *
	 * @return id
	 */
	public Long getId() {
		return id;
	}



	/**
	 * Sets the id.
	 *
	 * @param id setter
	 */
	public void setId(final Long id) {
		this.id = id;
	}



	/**
	 * Gets the hostname.
	 *
	 * @return hostname
	 */
	public String getHostname() {
		return hostname;
	}

	/**
	 * Sets the hostname.
	 *
	 * @param hostname to set
	 */
	public void setHostname(final String hostname) {
		this.hostname = hostname;
	}

	/**
	 * Gets the environment.
	 *
	 * @return the environment
	 */
	public String getEnvironment() {
		return environment;
	}

	/**
	 * Sets the environment.
	 *
	 * @param environment the new environment
	 */
	public void setEnvironment(final String environment) {
		this.environment = environment;
	}

	/**
	 * Gets the location.
	 *
	 * @return the location
	 */
	public String getLocation() {
		return location;
	}

	/**
	 * Sets the location.
	 *
	 * @param location the new location
	 */
	public void setLocation(final String location) {
		this.location = location;
	}


	/** Gets the hostType.
	 * @return the hostType
	 */
	public String getHostType() {
		return hostType;
	}

	/**
	* Sets the HostType.
	*
	* @param hostType the new hostType
	*/
	public void setHostType(final String hostType) {
		this.hostType = hostType;
	}


	/**
	 * Gets the databases.
	 *
	 * @return the databases
	 */
	public String getDatabases() {
		return databases;
	}


	/**
	 * Sets the databases.
	 *
	 * @param databases the new databases
	 */
	public void setDatabases(final String databases) {
		this.databases = databases;
	}


	/**
	 * Gets the schemas.
	 *
	 * @return the schemas
	 */
	public String getSchemas() {
		return schemas;
	}


	/**
	 * Sets the schemas.
	 *
	 * @param schemas the new schemas
	 */
	public void setSchemas(final String schemas) {
		this.schemas = schemas;
	}

	/**
	 * Gets the extra.
	 *
	 * @return JSON infos regarding databases info like Features, Tablespaces, Schemas, Patches
	 */
	public String getExtraInfo() {
		return extraInfo;
	}

	/**
	 * Sets the extra.
	 *
	 * @param extraInfo = JSON infos regarding databases info like Features, Tablespaces, Schemas, Patches.
	 */
	public void setExtraInfo(final String extraInfo) {
		this.extraInfo = extraInfo;
	}

	/**
	 * Gets the specifiche host.
	 *
	 * @return JSON infos regarding host OS and HW.
	 */
	public String getHostInfo() {
		return hostInfo;
	}

	/**
	 * Sets the specifiche host.
	 *
	 * @param hostInfo = JSON infos regarding host OS and HW to set
	 */
	public void setHostInfo(final String hostInfo) {
		this.hostInfo = hostInfo;
	}

	/**
	 * Gets the ultimo aggiornamento.
	 *
	 * @return update timestamp
	 */
	public Date getUpdated() {
		return updated;
	}

	/**
	 * Sets the ultimo aggiornamento.
	 *
	 * @param updated = update timestamp to set
	 */
	public void setUpdated(final Date updated) {
		this.updated = updated;
	}

	/**
	 * @return the associated cluster name
	 */
	public String getAssociatedClusterName() {
		return this.associatedClusterName;
	}
	/**
	 * Set the associated cluster name.
	 * @param associatedClusterName the associated cluster name
	 */
	public void setAssociatedClusterName(final String associatedClusterName) {
		this.associatedClusterName = associatedClusterName;
	}

}
