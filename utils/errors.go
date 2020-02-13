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

package utils

import (
	"errors"
)

// ErrHostNotFound contains "Host not found" error
var ErrHostNotFound = errors.New("Host not found")
var AerrHostNotFound AdvancedErrorInterface = NewAdvancedErrorPtr(ErrHostNotFound, "DB ERROR")
var ErrPatchingFunctionNotFound = errors.New("Patching Function not found")
var AerrPatchingFunctionNotFound AdvancedErrorInterface = NewAdvancedErrorPtr(ErrPatchingFunctionNotFound, "DB ERROR")

// ErrEventEnqueue contains "Failed to enqueue event" error
var ErrEventEnqueue = errors.New("Failed to enqueue event")

// ErrLicenseNotFound contains "License not found" error
var ErrLicenseNotFound = errors.New("License not found")
var AerrLicenseNotFound AdvancedErrorInterface = NewAdvancedErrorPtr(ErrLicenseNotFound, "DB ERROR")
