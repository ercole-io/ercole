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

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancing"
	"github.com/aws/aws-sdk-go-v2/service/rds"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	smithy "github.com/aws/smithy-go"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *AwsDataRetrieveJob) FetchObjectsCount(profile model.AwsProfile, seq uint64) error {
	cfg, err := job.loadDefaultConfig(profile)
	if err != nil {
		return err
	}
	bucketsCount, err := job.getBucketsCount(*cfg)
	if err != nil {
		return err
	}

	volumesCount, err := job.getVolumesCount(*cfg)
	if err != nil {
		return err
	}

	loadBalancersCount, err := job.getLoadBalancersCount(*cfg)
	if err != nil {
		return err
	}

	dbInstancesCount, err := job.getDbInstancesCount(*cfg)
	if err != nil {
		return err
	}

	vpcsCount, err := job.getVpcsCount(*cfg)
	if err != nil {
		return err
	}

	fileSystemsCount, err := job.getFileSystemsCount(*cfg)
	if err != nil {
		return err
	}

	res := model.AwsObject{
		SeqValue:    seq,
		ProfileID:   profile.ID,
		CreatedAt:   time.Now(),
		ProfileName: profile.Name,
		ObjectsCount: []model.AwsObjectCount{
			{
				Name:  "buckets",
				Count: *bucketsCount,
			},
			{
				Name:  "volumes",
				Count: *volumesCount,
			},
			{
				Name:  "load balancers",
				Count: *loadBalancersCount,
			},
			{
				Name:  "db instances",
				Count: *dbInstancesCount,
			},
			{
				Name:  "vpcs",
				Count: *vpcsCount,
			},
			{
				Name:  "file systems",
				Count: *fileSystemsCount,
			},
		},
	}
	if err := job.Database.AddAwsObject(res, "aws_objects"); err != nil {
		return err
	}

	return nil
}

func (job *AwsDataRetrieveJob) getBucketsCount(cfg aws.Config) (*int, error) {
	svc := s3.NewFromConfig(cfg)

	result, err := svc.ListBuckets(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Buckets)), nil
}

func (job *AwsDataRetrieveJob) getVolumesCount(cfg aws.Config) (*int, error) {
	svc := ec2.NewFromConfig(cfg)

	result, err := svc.DescribeVolumes(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Volumes)), nil
}

func (job *AwsDataRetrieveJob) getLoadBalancersCount(cfg aws.Config) (*int, error) {
	svc := elasticloadbalancing.NewFromConfig(cfg)

	result, err := svc.DescribeLoadBalancers(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.LoadBalancerDescriptions)), nil
}

func (job *AwsDataRetrieveJob) getDbInstancesCount(cfg aws.Config) (*int, error) {
	svc := rds.NewFromConfig(cfg)

	result, err := svc.DescribeDBInstances(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.DBInstances)), nil
}

func (job *AwsDataRetrieveJob) getVpcsCount(cfg aws.Config) (*int, error) {
	svc := ec2.NewFromConfig(cfg)

	result, err := svc.DescribeVpcs(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Vpcs)), nil
}

func (job *AwsDataRetrieveJob) getFileSystemsCount(cfg aws.Config) (*int, error) {
	svc := efs.NewFromConfig(cfg)

	result, err := svc.DescribeFileSystems(context.Background(), nil)
	if err != nil {
		var ae smithy.APIError
		if errors.As(err, &ae) {
			if ae.ErrorCode() == "AccessDeniedException" {
				job.Log.Warn(err)
				return aws.Int(0), nil
			}
		}

		return nil, err
	}

	return aws.Int(len(result.FileSystems)), nil
}
