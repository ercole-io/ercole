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
					{"INSTANCE_TYPE": ""},
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
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	svc := ec2.New(sess)

	instances, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		return err
	}

	instanceTypes, err := svc.DescribeInstanceTypes(&ec2.DescribeInstanceTypesInput{})
	if err != nil {
		return err
	}

	optimizableInstances := make([]ec2.Instance, 0)
	client := cloudwatch.New(sess)

	for _, reservation := range instances.Reservations {
		for _, v := range reservation.Instances {
			c := make(chan ec2.Instance)

			instanceTypesLen := len(instanceTypes.InstanceTypes)

			go job.checkInstanceTypeCPU(instanceTypes.InstanceTypes[:instanceTypesLen/2], v, c)
			go job.checkInstanceTypeCPU(instanceTypes.InstanceTypes[instanceTypesLen/2:], v, c)

			input := &cloudwatch.GetMetricStatisticsInput{
				EndTime:    aws.Time(time.Unix(time.Now().Unix(), 0)),
				StartTime:  aws.Time(time.Unix(time.Now().Add(time.Duration(-168)*time.Hour).Unix(), 0)),
				MetricName: aws.String("CPUUtilization"),
				Period:     aws.Int64(3600),
				Namespace:  aws.String("AWS/EC2"),
				Statistics: aws.StringSlice([]string{"Average", "Maximum"}),
				Dimensions: []*cloudwatch.Dimension{{Name: aws.String("InstanceId"), Value: v.InstanceId}},
			}

			stats, err := client.GetMetricStatistics(input)
			if err != nil {
				return err
			}

			datapointsLen := len(stats.Datapoints)

			go job.getAverageStatistic(stats.Datapoints[:datapointsLen/2], v, c)
			go job.getAverageStatistic(stats.Datapoints[datapointsLen/2:], v, c)

			optimizableInstances = append(optimizableInstances, <-c)
		}
	}

	for _, instance := range optimizableInstances {
		var objectName string

		for _, name := range instance.Tags {
			if *name.Key == "Name" {
				objectName = *name.Value
				break
			}
		}

		awsRecommendation := model.AwsRecommendation{
			SeqValue:   seqValue,
			ProfileID:  profile.ID,
			Category:   model.AwsResizeComputeInstance,
			ObjectType: model.AwsComputeInstance,
			ResourceID: *instance.InstanceId,
			Name:       objectName,
			Suggestion: model.AwsResizeComputeInstance,
			CreatedAt:  time.Now(),
			Details: []map[string]interface{}{
				{"INSTANCE_NAME": objectName},
			},
		}

		if errDb := job.Database.AddAwsObject(awsRecommendation, database.AwsRecommendationCollection); errDb != nil {
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

func (job *AwsDataRetrieveJob) getAverageStatistic(datapoints []*cloudwatch.Datapoint, instance *ec2.Instance, c chan ec2.Instance) {
	for _, d := range datapoints {
		if *d.Average > 50 && *d.Maximum > 50 {
			c <- *instance
		}
	}
}

func (job *AwsDataRetrieveJob) checkInstanceTypeCPU(list []*ec2.InstanceTypeInfo, instance *ec2.Instance, c chan ec2.Instance) {
	for _, v := range list {
		if v.InstanceType == instance.InstanceType && v.VCpuInfo.DefaultVCpus == aws.Int64(1) {
			c <- *instance
		}
	}
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
