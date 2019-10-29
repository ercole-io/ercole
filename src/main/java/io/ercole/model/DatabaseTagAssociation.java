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

@Entity
public class DatabaseTagAssociation {

    /** The id. */
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    private String hostname;
    private String dbname;
    private String tag;

    public String getHostname() {
        return hostname;
    }

    public String getTag() {
        return tag;
    }

    public void setTag(final String tag) {
        this.tag = tag;
    }

    public String getDbname() {
        return dbname;
    }

    public void setDbname(final String dbname) {
        this.dbname = dbname;
    }

    public void setHostname(final String hostname) {
        this.hostname = hostname;
    }

    public DatabaseTagAssociation() {
    }

    public Long getId() {
        return id;
    }

    public void setId(final Long id) {
        this.id = id;
    }

    public DatabaseTagAssociation(final String hostname, final String dbname, final String tag) {
        this.hostname = hostname;
        this.dbname = dbname;
        this.tag = tag;
    }
}
