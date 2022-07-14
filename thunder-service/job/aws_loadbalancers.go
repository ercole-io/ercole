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
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/elbv2"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *AwsDataRetrieveJob) GetAwsUnusedLoadBalancers(profiles []model.AwsProfile, seqValue uint64) {
	var are model.AwsRecommendationError

	var recommendation model.AwsRecommendation

	listRec := make([]interface{}, 0)
	errors := make([]interface{}, 0)

	for _, profile := range profiles {
		sess, err := session.NewSession(&aws.Config{
			Region:      aws.String(profile.Region),
			Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
		})
		if err != nil {
			recError := are.SetAwsRecommendationError(seqValue, "", time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		elbSvc := elb.New(sess)

		var resultelbSvc *elb.DescribeInstanceHealthOutput

		var errSvc error

		resultelbSvc_c, err := elbSvc.DescribeLoadBalancers(nil)
		if err != nil {
			recError := are.SetAwsRecommendationError(seqValue, "", time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		for _, l := range resultelbSvc_c.LoadBalancerDescriptions {
			if *l.LoadBalancerName != "" {
				params := elb.DescribeInstanceHealthInput{
					LoadBalancerName: l.LoadBalancerName,
				}
				resultelbSvc, errSvc = elbSvc.DescribeInstanceHealth(&params)

				if errSvc != nil {
					recError := are.SetAwsRecommendationError(seqValue, "", time.Now().UTC(), err.Error())
					errors = append(errors, recError)

					continue
				}

				for _, p := range resultelbSvc.InstanceStates {
					if *p.State == "OutOfService" {
						recommendation.SeqValue = seqValue
						recommendation.ProfileID = profile.ID.Hex()
						recommendation.Category = model.AwsUnusedResource
						recommendation.Suggestion = model.AwsDeleteLoadBalancerNotActive
						recommendation.Name = ""
						recommendation.ResourceID = *p.InstanceId
						recommendation.ObjectType = model.AwsObjectTypeLoadBalancer
						recommendation.Details = []map[string]interface{}{
							{"Resource Id": *p.InstanceId},
							{"Resource Type": "Load Balancer"},
							{"Resource Status": fmt.Sprintf("%v", *p.State)},
						}
						recommendation.CreatedAt = time.Now().UTC()

						listRec = append(listRec, recommendation)
					}
				}
			}
		}

		elbv2Svc := elbv2.New(sess)

		resultelbv2Svc, err := elbv2Svc.DescribeLoadBalancers(nil)
		if err != nil {
			recError := are.SetAwsRecommendationError(seqValue, "", time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		for _, l := range resultelbv2Svc.LoadBalancers {
			if *l.State.Code == "failed" {
				recommendation.SeqValue = seqValue
				recommendation.ProfileID = profile.ID.Hex()
				recommendation.Category = model.AwsUnusedResource
				recommendation.Suggestion = model.AwsDeleteLoadBalancerNotActive
				recommendation.Name = ""
				recommendation.ResourceID = *l.VpcId
				recommendation.ObjectType = model.AwsObjectTypeLoadBalancer
				recommendation.Details = []map[string]interface{}{
					{"Resource Id": *l.VpcId},
					{"Resource Type": "Load Balancer"},
					{"Resource Status": fmt.Sprintf("%v", *l.State.Code)},
				}
				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddAwsObjects(listRec, "aws_recommendations")

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}

	if len(errors) > 0 {
		errDb := job.Database.AddAwsObjects(errors, "aws_recommendations_errors")

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}
