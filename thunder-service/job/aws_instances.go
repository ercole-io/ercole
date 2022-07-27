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

	optimizableInstances := make([]ec2.Instance, 0)
	client := cloudwatch.New(sess)

	for _, reservation := range instances.Reservations {
		for _, v := range reservation.Instances {
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

			c := make(chan ec2.Instance)
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

func (job *AwsDataRetrieveJob) getAverageStatistic(datapoints []*cloudwatch.Datapoint, instance *ec2.Instance, c chan ec2.Instance) {
	for _, d := range datapoints {
		if *d.Average > 50 && *d.Maximum > 50 {
			c <- *instance
		}
	}
}
