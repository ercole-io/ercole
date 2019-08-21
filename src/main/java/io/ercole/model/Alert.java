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

import javax.persistence.Entity;
import javax.persistence.EnumType;
import javax.persistence.Enumerated;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;
import javax.persistence.Lob;

import org.hibernate.annotations.Type;

/**
 * The Class Alert.
 */
@Entity
public class Alert {


	/** The id. */
	@Id
	@GeneratedValue(strategy = GenerationType.IDENTITY)
	private Long id;

	/** The hostname. */
	private String hostname;

	/** The code. */
	@Enumerated(EnumType.STRING)
	private AlertCode code;

	/** The severity. */
	@Enumerated(EnumType.STRING)
	private AlertSeverity severity;
	
	/** The status. */
	@Enumerated(EnumType.STRING)
	private AlertStatus status;

	/** The description. */
	@Lob
	@Type(type = "org.hibernate.type.TextType")
	private String description;
	
	/** The date. */
	private Date date;
	
	/**
	 * Instantiates a new alert.
	 */
	public Alert() {
		
	}
	/**
	 * Instantiates a new alert.
	 *
	 * @param name the hostname
	 * @param cod the code
	 * @param descr the description
	 * @param sever the severity
	 */
	public Alert(final String name, final AlertCode cod, final String descr, 
			final AlertSeverity sever) {
		this.hostname = name;
		this.code = cod;
		this.description = descr;
		this.severity = sever;
		date = new Date();
		status = AlertStatus.NEW;
	}
	
	
	



	/**
	 * Gets the id.
	 *
	 * @return the id
	 */
	public Long getId() {
		return id;
	}

	/**
	 * Sets the id.
	 *
	 * @param id the new id
	 */
	public void setId(final Long id) {
		this.id = id;
	}

	/**
	 * Gets the hostname.
	 *
	 * @return the hostname
	 */
	public String getHostname() {
		return hostname;
	}

	/**
	 * Sets the hostname.
	 *
	 * @param hostname the new hostname
	 */
	public void setHostname(final String hostname) {
		this.hostname = hostname;
	}

	/**
	 * Gets the code.
	 *
	 * @return the code
	 */
	public AlertCode getCode() {
		return code;
	}

	/**
	 * Sets the code.
	 *
	 * @param code the new code
	 */
	public void setCode(final AlertCode code) {
		this.code = code;
	}

	/**
	 * Gets the severity.
	 *
	 * @return the severity
	 */
	public AlertSeverity getSeverity() {
		return severity;
	}

	/**
	 * Sets the severity.
	 *
	 * @param severity the new severity
	 */
	public void setSeverity(final AlertSeverity severity) {
		this.severity = severity;
	}

	/**
	 * Gets the description.
	 *
	 * @return the description
	 */
	public String getDescription() {
		return description;
	}

	/**
	 * Sets the description.
	 *
	 * @param description the new description
	 */
	public void setDescription(final String description) {
		this.description = description;
	}

	/**
	 * Gets the date.
	 *
	 * @return the date
	 */
	public Date getDate() {
		return date;
	}

	/**
	 * Sets the date.
	 *
	 * @param date the new date
	 */
	public void setDate(final Date date) {
		this.date = date;
	}
	
	/**
	 * Get the status.
	 * @return the alerts's status
	 */
	public AlertStatus getStatus() {
		return status;
	}
	
	/**
	 * Set the status.
	 * @param status the alert status
	 */
	public void setStatus(final AlertStatus status) {
		this.status = status;
	}
	
	
	
	

}
