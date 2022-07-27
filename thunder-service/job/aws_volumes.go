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

func (job *AwsDataRetrieveJob) FetchAwsVolumesNotUsed(profile model.AwsProfile, seqValue uint64) error {
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

	resultec2Svc, err := ec2Svc.DescribeVolumes(nil)
	if err != nil {
		return err
	}

	for _, w := range resultec2Svc.Volumes {
		if len(w.Attachments) == 0 {
			var objectName string

			for _, name := range w.Tags {
				if *name.Key == "Name" {
					objectName = *name.Value
					break
				}
			}

			recommendation.SeqValue = seqValue
			recommendation.ProfileID = profile.ID
			recommendation.Category = model.AwsUnusedResource
			recommendation.Suggestion = model.AwsDeleteBlockStorageNotUsed
			recommendation.Name = objectName
			recommendation.ResourceID = *w.VolumeId
			recommendation.ObjectType = model.AwsObjectVolume
			recommendation.Details = []map[string]interface{}{
				{"BLOCK_STORAGE_NAME": objectName},
				{"SIZE": *w.Size},
				{"ATTACCHED": "No"},
			}
			recommendation.CreatedAt = time.Now().UTC()

			listRec = append(listRec, recommendation)
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

func (job *AwsDataRetrieveJob) FetchAwsBlockStorageRightsizing(profile model.AwsProfile, seqValue uint64) error {
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

	resultec2Svc, err := ec2Svc.DescribeVolumes(nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -5)

	var volumeId string

	var iops, throughput int64

	for _, w := range resultec2Svc.Volumes {
		volumeId = *w.VolumeId

		if w.Iops != nil {
			iops = *w.Iops
		} else {
			iops = 0
		}

		if w.Throughput != nil {
			throughput = *w.Throughput
		} else {
			throughput = 0
		}

		iopsVolumeReadOps := GetIOPSthroughputMetricStatistics(sess, volumeId, "VolumeReadOps", "Count", timePast, timeNow)
		iopsVolumeWriteOps := GetIOPSthroughputMetricStatistics(sess, volumeId, "VolumeWriteOps", "Count", timePast, timeNow)
		throughputVolumeReadBytes := GetIOPSthroughputMetricStatistics(sess, volumeId, "VolumeReadBytes", "Bytes", timePast, timeNow)
		throughputiopVolumeWriteBytes := GetIOPSthroughputMetricStatistics(sess, volumeId, "VolumeWriteBytes", "Bytes", timePast, timeNow)
		maxIopsValue := GetMaximum(iopsVolumeReadOps, iopsVolumeWriteOps)
		maxThroughputValue := GetMaximum(throughputVolumeReadBytes, throughputiopVolumeWriteBytes)
		maxThroughputValue = maxThroughputValue / 1024 / 1024

		if iops < int64(maxIopsValue/2) && throughput < int64(maxThroughputValue/2) {
			var objectName string

			for _, name := range w.Tags {
				if *name.Key == "Name" {
					objectName = *name.Value
					break
				}
			}

			recommendation.SeqValue = seqValue
			recommendation.ProfileID = profile.ID
			recommendation.Category = model.AwsBlockStorageRightsizing
			recommendation.Suggestion = model.AwsResizeOversizedBlockStorage
			recommendation.Name = objectName
			recommendation.ResourceID = *w.VolumeId
			recommendation.ObjectType = model.AwsBlockStorage
			recommendation.Details = []map[string]interface{}{
				{"BLOCK_STORAGE_NAME": objectName},
				{"SIZE": *w.Size},
				{"TARGET": "THROUGHTPUT/IOPS"},
				{"THROUGHPUT_R/W_MAX_5DD ": maxThroughputValue},
				{"OPS_MAX_5DD": maxIopsValue},
			}

			recommendation.CreatedAt = time.Now().UTC()

			listRec = append(listRec, recommendation)
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

func GetIOPSthroughputMetricStatistics(sess *session.Session, volumeId string, metric string, unit string, startTime time.Time, endTime time.Time) *cloudwatch.GetMetricStatisticsOutput {
	svc := cloudwatch.New(sess)

	params := &cloudwatch.GetMetricStatisticsInput{
		EndTime:    aws.Time(endTime),
		MetricName: aws.String(metric),
		Namespace:  aws.String("AWS/EBS"),
		Period:     aws.Int64(432000),
		StartTime:  aws.Time(startTime),
		Statistics: []*string{
			aws.String("Maximum"),
		},
		Dimensions: []*cloudwatch.Dimension{
			{
				Name:  aws.String("VolumeId"),
				Value: aws.String(volumeId),
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

func GetMaximum(read *cloudwatch.GetMetricStatisticsOutput, write *cloudwatch.GetMetricStatisticsOutput) float64 {
	type ReadWrite struct {
		read  float64
		write float64
		sum   float64
	}

	var maxValue float64

	m := make(map[time.Time]ReadWrite)

	for _, r := range read.Datapoints {
		m[*r.Timestamp] = ReadWrite{read: *r.Maximum, sum: *r.Maximum}
	}

	for _, w := range write.Datapoints {
		if value, ok := m[*w.Timestamp]; ok {
			m[*w.Timestamp] = ReadWrite{read: value.read, write: *w.Maximum, sum: value.read + *w.Maximum}
		} else {
			m[*w.Timestamp] = ReadWrite{write: *w.Maximum, sum: *w.Maximum}
		}
	}

	for _, p := range m {
		if maxValue < p.sum {
			maxValue = p.sum
		}
	}

	return maxValue
}
