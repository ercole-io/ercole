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

package service

import (
	"testing"

	"github.com/ercole-io/ercole/config"
	gomock "github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestInit_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config: config.Configuration{
			ResourceFilePath: "../../resources",
		},
	}
	as.Init()

	assert.Len(t, as.TechnologyInfos, 14)
	assert.Equal(t, "Oracle/Database", as.TechnologyInfos[0].Product)
	assert.Equal(t, "Microsoft/SQLServer", as.TechnologyInfos[1].Product)
	assert.Equal(t, "iVBORw0K", as.TechnologyInfos[0].Logo[:8])
	assert.Equal(t, "iVBORw0K", as.TechnologyInfos[1].Logo[:8])
}
