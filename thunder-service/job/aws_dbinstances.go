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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
)

func (job *AwsDataRetrieveJob) FetchAwsUnusedDatabaseInstance(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	rdsSvc := rds.New(sess)

	resultrdsSvc, err := rdsSvc.DescribeDBInstances(nil)
	if err != nil {
		return err
	}

	for _, w := range resultrdsSvc.DBInstances {
		if *w.DBInstanceStatus == "inaccessible-encryption-credentials-recoverable" ||
			*w.DBInstanceStatus == "incompatible-option-group" ||
			*w.DBInstanceStatus == "incompatible-parameters" ||
			*w.DBInstanceStatus == "restore-error" {
			recommendation.SeqValue = seqValue
			recommendation.ProfileID = profile.ID
			recommendation.Category = model.AwsUnusedResource
			recommendation.Suggestion = model.AwsDeleteDatabaseInstanceNotActive
			recommendation.Name = *w.DBInstanceIdentifier
			recommendation.ResourceID = *w.DBInstanceIdentifier
			recommendation.ObjectType = model.AwsDatabaseInstance
			recommendation.Details = []map[string]interface{}{
				{"INSTANCE_NAME": *w.DBInstanceArn},
				{"DATABASE_NAME": *w.DBInstanceIdentifier},
				{"DB_INSTANCE_STATUS": *w.DBInstanceStatus},
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

func (job *AwsDataRetrieveJob) FetchAwsUnusedServiceDecommissioning3DB(profile model.AwsProfile, seqValue uint64) error {
	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)

	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	rdsSvc := rds.New(sess)

	resultrdsSvc, err := rdsSvc.DescribeDBInstances(nil)
	if err != nil {
		return err
	}

	ec2Svc := ec2.New(sess)

	instanceTypes, err := ec2Svc.DescribeInstanceTypes(nil)
	if err != nil {
		return err
	}

	timeNow := time.Now()
	timePast := timeNow.AddDate(0, 0, -7)

	for _, w := range resultrdsSvc.DBInstances {
		shape := strings.TrimLeft(*w.DBInstanceClass, "db.")

		var average, maximum float64

		averageCPU := GetMetricStatistics(sess, "DBInstanceIdentifier", *w.DBInstanceIdentifier, "CPUUtilization", "AWS/RDS", 3600, "Average", "Percent", timePast, timeNow)
		countAverageCPU := 0

		for _, op := range averageCPU.Datapoints {
			average += *op.Average

			if *op.Average > float64(50) {
				countAverageCPU += 1
			}
		}

		maxCPU := GetMetricStatistics(sess, "DBInstanceIdentifier", *w.DBInstanceIdentifier, "CPUUtilization", "AWS/RDS", 3600, "Maximum", "Percent", timePast, timeNow)
		countMaxCPU := 0

		for _, op := range maxCPU.Datapoints {
			maximum += *op.Maximum

			if *op.Maximum > float64(50) {
				countMaxCPU += 1
			}
		}

		isKoCPU := false
		for _, v := range instanceTypes.InstanceTypes {

			if *v.InstanceType == shape && *v.VCpuInfo.DefaultVCpus == 1 {
				isKoCPU = true
			}

			break
		}

		if countAverageCPU <= 3 && countMaxCPU <= 3 && !isKoCPU {

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
			recommendation.Category = model.AwsOversizedDatabase
			recommendation.Suggestion = model.AwsResizeOversizedDatabaseInstance
			recommendation.Name = *w.DBInstanceIdentifier
			recommendation.ResourceID = *w.DBInstanceIdentifier
			recommendation.ObjectType = model.AwsDatabaseInstance
			recommendation.Details = []map[string]interface{}{
				{"INSTANCE_NAME": *w.DBInstanceIdentifier},
				{"SHAPE": shape},
				{"%_CPU_AVERAGE_7DD(DAILY)": average},
				{"NUMBER_OF_THRESHOLD_REACHED_(>50%)": "3"},
				{"%_CPU_AVERAGE_7DD(MINUTES)": maximum},
				{"NUMBER_OF_THRESHOLD_REACHED_(>50%)": "3"},
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
