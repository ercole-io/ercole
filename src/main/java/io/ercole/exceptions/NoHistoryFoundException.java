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
 * Exception to be thrown when we don't have any historicalHost related to the searched hostname.
 */
public class NoHistoryFoundException extends Exception {

	private static final long serialVersionUID = -3038769143800085998L;

	/**
	 * @param msg to be thrown
	 */
	public NoHistoryFoundException(final String msg) {
		super(msg);
	}

}
