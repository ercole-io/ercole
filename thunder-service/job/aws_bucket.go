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
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	smithy "github.com/aws/smithy-go"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (job *AwsDataRetrieveJob) FetchObjectStorageOptimization(profile model.AwsProfile, seq uint64) error {
	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}

	s3client := s3.NewFromConfig(*cfg)
	input := &s3.ListBucketsInput{}

	result, err := s3client.ListBuckets(context.Background(), input)
	if err != nil {
		return err
	}

	c := make(chan error)

	for _, v := range result.Buckets {
		go func(name string, profileID primitive.ObjectID, seq uint64) {
			if err := job.FetchBucketLifecycleConfiguration(name, profile.ID, seq, s3client); err != nil {
				c <- err
			}
		}(*v.Name, profile.ID, seq)
	}

	return <-c
}

func (job *AwsDataRetrieveJob) FetchBucketLifecycleConfiguration(bucketName string, profileID primitive.ObjectID, seqValue uint64, s3client *s3.Client) error {
	inputBucket := s3.GetBucketLifecycleConfigurationInput{
		Bucket: &bucketName,
	}

	_, err := s3client.GetBucketLifecycleConfiguration(context.Background(), &inputBucket)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			switch ae.ErrorCode() {
			case "NoSuchLifecycleConfiguration":
				bucketObject, errObj := s3client.ListObjects(context.Background(), &s3.ListObjectsInput{Bucket: &bucketName})
				if errObj != nil {
					return errObj
				}

				objLen := len(bucketObject.Contents)

				c := make(chan int64)
				go job.sumSize(bucketObject.Contents[:objLen/2], c)
				go job.sumSize(bucketObject.Contents[objLen/2:], c)
				objSum := <-c + <-c

				if objSum < 500_000_000_000 {
					return nil
				}

				awsRecommendation := model.AwsRecommendation{
					SeqValue:   seqValue,
					ProfileID:  profileID,
					Category:   model.AwsObjectStorageOptimization,
					Suggestion: model.AwsObjectStorageOptimizationSuggestion,
					ObjectType: model.AwsObjectStorageOptimizationType,
					Name:       bucketName,
					Details: []map[string]interface{}{
						{"OPTIMIZATION": "ENABLE AUTO-TIERING"},
						{"BUCKET_NAME": bucketName},
						{"OBJECTS": objLen},
						{"SIZE": objSum},
					},
					CreatedAt: time.Now(),
					Errors: []map[string]string{
						{ae.ErrorCode(): ae.Error()},
					},
				}
				if errDb := job.Database.AddAwsObject(awsRecommendation, database.AwsRecommendationCollection); errDb != nil {
					return errDb
				}
			default:
				return ae
			}
		} else {
			return err
		}
	}

	return nil
}

func (job *AwsDataRetrieveJob) sumSize(contents []types.Object, c chan int64) {
	var sum int64

	for _, cnt := range contents {
		sum += int64(*cnt.Size)
	}

	c <- sum
}
