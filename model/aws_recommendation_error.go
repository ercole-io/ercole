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

package model

import "time"

type AwsRecommendationError struct {
	SeqValue  uint64    `json:"seqValue" bson:"seqValue"`
	ProfileID string    `json:"profileID" bson:"profileID"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Error     string    `json:"error" bson:"error"`
}

func (are AwsRecommendationError) SetAwsRecommendationError(seqValue uint64, profileID string, createdAt time.Time, strError string) AwsRecommendationError {
	recError := AwsRecommendationError{
		SeqValue:  seqValue,
		ProfileID: profileID,
		CreatedAt: createdAt,
		Error:     strError,
	}

	return recError
}
