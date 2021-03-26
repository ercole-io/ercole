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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/ercole-io/ercole/v2/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (as *APIService) AddMySQLAgreement(agreement model.MySQLAgreement) (primitive.ObjectID, error) {
	agreement.ID = as.NewObjectID()
	id, err := as.Database.AddMySQLAgreement(agreement)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return id, nil
}

func (as *APIService) UpdateMySQLAgreement(agreement model.MySQLAgreement) error {
	if err := as.Database.UpdateMySQLAgreement(agreement); err != nil {
		return err
	}

	return nil
}

func (as *APIService) GetMySQLAgreements() ([]model.MySQLAgreement, error) {
	agreements, err := as.Database.GetMySQLAgreements()
	if err != nil {
		return nil, err
	}
	return agreements, nil
}

func (as *APIService) DeleteMySQLAgreement(id primitive.ObjectID) error {
	if err := as.Database.DeleteMySQLAgreement(id); err != nil {
		return err
	}
	return nil
}
