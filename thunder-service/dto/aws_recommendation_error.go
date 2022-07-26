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

package dto

import (
	"time"

	"github.com/ercole-io/ercole/v2/model"
)

type AwsRecommendationErrorDto struct {
	SeqValue  uint64    `json:"seqValue"`
	Category  string    `json:"category"`
	Error     string    `json:"error"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToAwsRecommendationErrorDto(model model.AwsRecommendation) []AwsRecommendationErrorDto {
	res := make([]AwsRecommendationErrorDto, 0, len(model.Errors))

	for _, e := range model.Errors {
		for _, v := range e {
			dto := AwsRecommendationErrorDto{
				SeqValue:  model.SeqValue,
				Category:  model.Category,
				CreatedAt: model.CreatedAt,
				Error:     v,
			}

			res = append(res, dto)
		}
	}

	return res
}

func ToAwsRecommendationsErrorsDto(models []model.AwsRecommendation) []AwsRecommendationErrorDto {
	res := make([]AwsRecommendationErrorDto, 0)

	for _, m := range models {
		res = append(res, ToAwsRecommendationErrorDto(m)...)
	}

	return res
}
