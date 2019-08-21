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

package io.ercole.config;

import java.util.Collection;
import java.util.EnumSet;
import java.util.Set;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.security.core.GrantedAuthority;
import org.springframework.security.core.authority.mapping.GrantedAuthoritiesMapper;
import org.springframework.stereotype.Component;

/**
 * LDAP Authorities mapper.
 * 
 */
@Component
public class LdapAuthoritiesMapper implements GrantedAuthoritiesMapper {
	private static final Log LOGGER = LogFactory.getLog(LdapAuthoritiesMapper.class);

	@Value("${auth.ad.role}")
	private String adRole;

	/* (non-Javadoc)
	 * @see org.springframework.security.core.authority.mapping.GrantedAuthoritiesMapper
	 * #mapAuthorities(java.util.Collection)
	 */
	@Override
	public Collection<LdapAuthority> mapAuthorities(final Collection<? extends GrantedAuthority> authorities) {
		Set<LdapAuthority> roles = EnumSet.noneOf(LdapAuthority.class); // empty EnumSet
		for (GrantedAuthority authority : authorities) {
			LOGGER.error(authority.getAuthority());
			if (adRole.equalsIgnoreCase(authority.getAuthority())) {
				roles.add(LdapAuthority.ROLE_USER);
			}
		}
		return roles;
	}

}
