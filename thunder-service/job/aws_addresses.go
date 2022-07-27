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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
)

func (job *AwsDataRetrieveJob) FetchAwsUnusedIPAddresses(profile model.AwsProfile, seqValue uint64) error {
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

	resultec2Svc, err := ec2Svc.DescribeAddresses(nil)
	if err != nil {
		return err
	}

	for _, w := range resultec2Svc.Addresses {
		if *w.AssociationId == "" {
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
			recommendation.Suggestion = model.AwsDeletePublicIPAddressNotAssociated
			recommendation.Name = objectName
			recommendation.ResourceID = *w.AllocationId
			recommendation.ObjectType = model.AwsPublicID
			recommendation.Details = []map[string]interface{}{
				{"RESOURCE_NAME": objectName},
				{"RESOURCE_TYPE": "Public IP"},
				{"RESOURCE_STATUS": "Not associated"},
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
