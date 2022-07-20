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

// Package service is a package that provides methods for querying data
package service

import (
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/job"
)

func (as *ThunderService) GetAwsRecommendations() ([]model.AwsRecommendation, error) {
	selectedProfiles, err := as.Database.GetSelectedAwsProfiles()
	if err != nil {
		return nil, err
	}

	awsRecommendations := make([]model.AwsRecommendation, 0)
	if len(selectedProfiles) > 0 {
		awsRecommendations, err = as.Database.GetAwsRecommendations(selectedProfiles)

		if err != nil {
			return nil, err
		}
	}

	return awsRecommendations, err
}

func (as *ThunderService) ForceGetAwsRecommendations() error {
	log := logger.NewLogger("THUN", logger.LogVerbosely(true))

	j := &job.AwsDataRetrieveJob{
		Database: as.Database,
		Config:   as.Config,
		Log:      log,
	}
	j.Run()

	return nil
}
