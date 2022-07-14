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
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *AwsDataRetrieveJob) RetrieveObjectStorageOptimization() error {
	awsProfiles, err := job.Database.GetAwsProfiles(false)
	if err != nil {
		return err
	}

	c := make(chan error)

	for _, profile := range awsProfiles {
		go func(profile model.AwsProfile) {
			if err := job.RetrieveBuckets(profile); err != nil {
				c <- err
			}
		}(profile)
	}

	return <-c
}

func (job *AwsDataRetrieveJob) RetrieveBuckets(profile model.AwsProfile) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	svc := s3.New(sess)
	input := &s3.ListBucketsInput{}

	result, err := svc.ListBuckets(input)
	if err != nil {
		return err
	}

	c := make(chan error)

	for _, v := range result.Buckets {
		go func(name string, accessKeyId string) {
			if err := job.RetrieveBucketLifecycleConfiguration(name, accessKeyId, svc); err != nil {
				c <- err
			}
		}(*v.Name, profile.AccessKeyId)
	}

	return <-c
}

func (job *AwsDataRetrieveJob) RetrieveBucketLifecycleConfiguration(bucketName string, accessKeyId string, svc *s3.S3) error {
	inputBucket := s3.GetBucketLifecycleConfigurationInput{
		Bucket: aws.String(bucketName),
	}

	_, err := svc.GetBucketLifecycleConfiguration(&inputBucket)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "NoSuchLifecycleConfiguration":
				bucketObject, errObj := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucketName)})
				if errObj != nil {
					return errObj
				}

				objLen := len(bucketObject.Contents)

				c := make(chan int64)
				go job.sumSize(bucketObject.Contents[:objLen/2], c)
				go job.sumSize(bucketObject.Contents[objLen/2:], c)
				objSum := <-c + <-c

				awsRecommendation := model.AwsRecommendation{
					ProfileID: accessKeyId,
					Details: []map[string]interface{}{
						{"OPTIMIZATION": "ENABLE AUTO-TIERING"},
						{"BUCKET_NAME": bucketName},
						{"OBJECTS": objLen},
						{"SIZE": objSum},
					},
					CreatedAt: time.Now(),
				}
				if errDb := job.Database.AddAwsObject(awsRecommendation, "aws_recommendations"); errDb != nil {
					job.Log.Error(errDb)
					return errDb
				}
			default:
				job.Log.Warn(err)
			}
		} else {
			job.Log.Warn(err)
		}
	}

	return nil
}

func (job *AwsDataRetrieveJob) sumSize(contents []*s3.Object, c chan int64) {
	var sum int64

	for _, cnt := range contents {
		sum += int64(*cnt.Size)
	}

	c <- sum
}
