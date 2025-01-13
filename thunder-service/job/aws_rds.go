// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
)

func (job *AwsDataRetrieveJob) FetchRDS(profile model.AwsProfile, seqValue uint64) error {
	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}

	rdsclient := rds.NewFromConfig(*cfg)

	res, err := rdsclient.DescribeDBInstances(context.Background(), &rds.DescribeDBInstancesInput{})
	if err != nil {
		return err
	}

	ec2client := ec2.NewFromConfig(*cfg)

	awsdbinstances := make([]model.AwsDbInstance, 0, len(res.DBInstances))

	for _, dbinstance := range res.DBInstances {
		detail, err := ec2client.DescribeInstanceTypes(context.Background(),
			&ec2.DescribeInstanceTypesInput{InstanceTypes: []types.InstanceType{
				types.InstanceType(strings.Replace(*dbinstance.DBInstanceClass, "db.", "", 1)),
			}})
		if err != nil {
			return err
		}

		if len(detail.InstanceTypes) != 0 {
			awsdbinstances = append(awsdbinstances,
				model.AwsDbInstance{DBInstance: dbinstance, InstanceTypeDetail: &detail.InstanceTypes[0]})
		}
	}

	if err := job.Database.AddAwsObject(model.AwsRDS{
		SeqValue:    seqValue,
		ProfileID:   profile.ID,
		ProfileName: profile.Name,
		Instances:   awsdbinstances,
		CreatedAt:   time.Now(),
	}, database.AwsRDSCollection); err != nil {
		return err
	}

	return nil
}
