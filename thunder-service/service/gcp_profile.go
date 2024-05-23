// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ts *ThunderService) AddGcpProfile(profile dto.GcpProfileRequest) error {
	profileModel := model.GcpProfile{
		ID:          ts.NewObjectID(),
		Name:        profile.Name,
		Selected:    false,
		PrivateKey:  profile.PrivateKey,
		ClientEmail: profile.ClientEmail,
	}

	err := ts.Database.AddGcpProfile(profileModel)
	if err != nil {
		return err
	}

	return nil
}

func (ts *ThunderService) GetGcpProfiles() ([]model.GcpProfile, error) {
	return ts.Database.ListGcpProfiles()
}

func (ts *ThunderService) SelectGcpProfile(idhex string, selected bool) error {
	id, err := primitive.ObjectIDFromHex(idhex)
	if err != nil {
		return err
	}

	if err := ts.Database.SelectGcpProfile(id, selected); err != nil {
		return err
	}

	return nil
}

func (ts *ThunderService) UpdateGcpProfile(profileID string, profile dto.GcpProfileRequest) error {
	id, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	return ts.Database.UpdateGcpProfile(id, model.GcpProfile{
		ID:          id,
		Name:        profile.Name,
		PrivateKey:  profile.PrivateKey,
		ClientEmail: profile.ClientEmail,
	})
}

func (ts *ThunderService) RemoveGcpProfile(profileID string) error {
	id, err := primitive.ObjectIDFromHex(profileID)
	if err != nil {
		return err
	}

	return ts.Database.RemoveGcpProfile(id)
}
