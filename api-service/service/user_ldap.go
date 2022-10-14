// Copyright (c) 2022 Sorint.lab S.p.A.
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

package service

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/go-ldap/ldap"
)

func (as *APIService) GetLDAPUsers(user string) ([]model.UserLDAP, error) {
	ldapBase := as.Config.APIService.AuthenticationProvider.LDAPBase

	var filter string

	var resultSR *ldap.SearchRequest

	var sr *ldap.SearchResult

	var errSearch error

	l, err := ldap.DialURL(fmt.Sprintf("%s%s%s%s", "ldap://", as.Config.APIService.AuthenticationProvider.Host, ":", strconv.Itoa(as.Config.APIService.AuthenticationProvider.Port)))
	if err != nil {
		return nil, utils.NewError(utils.ErrConnectLDAPServer, "LDAP ERROR")
	}

	err = l.StartTLS(&tls.Config{InsecureSkipVerify: true})
	if err != nil {
		return nil, err
	}

	err = l.Bind(as.Config.APIService.AuthenticationProvider.LDAPBindDN, as.Config.APIService.AuthenticationProvider.LDAPBindPassword)
	if err != nil {
		return nil, err
	}

	filter = fmt.Sprintf("(uid=%s)", user)

	resultSR = ldapSearchRequest(filter, ldapBase)

	sr, errSearch = l.Search(resultSR)
	if errSearch != nil || len(sr.Entries) == 0 {
		if len(sr.Entries) != 1 {
			filter = fmt.Sprintf("(sn=%s)", user)

			resultSR = ldapSearchRequest(filter, ldapBase)

			sr, errSearch = l.Search(resultSR)
			if errSearch != nil || len(sr.Entries) == 0 {
				if len(sr.Entries) != 1 {
					filter = fmt.Sprintf("(givenName=%s)", user)

					resultSR = ldapSearchRequest(filter, ldapBase)

					sr, errSearch = l.Search(resultSR)
					if errSearch != nil || len(sr.Entries) == 0 {
						if len(sr.Entries) != 1 {
							return nil, errors.New("User does not exist")
						} else {
							return nil, errSearch
						}
					}
				} else {
					return nil, errSearch
				}
			}
		} else {
			return nil, errSearch
		}
	}

	results := make([]model.UserLDAP, 0)

	for _, entry := range sr.Entries {
		user := new(model.UserLDAP)

		for _, at := range entry.Attributes {
			switch at.Name {
			case "givenName":
				user.GivenName = at.Values[0]
			case "sn":
				user.Sn = at.Values[0]
			case "uid":
				user.Uid = at.Values[0]
			}
		}

		results = append(results, *user)
	}

	return results, nil
}

func ldapSearchRequest(filter string, ldapBase string) *ldap.SearchRequest {
	searchRequest := ldap.NewSearchRequest(
		ldapBase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=*)%s)", filter),
		[]string{"givenName", "sn", "uid"},
		nil,
	)

	return searchRequest
}

func (as *APIService) AddUserByLDAP(userLDAP model.UserLDAP, groups []string) error {
	return as.AddUser(model.User{
		Username:  userLDAP.Uid,
		FirstName: userLDAP.GivenName,
		LastName:  userLDAP.Sn,
		Groups:    groups,
	})
}
