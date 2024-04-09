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
package dto

import (
	ec2types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/ercole-io/ercole/v2/model"
)

type AwsRDSResponse struct {
	AwsDbInstances []AwsDbInstance `json:"awsDbInstances"`
}

type AwsDbInstance struct {
	DbName              string `json:"dbName"`
	DbInstanceClass     string `json:"dbInstanceClass"`
	Engine              string `json:"engine"`
	EngineVersion       string `json:"engineVersion"`
	DbInstanceStatus    string `json:"dbInstanceStatus"`
	LicenseModel        string `json:"licenseModel"`
	StorageType         string `json:"storageType"`
	AllocatedStorage    int    `json:"allocatedStorage"`
	MaxAllocatedStorage int    `json:"maxAllocatedStorage"`

	AwsInstanceTypeDetail
}

type AwsInstanceTypeDetail struct {
	InstanceType          string `json:"instanceType"`
	ProcessorManufacturer string `json:"processorManufacturer"`
	DefaultCore           int    `json:"defaultCore"`
	DefaultThreadsPerCore int    `json:"defaultThreadsPerCore"`
	DefaultVCpus          int    `json:"defaultVCpus"`
	MemorySizeInMib       int    `json:"memorySizeInMib"`
}

func ToAwsInstanceTypeDetail(m *ec2types.InstanceTypeInfo) *AwsInstanceTypeDetail {
	if m != nil {
		return &AwsInstanceTypeDetail{
			InstanceType:          string(m.InstanceType),
			ProcessorManufacturer: *m.ProcessorInfo.Manufacturer,
			DefaultCore:           int(*m.VCpuInfo.DefaultCores),
			DefaultThreadsPerCore: int(*m.VCpuInfo.DefaultThreadsPerCore),
			DefaultVCpus:          int(*m.VCpuInfo.DefaultVCpus),
			MemorySizeInMib:       int(*m.MemoryInfo.SizeInMiB),
		}
	}

	return nil
}

func ToAwsDbInstance(m *model.AwsDbInstance) *AwsDbInstance {
	if m != nil {
		return &AwsDbInstance{
			DbName:              *m.DBName,
			DbInstanceClass:     *m.DBInstanceClass,
			Engine:              *m.Engine,
			EngineVersion:       *m.EngineVersion,
			DbInstanceStatus:    *m.DBInstanceStatus,
			LicenseModel:        *m.LicenseModel,
			StorageType:         *m.StorageType,
			AllocatedStorage:    int(*m.AllocatedStorage),
			MaxAllocatedStorage: int(*m.MaxAllocatedStorage),

			AwsInstanceTypeDetail: *ToAwsInstanceTypeDetail(m.InstanceTypeDetail),
		}
	}

	return nil
}

func ToAwsDbInstances(list []model.AwsDbInstance) []AwsDbInstance {
	res := make([]AwsDbInstance, 0, len(list))

	for _, v := range list {
		res = append(res, *ToAwsDbInstance(&v))
	}

	return res
}

func ToAwsRDSResponse(list []model.AwsRDS) AwsRDSResponse {
	instances := make([]AwsDbInstance, 0, len(list))

	for _, rds := range list {
		instances = append(instances, ToAwsDbInstances(rds.Instances)...)
	}

	return AwsRDSResponse{instances}
}
