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

package io.ercole.repositories;

import java.util.List;
import org.springframework.data.repository.PagingAndSortingRepository;

import io.ercole.model.DatabaseTagAssociation;

public interface DatabaseTagAssociationRepository extends PagingAndSortingRepository<DatabaseTagAssociation, Long> {
    List<DatabaseTagAssociation> findByHostnameAndDbname(String hostname, String dbname);
    List<DatabaseTagAssociation> findTagsByHostnameOrderByDbname(String hostname);
    DatabaseTagAssociation findOneByHostnameAndDbnameAndTag(String hostname, String dbname, String tag); 
}
