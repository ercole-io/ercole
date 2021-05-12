// Copyright (c) 2020 Sorint.lab S.p.A.
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
var ErrPatchingFunctionNotFound = errors.New("Patching Function not found")

// ErrEventEnqueue contains "Failed to enqueue event" error
var ErrEventEnqueue = errors.New("Failed to enqueue event")

// ErrLicenseNotFound contains "License not found" error
var ErrLicenseNotFound = errors.New("License not found")

// ErrAlertNotFound contains "Alert not found" error
var ErrAlertNotFound = errors.New("Alert not found")

// ErrClusterNotFound contains "Cluster not found" error
var ErrClusterNotFound = errors.New("Cluster not found")

var ErrOracleDatabaseLicenseTypeIDNotFound = errors.New("Oracle Database LicenseTypeID not found")

// ErrOracleDatabaseAgreementNotFound contains "Agreement not found" error
var ErrOracleDatabaseAgreementNotFound = errors.New("Agreement not found")

// ErrNotInClusterHostNotFound contains "Baremetal host not found" error
var ErrNotInClusterHostNotFound = errors.New("Not in cluster host not found")

var ErrInvalidHostdata = errors.New("Invalid hostdata")
