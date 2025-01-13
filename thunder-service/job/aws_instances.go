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
package job

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatch/types"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
)

func (job *AwsDataRetrieveJob) FetchAwsNotActiveInstances(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(*cfg)

	resultec2Svc, err := ec2Client.DescribeInstances(context.Background(), nil)
	if err != nil {
		return err
	}

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			if i.State == nil {
				continue
			}
			if i.State.Name == ec2Types.InstanceStateNameStopped {
				var objectName string

				for _, name := range i.Tags {
					if *name.Key == "Name" {
						objectName = *name.Value
						break
					}
				}

				recommendation.SeqValue = seqValue
				recommendation.ProfileID = profile.ID
				recommendation.Category = model.AwsNotActiveResource
				recommendation.Suggestion = model.AwsDeleteComputeInstanceNotActive
				recommendation.Name = objectName
				recommendation.ResourceID = *i.InstanceId
				recommendation.ObjectType = model.AwsComputeInstance
				recommendation.Details = []map[string]interface{}{
					{"INSTANCE_NAME": objectName},
					{"INSTANCE_TYPE": i.InstanceType},
					{"STATUS": ec2Types.InstanceStateNameStopped},
				}
				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddAwsObjects(listRec, database.AwsRecommendationCollection)
		if errDb != nil {
			return errDb
		}
	}

	return nil
}

func (job *AwsDataRetrieveJob) FetchAwsComputeInstanceRightsizing(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(*cfg)

	resultec2Svc, err := ec2Client.DescribeInstances(context.Background(), nil)
	if err != nil {
		return err
	}

	instanceTypes, err := ec2Client.DescribeInstanceTypes(context.Background(), nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -7)

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			var average, maximum float64

			averageCPU := GetMetricStatistics(*cfg, cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(timeNow),
				MetricName: aws.String("CPUUtilization"),
				Namespace:  aws.String("AWS/EC2"),
				Period:     aws.Int32(3600),
				StartTime:  aws.Time(timePast),
				Statistics: []types.Statistic{types.StatisticAverage},
				Dimensions: []types.Dimension{{Name: aws.String("InstanceId"), Value: i.InstanceId}},
				Unit:       types.StandardUnitPercent,
			})
			countAverageCPU := 0

			for _, op := range averageCPU.Datapoints {
				average += *op.Average

				if *op.Average > float64(50) {
					countAverageCPU += 1
				}
			}

			maxCPU := GetMetricStatistics(*cfg, cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(timeNow),
				MetricName: aws.String("CPUUtilization"),
				Namespace:  aws.String("AWS/EC2"),
				Period:     aws.Int32(3600),
				StartTime:  aws.Time(timePast),
				Statistics: []types.Statistic{types.StatisticMaximum},
				Dimensions: []types.Dimension{{Name: aws.String("InstanceId"), Value: i.InstanceId}},
				Unit:       types.StandardUnitPercent,
			})
			countMaxCPU := 0

			for _, op := range maxCPU.Datapoints {
				maximum += *op.Maximum

				if *op.Maximum > float64(50) {
					countMaxCPU += 1
				}
			}

			isKoCPU := false

			for _, v := range instanceTypes.InstanceTypes {
				if v.InstanceType == i.InstanceType && *v.VCpuInfo.DefaultVCpus == 1 {
					isKoCPU = true
				}

				break
			}

			if countAverageCPU <= 3 && countMaxCPU <= 180 && !isKoCPU {
				var objectName string

				for _, name := range i.Tags {
					if *name.Key == "Name" {
						objectName = *name.Value
						break
					}
				}

				lenDatapoints := len(maxCPU.Datapoints)

				if lenDatapoints > 0 {
					average = average / float64(lenDatapoints)
					maximum = maximum / float64(lenDatapoints)
				} else {
					average = 0
					maximum = 0
				}

				recommendation.SeqValue = seqValue
				recommendation.ProfileID = profile.ID
				recommendation.Category = model.AwsComputeInstance
				recommendation.Suggestion = model.AwsResizeComputeInstance
				recommendation.Name = objectName
				recommendation.ResourceID = *i.InstanceId
				recommendation.ObjectType = model.AwsResizeComputeInstance
				recommendation.Details = []map[string]interface{}{
					{"INSTANCE_NAME": objectName},
					{"SHAPE": i.InstanceType},
					{"%_CPU_AVERAGE_7DD(DAILY)": average},
					{"NUMBER_OF_THRESHOLD_REACHED_(>50%)": "3"},
					{"%_CPU_AVERAGE_7DD(MINUTES)": maximum},
					{"NUMBER_OF_THRESHOLD_REACHED_(>50%)": "3"},
				}

				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddAwsObjects(listRec, database.AwsRecommendationCollection)
		if errDb != nil {
			return errDb
		}
	}

	return nil
}

func (job *AwsDataRetrieveJob) FetchAwsInstanceDecommissioning2(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}

	ec2Client := ec2.NewFromConfig(*cfg)

	resultec2Svc, err := ec2Client.DescribeInstances(context.Background(), nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -7)

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			var average, maximum float64

			averageCPU := GetMetricStatistics(*cfg, cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(timeNow),
				MetricName: aws.String("CPUUtilization"),
				Namespace:  aws.String("AWS/EC2"),
				Period:     aws.Int32(86400),
				StartTime:  aws.Time(timePast),
				Statistics: []types.Statistic{types.StatisticAverage},
				Dimensions: []types.Dimension{{Name: aws.String("InstanceId"), Value: i.InstanceId}},
				Unit:       types.StandardUnitPercent,
			})
			countAverageCPU := 0

			for _, op := range averageCPU.Datapoints {
				average += *op.Average

				if *op.Average > float64(5) {
					countAverageCPU += 1
				}
			}

			maxCPU := GetMetricStatistics(*cfg, cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(timeNow),
				MetricName: aws.String("CPUUtilization"),
				Namespace:  aws.String("AWS/EC2"),
				Period:     aws.Int32(3600),
				StartTime:  aws.Time(timePast),
				Statistics: []types.Statistic{types.StatisticMaximum},
				Dimensions: []types.Dimension{{Name: aws.String("InstanceId"), Value: i.InstanceId}},
				Unit:       types.StandardUnitPercent,
			})
			countMaxCPU := 0

			for _, op := range maxCPU.Datapoints {
				maximum += *op.Maximum

				if *op.Maximum > float64(5) {
					countMaxCPU += 1
				}
			}

			if countAverageCPU <= 3 && countMaxCPU <= 180 {
				var objectName string

				for _, name := range i.Tags {
					if *name.Key == "Name" {
						objectName = *name.Value
						break
					}
				}

				lenDatapoints := len(maxCPU.Datapoints)

				if lenDatapoints > 0 {
					average = average / float64(lenDatapoints)
					maximum = maximum / float64(lenDatapoints)
				} else {
					average = 0
					maximum = 0
				}

				recommendation.SeqValue = seqValue
				recommendation.ProfileID = profile.ID
				recommendation.Category = model.AwsComputeInstance
				recommendation.Suggestion = model.AwsDeleteComputeInstanceNotUsed
				recommendation.Name = objectName
				recommendation.ResourceID = *i.InstanceId
				recommendation.ObjectType = model.AwsComputeInstanceNotUsed
				recommendation.Details = []map[string]interface{}{
					{"INSTANCE_NAME": objectName},
					{"SHAPE": i.InstanceType},
					{"%_CPU_AVERAGE_7DD(DAILY)": average},
					{"NUMBER_OF_THRESHOLD_REACHED_(>5%)": "3"},
					{"%_CPU_AVERAGE_7DD(MINUTES)": maximum},
					{"NUMBER_OF_THRESHOLD_REACHED_(>5%)": "180"},
				}

				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddAwsObjects(listRec, database.AwsRecommendationCollection)
		if errDb != nil {
			return errDb
		}
	}

	return nil
}
