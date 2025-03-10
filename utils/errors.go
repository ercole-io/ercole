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

package utils

import (
	"errors"
)

// ErrNotFound generic not found
var ErrNotFound = errors.New("Not found")

// ErrHostNotFound contains "Host not found" error
var ErrHostNotFound = errors.New("Host not found")

var ErrUserNotFound = errors.New("User not found")

// ErrEventEnqueue contains "Failed to enqueue event" error
var ErrEventEnqueue = errors.New("Failed to enqueue event")

// ErrLicenseNotFound contains "License not found" error
var ErrLicenseNotFound = errors.New("License not found")

// ErrAlertNotFound contains "Alert not found" error
var ErrAlertNotFound = errors.New("Alert not found")

// ErrClusterNotFound contains "Cluster not found" error
var ErrClusterNotFound = errors.New("Cluster not found")

// ErrOracleDatabaseLicenseTypeIDNotFound contains "LicenseTypeID Not Found" error
var ErrOracleDatabaseLicenseTypeIDNotFound = errors.New("Oracle Database LicenseTypeID not found")

// ErrContractNotFound contains "Contract not found" error
var ErrContractNotFound = errors.New("Contract not found")

// ErrNotInClusterHostNotFound contains "Baremetal host not found" error
var ErrNotInClusterHostNotFound = errors.New("Not in cluster host not found")

var ErrInvalidHostdata = errors.New("Invalid hostdata")

var ErrInvalidJSON = errors.New("invalid JSON")

var ErrInvalidLocation = errors.New("Invalid location")

var ErrInvalidLicenseType = errors.New("Invalid license type")

var ErrInvalidRole = errors.New("Invalid role")

var ErrInvalidGroup = errors.New("Invalid group")

var ErrInvalidUser = errors.New("Invalid user")

var ErrInvalidAck = errors.New("Alert(s) cannot be acknowledged")

var ErrInvalidToken = errors.New("invalid token")

var ErrHostNotInCluster = errors.New("host not in cluster")

var ErrInvalidProfileId = errors.New("invalid profile id")

var ErrConnDB = errors.New("Can't connect to the database!")

var ErrRoleNotFound = errors.New("Role not found")

var ErrRoleAlreadyExists = errors.New("Role already exists")

var ErrGroupNotFound = errors.New("Group not found")

var ErrGroupAlreadyExists = errors.New("Group already exists")

var ErrGroupCannotBeDeleted = errors.New("The group cannot be deleted because it is associated with one or more users")

var ErrRoleCannotBeDeleted = errors.New("The role cannot be deleted because it is associated with one or more groups")

var ErrConnectLDAPServer = errors.New("Cannot connect to LDAP server")

var ErrSuperUserCannotBeDeleted = errors.New("Super User cannot be deleted")

var ErrPermissionDenied = errors.New("Permission denied")

var ErrInvalidExadata = errors.New("invalid exadata")

var ErrInvalidOracleContract = errors.New("invalid oracle contract")

var ErrMissingDatabaseNotFound = errors.New("Missing database not found")
