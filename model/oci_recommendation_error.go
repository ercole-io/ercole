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

package model

import "time"

type OciRecommendationError struct {
	SeqValue  uint64    `json:"seqValue" bson:"seqValue"`
	ProfileID string    `json:"profileID" bson:"profileID"`
	Category  string    `json:"category" bson:"category"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Error     string    `json:"error" bson:"error"`
}

func (ore OciRecommendationError) SetOciRecommendationError(seqValue uint64, profileID string, category string, createdAt time.Time, strError string) OciRecommendationError {
	recError := OciRecommendationError{
		SeqValue:  seqValue,
		ProfileID: profileID,
		Category:  category,
		CreatedAt: createdAt,
		Error:     strError,
	}

	return recError
}
