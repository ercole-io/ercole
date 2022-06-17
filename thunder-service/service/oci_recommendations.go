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
	"time"

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/job"
)

func (as *ThunderService) GetOciRecommendations() ([]model.OciRecommendation, error) {
	selectedProfiles, err := as.Database.GetSelectedOciProfiles()
	if err != nil {
		return nil, err
	}

	ociRecommendations := make([]model.OciRecommendation, 0)
	if len(selectedProfiles) > 0 {
		ociRecommendations, err = as.Database.GetOciRecommendations(selectedProfiles)

		if err != nil {
			return nil, err
		}
	}

	return ociRecommendations, err
}

func (as *ThunderService) ForceGetOciRecommendations() error {
	log := logger.NewLogger("THUN", logger.LogVerbosely(true))

	j := &job.OciDataRetrieveJob{
		Database: as.Database,
		TimeNow:  time.Now,
		Config:   as.Config,
		Log:      log,
	}
	j.Run()

	return nil
}
