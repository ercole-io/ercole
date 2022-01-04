// Copyright (c) 2020 Sorint.lab S.p.A.
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
package service

import (
	"context"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/core"
	"github.com/oracle/oci-go-sdk/v45/monitoring"

	"github.com/ercole-io/ercole/v2/model"
)

type Instance struct {
	CompartmentID   string `json:"compartmentID"`
	CompartmentName string `json:"compartmentName"`
	ResourceID      string `json:"resourceID"`
	Name            string `json:"name"`
	Shape           string `json:"shape"`
	Cnt             int    `json:"cnt"`
}

func (as *ThunderService) GetOciComputeInstancesIdle(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var listRec []model.OciErcoleRecommendation
	var merr error
	var listCompartments []model.OciCompartment

	listRec = make([]model.OciErcoleRecommendation, 0)

	for _, profileId := range profiles {

		customConfigProvider, tenancyOCID, err := as.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = as.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		// retrieve metrics data for each compartment
		var strNamespace = "oci_compute_infrastructure_health"

		for _, compartment := range listCompartments {

			instances, err := as.getOciInstances(customConfigProvider, compartment.CompartmentID)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			// query for instance status in the last 8 days
			var strQueryInstanceIdle = "instance_status[5m].mean()==0"

			sTime := common.SDKTime{time.Now().Local().AddDate(0, 0, -8)}
			eTime := common.SDKTime{time.Now().Local()}

			monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			req := monitoring.SummarizeMetricsDataRequest{
				CompartmentId: &compartment.CompartmentID,
				SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
					StartTime: &sTime,
					EndTime:   &eTime,
					Namespace: &strNamespace,
					Query:     &strQueryInstanceIdle,
				},
			}

			resp, err := monClient.SummarizeMetricsData(context.Background(), req)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			var recommendation model.OciErcoleRecommendation

		items:
			for _, s := range resp.Items {
				for _, a := range s.AggregatedDatapoints {
					if *a.Value == 1.0 {
						delete(instances, s.Dimensions["resourceId"])
						continue items
					}
				}
			}

			for id, value := range instances {
				recommendation.Type = model.RecommendationTypeComputeInstanceIdle
				recommendation.CompartmentID = compartment.CompartmentID
				recommendation.CompartmentName = compartment.Name
				recommendation.Name = value
				recommendation.ResourceID = id
				listRec = append(listRec, recommendation)
			}
		}
	}
	return listRec, merr
}

func (as *ThunderService) GetOciComputeInstanceRightsizing(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var listRec []model.OciErcoleRecommendation
	var merr error
	var err error
	var listCompartments []model.OciCompartment
	var instances map[string]Instance

	instances = make(map[string]Instance)

	listCompartments, err = as.GetOciCompartments(profiles)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	// retrieve metrics data for each compartment
	var strNamespace = "oci_computeagent"
	var AvgCPUThreshold = 3
	var PeakCPUThreshold = 180
	var MemoryThreshold = 1

	for _, compartment := range listCompartments {

		// first query is about CPU
		var strQueryAvgCPU = "CpuUtilization[1d].avg()>50" // average CPU utilization in the last  90 days

		sTime := common.SDKTime{time.Now().Local().AddDate(0, 0, -90)}
		eTime := common.SDKTime{time.Now().Local()}

		monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(common.DefaultConfigProvider())
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		req := monitoring.SummarizeMetricsDataRequest{
			CompartmentId: &compartment.CompartmentID,
			SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
				StartTime: &sTime,
				EndTime:   &eTime,
				Namespace: &strNamespace,
				Query:     &strQueryAvgCPU,
			},
		}

		instances, err = as.CountEventsOccurence(monClient, req, instances, compartment.CompartmentID, AvgCPUThreshold)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		var strQueryPeakCPU = "CpuUtilization[1m].max()>50" // CPU utilization peak in the last 7 days
		sTime = common.SDKTime{time.Now().Local().AddDate(0, 0, -7)}
		eTime = common.SDKTime{time.Now().Local()}

		req = monitoring.SummarizeMetricsDataRequest{
			CompartmentId: &compartment.CompartmentID,
			SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
				StartTime: &sTime,
				EndTime:   &eTime,
				Namespace: &strNamespace,
				Query:     &strQueryPeakCPU,
			},
		}

		instances, err = as.CountEventsOccurence(monClient, req, instances, compartment.CompartmentID, PeakCPUThreshold)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		var strQueryMemory = "MemoryUtilization[1m].max()>90" // memory utilization in the last 7 days
		sTime = common.SDKTime{time.Now().Local().AddDate(0, 0, -7)}
		eTime = common.SDKTime{time.Now().Local()}

		req = monitoring.SummarizeMetricsDataRequest{
			CompartmentId: &compartment.CompartmentID,
			SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
				StartTime: &sTime,
				EndTime:   &eTime,
				Namespace: &strNamespace,
				Query:     &strQueryMemory,
			},
		}

		instances, err = as.CountEventsOccurence(monClient, req, instances, compartment.CompartmentID, MemoryThreshold)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}
	}

	// if an instance has been counting far all the typologies it can be optimized
	var recommendation model.OciErcoleRecommendation

	for _, a := range instances {
		if a.Cnt == 3 {
			recommendation.CompartmentID = a.CompartmentID
			recommendation.CompartmentName = a.CompartmentName
			recommendation.ResourceID = a.ResourceID
			recommendation.Name = a.Name
			listRec = append(listRec, recommendation)
		}

	}
	return listRec, merr
}

func (as *ThunderService) CountEventsOccurence(client monitoring.MonitoringClient, req monitoring.SummarizeMetricsDataRequest, instances map[string]Instance, compartmentId string, threshold int) (map[string]Instance, error) {

	// Send the request using the service client
	resp, err := client.SummarizeMetricsData(context.Background(), req)
	if err != nil {
		return instances, err
	}

	var instance Instance
	var cnt int
	for _, s := range resp.Items {
		for _, a := range s.AggregatedDatapoints {
			if *a.Value == 1.0 {
				cnt++
			}
		}
		if cnt > threshold {
			// the instance is eligible for optimization
			if s.Dimensions["shape"] != "VM.StandardE2.1" {
				if val, ok := instances[s.Dimensions["resourceId"]]; ok {
					val.Cnt += 1
					instances[s.Dimensions["resourceId"]] = val

				} else {
					instance.CompartmentID = compartmentId
					instance.ResourceID = s.Dimensions["resourceId"]
					instance.Name = s.Dimensions["resourceDisplayName"]
					instance.Shape = s.Dimensions["shape"]
					instances[s.Dimensions["resourceId"]] = instance
				}
			}
		}
	}
	return instances, nil
}

func (as *ThunderService) GetMetricResponse(client monitoring.MonitoringClient, compartmentId string, namespace string, query string) (*monitoring.SummarizeMetricsDataResponse, error) {
	var merr error

	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: &compartmentId,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Namespace: &namespace,
			Query:     &query,
		},
	}

	resp, err := client.SummarizeMetricsData(context.Background(), req)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	return &resp, nil
}

func (as *ThunderService) GetOciBlockStorageRightsizing(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var listRec []model.OciErcoleRecommendation
	var merr error
	var listCompartments []model.OciCompartment
	var recommendation model.OciErcoleRecommendation

	var vol model.OciResourcePerformance
	listRec = make([]model.OciErcoleRecommendation, 0)

	for _, profileId := range profiles {

		customConfigProvider, tenancyOCID, err := as.getOciCustomConfigProviderAndTenancy(profileId)

		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = as.getOciProfileCompartments(tenancyOCID, customConfigProvider)

		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		var resTmp model.OciResourcePerformance
		var ok bool

		monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			return nil, merr
		}

		// retrieve metrics data for each compartment
		for _, compartment := range listCompartments {

			var vols = make(map[string]model.OciResourcePerformance)

			coreClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			req := core.ListVolumesRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp1, err := coreClient.ListVolumes(context.Background(), req)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			if len(resp1.Items) > 0 {
				for _, r := range resp1.Items {
					vol = model.OciResourcePerformance{
						ResourceID: *r.Id,
						Name:       *r.DisplayName,
						Size:       int(*r.SizeInGBs),
						VpusPerGB:  int(*r.VpusPerGB),
						Throughput: 0.0,
						Iops:       0,
					}
					vols[*r.Id] = vol
				}

				// first query is about Read Throughput
				resp, err := as.GetMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeReadThroughput[5d].max()")

				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				for _, s := range resp.Items {
					tempId := s.Dimensions["resourceId"]
					if resTmp, ok = vols[tempId]; ok {
						if s.Metadata["unit"] == "bytes" {
							resTmp.Throughput += *s.AggregatedDatapoints[0].Value / 1024 / 1024
						}
						vols[tempId] = resTmp
					}
				}

				// second query is about Write Throughput
				resp, err = as.GetMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeWriteThroughput[5d].max()")

				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				for _, s := range resp.Items {
					tempId := s.Dimensions["resourceId"]

					if resTmp, ok = vols[tempId]; ok {
						if s.Metadata["unit"] == "bytes" {
							resTmp.Throughput += *s.AggregatedDatapoints[0].Value / 1024 / 1024
						}
						vols[tempId] = resTmp
					}
				}

				// third query is about Read Ops
				resp, err = as.GetMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeReadOps[5d].max()")

				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				for _, s := range resp.Items {
					tempId := s.Dimensions["resourceId"]
					if resTmp, ok = vols[tempId]; ok {
						if s.Metadata["unit"] == "operations" {
							resTmp.Iops += int(*s.AggregatedDatapoints[0].Value)
						}
						vols[tempId] = resTmp
					}
				}

				// fourth query is about Write Ops
				resp, err = as.GetMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeWriteOps[5d].max()")

				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				for _, s := range resp.Items {
					tempId := s.Dimensions["resourceId"]

					if resTmp, ok = vols[tempId]; ok {
						if s.Metadata["unit"] == "operations" {
							resTmp.Iops += int(*s.AggregatedDatapoints[0].Value)
						}
						vols[tempId] = resTmp
					}
				}

				// N.B devo resettare la mappa ad ogni compartment

				if len(vols) != 0 {
					for _, v := range vols {
						isOpt, err := as.isOptimizable(v)
						if err != nil {
							merr = multierror.Append(merr, err)
							continue
						}
						if isOpt {
							recommendation.Type = model.RecommendationTypeBlockStorage
							recommendation.CompartmentID = compartment.CompartmentID
							recommendation.CompartmentName = compartment.Name
							recommendation.ResourceID = v.ResourceID
							recommendation.Name = v.Name
							listRec = append(listRec, recommendation)
						}
					}
				}
			}
		}
	}

	return listRec, merr
}

func (as *ThunderService) isOptimizable(res model.OciResourcePerformance) (bool, error) {
	var ociPerfs *model.OciVolumePerformance

	if res.VpusPerGB == 0 {
		return false, nil
	}

	ociPerfs = as.getOciVolumePerformance(res.VpusPerGB, res.Size)

	if res.Throughput < (ociPerfs.Performances[0].Values.MaxThroughput/2.0) && res.Iops < (ociPerfs.Performances[0].Values.MaxIOPS)/2.0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (as *ThunderService) getOciVolumePerformance(vpu int, size int) *model.OciVolumePerformance {
	var baseIopsPerGB float64
	var maxIops int
	var baseThroughput float64
	var maxTroughput float64
	var retThroughput float64
	var retIOPS int

	if vpu != 0 {
		baseIopsPerGB = 1.5*float64(vpu) + 45
		maxIops = 2500 * vpu
		baseThroughput = (12*float64(vpu) + 360) / 1000
		maxTroughput = 20*float64(vpu) + 280
	} else {
		baseIopsPerGB = 2
		maxIops = 3000
		baseThroughput = 240.0 / 15.0 / 1000.0
		maxTroughput = 480 / 15
	}

	var valRet model.OciVolumePerformance
	var perfTmp model.OciPerformance
	var valTmp model.OciPerfValues

	valTmp.MaxThroughput = baseThroughput * float64(size)
	if retThroughput > maxTroughput {
		valTmp.MaxThroughput = maxTroughput
	}

	valTmp.MaxIOPS = int(baseIopsPerGB) * size
	if retIOPS > maxIops {
		valTmp.MaxIOPS = maxIops
	}

	valRet.Vpu = vpu
	perfTmp.Size = size
	perfTmp.Values = valTmp
	valRet.Performances = append(valRet.Performances, perfTmp)

	return &valRet
}

func (as *ThunderService) getOciInstances(customConfigProvider common.ConfigurationProvider, compartmentID string) (map[string]string, error) {
	var merr error
	var retList map[string]string
	retList = make(map[string]string)

	client, err := core.NewComputeClientWithConfigurationProvider(customConfigProvider)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	req := core.ListInstancesRequest{
		CompartmentId: &compartmentID,
	}

	resp, err := client.ListInstances(context.Background(), req)

	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	for _, s := range resp.Items {
		retList[*s.Id] = *s.DisplayName
	}

	return retList, nil

}
