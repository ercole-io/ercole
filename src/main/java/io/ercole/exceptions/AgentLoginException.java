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

package io.ercole.exceptions;

/**
 * Exception for login problems with the agents.
 */
public class AgentLoginException extends Exception {

	private static final long serialVersionUID = 7742041245065496809L;

	/**
	 * @param msg to throw
	 */
	public AgentLoginException(final String msg) {
		super(msg);
	}
}
