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
	"github.com/ercole-io/ercole/v2/utils"
)

const (
	standardEdition   = "standard"
	enterpriseEdition = "enterprise"
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
	Edition             string `json:"edition"`
	LicensesCount       int    `json:"licensesCount"`

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
		res := &AwsInstanceTypeDetail{
			InstanceType: string(m.InstanceType),
		}

		processorManufacturer := m.ProcessorInfo.Manufacturer
		if processorManufacturer != nil {
			res.ProcessorManufacturer = *processorManufacturer
		}

		defaultCore := m.VCpuInfo.DefaultCores
		if defaultCore != nil {
			res.DefaultCore = int(*defaultCore)
		}

		defaultThreadsPerCore := m.VCpuInfo.DefaultThreadsPerCore
		if defaultThreadsPerCore != nil {
			res.DefaultThreadsPerCore = int(*defaultThreadsPerCore)
		}

		defaultVCpus := m.VCpuInfo.DefaultVCpus
		if defaultVCpus != nil {
			res.DefaultVCpus = int(*defaultVCpus)
		}

		sizeInMiB := m.MemoryInfo.SizeInMiB
		if sizeInMiB != nil {
			res.MemorySizeInMib = int(*sizeInMiB)
		}

		return res
	}

	return nil
}

func ToAwsDbInstance(m *model.AwsDbInstance) *AwsDbInstance {
	if m != nil {
		res := &AwsDbInstance{}

		dbName := m.DBName
		if dbName != nil {
			res.DbName = *dbName
		}

		dbInstanceClass := m.DBInstanceClass
		if dbInstanceClass != nil {
			res.DbInstanceClass = *dbInstanceClass
		}

		engine := m.Engine
		if engine != nil {
			res.Engine = *engine
		}

		res.Edition = getAwsRdsEdition(engine)

		engineVersion := m.EngineVersion
		if engineVersion != nil {
			res.EngineVersion = *engineVersion
		}

		dbInstanceStatus := m.DBInstanceStatus
		if dbInstanceStatus != nil {
			res.DbInstanceStatus = *dbInstanceStatus
		}

		licenseModel := m.LicenseModel
		if licenseModel != nil {
			res.LicenseModel = *licenseModel
		}

		storageType := m.StorageType
		if storageType != nil {
			res.StorageType = *storageType
		}

		allocatedStorage := m.AllocatedStorage
		if allocatedStorage != nil {
			res.AllocatedStorage = int(*allocatedStorage)
		}

		maxAllocatedStorage := m.MaxAllocatedStorage
		if maxAllocatedStorage != nil {
			res.MaxAllocatedStorage = int(*maxAllocatedStorage)
		}

		awsInstanceTypeDetail := ToAwsInstanceTypeDetail(m.InstanceTypeDetail)
		if awsInstanceTypeDetail != nil {
			res.AwsInstanceTypeDetail = *awsInstanceTypeDetail
		}

		res.LicensesCount = getLicensesCount(res.Edition, res.DefaultVCpus)

		return res
	}

	return nil
}

func ToAwsDbInstances(list []model.AwsDbInstance) []AwsDbInstance {
	res := make([]AwsDbInstance, 0, len(list))

	for _, v := range list {
		awsDbInstance := ToAwsDbInstance(&v)

		if awsDbInstance != nil {
			res = append(res, *awsDbInstance)
		}
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

func getAwsRdsEdition(engine *string) string {
	if engine == nil {
		return ""
	}

	standards := []string{
		"custom-oracle-se2",
		"custom-oracle-se2-cdb",
		"oracle-se2",
		"oracle-se2-cdb",
	}

	enterprises := []string{
		"custom-oracle-ee",
		"custom-oracle-ee-cdb",
		"oracle-ee",
		"oracle-ee-cdb",
	}

	if utils.Contains(standards, *engine) {
		return standardEdition
	}

	if utils.Contains(enterprises, *engine) {
		return enterpriseEdition
	}

	return ""
}

func getLicensesCount(edition string, vcpus int) int {
	var licensesCount int

	if edition == standardEdition {
		licensesCount = vcpus / 4
	}

	if edition == enterpriseEdition {
		licensesCount = vcpus / 2
	}

	if licensesCount == 0 {
		licensesCount = 1
	}

	return licensesCount
}
