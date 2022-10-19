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
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/efs"
	"github.com/aws/aws-sdk-go/service/elb"
	"github.com/aws/aws-sdk-go/service/rds"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/ercole-io/ercole/v2/model"
)

func (job *AwsDataRetrieveJob) FetchObjectsCount(profile model.AwsProfile, seq uint64) error {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(profile.Region),
		Credentials: credentials.NewStaticCredentials(profile.AccessKeyId, *profile.SecretAccessKey, ""),
	})
	if err != nil {
		return err
	}

	bucketsCount, err := job.getBucketsCount(sess)
	if err != nil {
		return err
	}

	volumesCount, err := job.getVolumesCount(sess)
	if err != nil {
		return err
	}

	loadBalancersCount, err := job.getLoadBalancersCount(sess)
	if err != nil {
		return err
	}

	dbInstancesCount, err := job.getDbInstancesCount(sess)
	if err != nil {
		return err
	}

	vpcsCount, err := job.getVpcsCount(sess)
	if err != nil {
		return err
	}

	fileSystemsCount, err := job.getFileSystemsCount(sess)
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

func (job *AwsDataRetrieveJob) getBucketsCount(session *session.Session) (*int, error) {
	svc := s3.New(session)

	result, err := svc.ListBuckets(nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Buckets)), nil
}

func (job *AwsDataRetrieveJob) getVolumesCount(session *session.Session) (*int, error) {
	svc := ec2.New(session)

	result, err := svc.DescribeVolumes(nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Volumes)), nil
}

func (job *AwsDataRetrieveJob) getLoadBalancersCount(session *session.Session) (*int, error) {
	svc := elb.New(session)

	result, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.LoadBalancerDescriptions)), nil
}

func (job *AwsDataRetrieveJob) getDbInstancesCount(session *session.Session) (*int, error) {
	svc := rds.New(session)

	result, err := svc.DescribeDBInstances(nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.DBInstances)), nil
}

func (job *AwsDataRetrieveJob) getVpcsCount(session *session.Session) (*int, error) {
	svc := ec2.New(session)

	result, err := svc.DescribeVpcs(nil)
	if err != nil {
		return nil, err
	}

	return aws.Int(len(result.Vpcs)), nil
}

func (job *AwsDataRetrieveJob) getFileSystemsCount(session *session.Session) (*int, error) {
	svc := efs.New(session)

	result, err := svc.DescribeFileSystems(nil)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			if aerr.Code() == "AccessDeniedException" {
				job.Log.Warn(err)
				return aws.Int(0), nil
			}
		}

		return nil, err
	}

	return aws.Int(len(result.FileSystems)), nil
}
