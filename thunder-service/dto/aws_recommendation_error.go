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
