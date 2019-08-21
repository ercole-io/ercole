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
/**
 * The Enum AlertCode.
 */
public enum AlertCode {
	
	/** The new database. */
	NEW_DATABASE("New Database"), 
	
	/** The new option. */
	NEW_OPTION("New Option"), 
	
	/** The new license. */
	NEW_LICENSE("New License"), 
	
	/** The new server. */
	NEW_SERVER("New Server"),
	
	/** The no data from agent. */
	NO_DATA("No Data");

	/**
	 * Instantiates a new alert code.
	 *
	 * @param title the title
	 */
	AlertCode(final String title) {
		this.briefDescr = title;
	}
	
	/** The brief descr. */
	private final String briefDescr;
	
	/**
	 * Gets the brief descr.
	 *
	 * @return the brief descr
	 */
	public String getBriefDescr() {
		return briefDescr;
	}
	
}
