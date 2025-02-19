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

package job

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
)

func GetMetricStatistics(sess *session.Session, dimensionName string, dimensionValue string, metric string, namespace string, period int64, statistics string, unit string, startTime time.Time, endTime time.Time) *cloudwatch.GetMetricStatisticsOutput {
	svc := cloudwatch.New(sess)

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(endTime),
		MetricName: aws.String(metric),
		Namespace:  aws.String(namespace),
		Period:     aws.Int64(period),
		StartTime:  aws.Time(startTime),
		Statistics: []*string{
			aws.String(statistics),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String(dimensionName),
				Value: aws.String(dimensionValue),
			},
		},
		Unit: aws.String(unit),
	}
	resp, err := svc.GetMetricStatistics(params)

	if err != nil {
		return nil
	}

	return resp
}
