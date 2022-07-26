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

type AwsRecommendationDto struct {
	SeqValue   uint64                   `json:"seqValue"`
	Category   string                   `json:"category"`
	Suggestion string                   `json:"suggestion"`
	Name       string                   `json:"name"`
	ObjectType string                   `json:"objectType"`
	Details    []map[string]interface{} `json:"details"`
	CreatedAt  time.Time                `json:"createdAt"`
}

func ToAwsRecommendationDto(model model.AwsRecommendation) AwsRecommendationDto {
	return AwsRecommendationDto{
		SeqValue:   model.SeqValue,
		Category:   model.Category,
		Suggestion: model.Suggestion,
		Name:       model.Name,
		ObjectType: model.ObjectType,
		Details:    model.Details,
		CreatedAt:  model.CreatedAt,
	}
}

func ToAwsRecommendationsDto(models []model.AwsRecommendation) []AwsRecommendationDto {
	res := make([]AwsRecommendationDto, 0)

	for _, m := range models {
		res = append(res, ToAwsRecommendationDto(m))
	}

	return res
}
