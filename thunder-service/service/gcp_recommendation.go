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
	"github.com/ercole-io/ercole/v2/thunder-service/job"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ts *ThunderService) ListGcpRecommendations() ([]model.GcpRecommendation, error) {
	selectedProfiles, err := ts.Database.GetActiveGcpProfiles()
	if err != nil {
		return nil, err
	}

	profileIDs := make([]primitive.ObjectID, 0, len(selectedProfiles))
	for _, p := range selectedProfiles {
		profileIDs = append(profileIDs, p.ID)
	}

	return ts.Database.ListGcpRecommendationsByProfiles(profileIDs)
}

func (ts *ThunderService) ForceGetGcpRecommendations() {
	j := &job.GcpDataRetrieveJob{
		Database: ts.Database,
		Config:   ts.Config,
		Log:      ts.Log,
		Opt:      nil,
	}

	j.Run()
}
