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
	"github.com/ercole-io/ercole/v2/model"
)

func (ts *ThunderService) GetOciRecommendationErrors(seqNum uint64) ([]model.OciRecommendationError, error) {
	if seqNum == 0 {
		seqValue, err := ts.Database.GetLastOciRecommendationErrorsSeqValue()
		if err != nil {
			return nil, err
		}

		seqNum = seqValue
	}

	ociRecommendationErrors, err := ts.Database.GetOciRecommendationErrors(seqNum)
	if err != nil {
		return nil, err
	}

	return ociRecommendationErrors, nil
}
