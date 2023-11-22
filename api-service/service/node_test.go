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
	"testing"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestGetNodes_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	group := model.Group{
		Name:  model.GroupAdmin,
		Roles: []string{model.AdminPermission},
	}

	expected := []model.Node{
		{
			Name:   "test",
			Roles:  []string{"admin"},
			Parent: "",
		},
	}

	db.EXPECT().GetNodesByRoles([]string{model.AdminPermission}).Return(expected, nil)
	db.EXPECT().GetGroup("admin").Return(&group, nil)

	res, err := as.GetNodes([]string{"admin"})
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestGetNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := &model.Node{
		Name:   "test",
		Roles:  []string{"admin"},
		Parent: "",
	}

	db.EXPECT().GetNodeByName("test").Return(expected, nil)

	res, err := as.GetNode("test")
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}

func TestAddNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := model.Node{
		Name:   "test",
		Roles:  []string{"admin"},
		Parent: "",
	}

	db.EXPECT().AddNode(expected).Return(nil)

	err := as.AddNode(expected)
	require.NoError(t, err)
}

func TestUpdateNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	expected := model.Node{
		Name:   "test",
		Roles:  []string{"admin"},
		Parent: "",
	}

	db.EXPECT().UpdateNode(expected).Return(nil)

	err := as.UpdateNode(expected)
	require.NoError(t, err)
}

func TestRemoveNode_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	db.EXPECT().RemoveNode("test").Return(nil)

	err := as.RemoveNode("test")
	require.NoError(t, err)
}
