package io.ercole.repositories;

import java.util.List;

import org.springframework.data.repository.PagingAndSortingRepository;

import io.ercole.model.LicenseModifier;

public interface LicenseModifierRepository extends PagingAndSortingRepository<LicenseModifier, Long> {
    LicenseModifier findByHostnameAndDbnameAndLicenseName(String hostname, String dbname, String licenseName);
    List<LicenseModifier> findByHostnameAndDbname(String hostname, String dbname);
}
