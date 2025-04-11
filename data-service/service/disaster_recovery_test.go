// Copyright (c) 2025 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/require"
	gomock "go.uber.org/mock/gomock"
)

func TestCreateDR(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	hds := HostDataService{
		Database: db,
	}

	drname := "test_DR"

	db.EXPECT().ExistsDR(drname).Return(true).Times(1)
	db.EXPECT().DismissHost(drname).Return(nil).Times(1)
	db.EXPECT().InsertHostData(gomock.Any()).Return(nil).Times(1)

	err := hds.createDR(model.HostDataBE{
		Hostname: "test",
		Archived: false,
		IsDR:     true,
	})

	require.NoError(t, err)
}
