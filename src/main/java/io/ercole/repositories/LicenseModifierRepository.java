package io.ercole.repositories;

import org.springframework.data.repository.PagingAndSortingRepository;

import io.ercole.model.LicenseModifier;

public interface LicenseModifierRepository extends PagingAndSortingRepository<LicenseModifier, Long> {
    LicenseModifier findByHostnameAndDbnameAndLicenseName(String hostname, String dbname, String licenseName);
}
