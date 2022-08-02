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
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatch"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
)

func (job *AwsDataRetrieveJob) FetchAwsNotActiveInstances(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	ec2Svc := ec2.New(sess)

	resultec2Svc, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return err
	}

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			if *i.State.Name == "stopped" {
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
					{"INSTANCE_TYPE": *i.InstanceType},
					{"STATUS": "stopped"},
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

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	ec2Svc := ec2.New(sess)

	resultec2Svc, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return err
	}

	instanceTypes, err := ec2Svc.DescribeInstanceTypes(nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -7)

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			var average, maximum float64

			averageCPU := getCPUMetricStatistics(sess, *i.InstanceId, "CPUUtilization", 3600, "Average", "Percent", timePast, timeNow)
			countAverageCPU := 0

			for _, op := range averageCPU.Datapoints {
				average += *op.Average

				if *op.Average > float64(50) {
					countAverageCPU += 1
				}
			}

			maxCPU := getCPUMetricStatistics(sess, *i.InstanceId, "CPUUtilization", 3600, "Maximum", "Percent", timePast, timeNow)
			countMaxCPU := 0

			for _, op := range maxCPU.Datapoints {
				maximum += *op.Maximum

				if *op.Maximum > float64(50) {
					countMaxCPU += 1
				}
			}

			isKoCPU := false

			for _, v := range instanceTypes.InstanceTypes {
				if *v.InstanceType == *i.InstanceType && *v.VCpuInfo.DefaultVCpus == 1 {
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
					{"SHAPE": *i.InstanceType},
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

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	ec2Svc := ec2.New(sess)

	resultec2Svc, err := ec2Svc.DescribeInstances(nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -7)

	for _, w := range resultec2Svc.Reservations {
		for _, i := range w.Instances {
			var average, maximum float64

			averageCPU := getCPUMetricStatistics(sess, *i.InstanceId, "CPUUtilization", 86400, "Average", "Percent", timePast, timeNow)
			countAverageCPU := 0

			for _, op := range averageCPU.Datapoints {
				average += *op.Average

				if *op.Average > float64(5) {
					countAverageCPU += 1
				}
			}

			maxCPU := getCPUMetricStatistics(sess, *i.InstanceId, "CPUUtilization", 3600, "Maximum", "Percent", timePast, timeNow)
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
					{"SHAPE": *i.InstanceType},
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

func getCPUMetricStatistics(sess *session.Session, instanceID string, metric string, period int64, statistics string, unit string, startTime time.Time, endTime time.Time) *cloudwatch.GetMetricStatisticsOutput {
	svc := cloudwatch.New(sess)

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(endTime),
		MetricName: aws.String(metric),
		Namespace:  aws.String("AWS/EC2"),
		Period:     aws.Int64(period),
		StartTime:  aws.Time(startTime),
		Statistics: []*string{
			aws.String(statistics),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("InstanceId"),
				Value: aws.String(instanceID),
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
