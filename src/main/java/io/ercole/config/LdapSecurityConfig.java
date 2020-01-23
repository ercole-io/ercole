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

import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.security.authentication.AuthenticationProvider;
import org.springframework.security.config.annotation.authentication.builders.AuthenticationManagerBuilder;
import org.springframework.security.config.annotation.web.builders.HttpSecurity;
import org.springframework.security.config.annotation.web.configuration.EnableWebSecurity;
import org.springframework.security.config.annotation.web.configuration.WebSecurityConfigurerAdapter;
import org.springframework.security.ldap.authentication.ad.ActiveDirectoryLdapAuthenticationProvider;
import org.springframework.security.web.util.matcher.AntPathRequestMatcher;

/**
 * Class for LDAP authentication.
 */

@Configuration
//@EnableGlobalMethodSecurity(securedEnabled = true)
@EnableWebSecurity
public class LdapSecurityConfig extends WebSecurityConfigurerAdapter {
	
	@Value("${agent.api.update}")
	private String pathUpdate;
	
	@Value("${user.normal.name}")
	private String normalUser;
	
	@Value("${user.normal.password}")
	private String normalPassword;
	
	@Value("${remember.me.seconds}")
	private String rememberMe;

	@Value("${auth.ad.enabled}")
	private String adEnabled;

	@Value("${auth.ad.domain}")
	private String adDomain;

	@Value("${auth.ad.url}")
	private String adUrl;

	@Value("${auth.ad.root}")
	private String adRoot;

	@Value("${auth.ad.search}")
	private String adSearch;

	@Autowired
	private LdapAuthoritiesMapper ldapAuthoritiesMapper;

	/**
	 * Authentication through LDAP without SALT crypt. 
	 */
	@Override
	protected void configure(final AuthenticationManagerBuilder auth) throws Exception {
		
		if (Boolean.parseBoolean(adEnabled)) {
			auth.authenticationProvider(activeDirectoryLdapAuthenticationProvider());
		} else {
	        auth
	        .inMemoryAuthentication()
	            .withUser(normalUser).password("{noop}" + normalPassword).roles("USER");
		}
	}

	/**
	 * Initializes the Auth provider.
	 * 
	 * @return the auth provider
	 */
	@Bean
    public AuthenticationProvider activeDirectoryLdapAuthenticationProvider() {
        ActiveDirectoryLdapAuthenticationProvider provider =
        		new ActiveDirectoryLdapAuthenticationProvider(adDomain, adUrl, adRoot);
        provider.setConvertSubErrorCodesToExceptions(true);
		provider.setUseAuthenticationRequestCredentials(true);
		provider.setAuthoritiesMapper(ldapAuthoritiesMapper);
		provider.setSearchFilter(adSearch);
        return provider;
	}

	/**
	 * Configuration of end point access
	 */
	@Override
	protected void configure(final HttpSecurity http) throws Exception {	
		// commenting the following row, breaks agent basic auth
		http.csrf().disable();
		
		http.authorizeRequests()
			.antMatchers("/assets/**").permitAll()
			.antMatchers("/packages/**").permitAll()
			.antMatchers(pathUpdate).permitAll()
			.antMatchers("/alerts/missing-host/**").permitAll()
			.anyRequest().fullyAuthenticated()
			.and()
		.formLogin()
			.loginPage("/login")
			.permitAll()
			.and()
		.logout()
			.logoutRequestMatcher(new AntPathRequestMatcher("/logout"))
			.logoutSuccessUrl("/login")
			.and()
		.rememberMe()
			.tokenValiditySeconds(Integer.parseInt(rememberMe));
	}

}
