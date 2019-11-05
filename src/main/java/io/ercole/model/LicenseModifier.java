package io.ercole.model;

import javax.persistence.Entity;
import javax.persistence.GeneratedValue;
import javax.persistence.GenerationType;
import javax.persistence.Id;

@Entity
public class LicenseModifier {

    /** The id. */
    @Id
    @GeneratedValue(strategy = GenerationType.SEQUENCE)
    private Long id;

    private String hostname;
    private String dbname;
    private String licenseName;
    private int newValue;

    public Long getId() {
        return id;
    }

    public void setId(final Long id) {
        this.id = id;
    }

    public String getHostname() {
        return hostname;
    }

    public void setHostname(final String hostname) {
        this.hostname = hostname;
    }

    public String getDbname() {
        return dbname;
    }

    public void setDbname(final String dbname) {
        this.dbname = dbname;
    }

    public String getLicenseName() {
        return licenseName;
    }

    public void setLicenseName(final String licenseName) {
        this.licenseName = licenseName;
    }

    public int getNewValue() {
        return newValue;
    }

    public void setNewValue(final int newValue) {
        this.newValue = newValue;
    }

    public LicenseModifier(final String hostname, final String dbname, final String licenseName, final int newValue) {
        this.hostname = hostname;
        this.dbname = dbname;
        this.licenseName = licenseName;
        this.newValue = newValue;
    }

    public LicenseModifier() {
    }
}
