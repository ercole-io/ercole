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
	Type            string `json:"type"`
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

			sTime := common.SDKTime{Time: time.Now().Local().AddDate(0, 0, -8)}
			eTime := common.SDKTime{Time: time.Now().Local()}

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
	var listCompartments []model.OciCompartment
	var instancesNotOptimizable map[string]Instance
	var allInstancesWithMetrics map[string]Instance
	var allInstances map[string]Instance
	var AvgCPUThreshold = 3
	var PeakCPUThreshold = 180
	var MemoryThreshold = 1
	instancesNotOptimizable = make(map[string]Instance)
	allInstancesWithMetrics = make(map[string]Instance)
	allInstances = make(map[string]Instance)

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

		monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		// retrieve metrics data for each compartment
		for _, compartment := range listCompartments {

			allInstances, err = as.getOciInstancesWithoutBasicShape(allInstances, compartment, tenancyOCID, customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			allInstancesWithMetrics, err = as.getOciInstancesWithMetricsWithoutBasicShape(allInstancesWithMetrics, compartment, tenancyOCID, customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			// first query is about average CPU utilization in the last  90 days
			var strQueryAvgCPU = "CpuUtilization[1d].avg()>50"

			sTime := common.SDKTime{Time: time.Now().Local().AddDate(0, 0, -89)}
			eTime := common.SDKTime{Time: time.Now().Local()}

			instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryAvgCPU, sTime, eTime, instancesNotOptimizable, compartment, AvgCPUThreshold)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			// second query is about CPU utilization peak in the last 7 days
			var strQueryPeakCPU = "CpuUtilization[1m].max()>50"
			sTime = common.SDKTime{Time: time.Now().Local().AddDate(0, 0, -7)}
			eTime = common.SDKTime{Time: time.Now().Local()}

			instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryPeakCPU, sTime, eTime, instancesNotOptimizable, compartment, PeakCPUThreshold)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			// third query is about memory utilization in the last 7 days
			var strQueryMemory = "MemoryUtilization[1m].max()>90"
			sTime = common.SDKTime{Time: time.Now().Local().AddDate(0, 0, -7)}
			eTime = common.SDKTime{Time: time.Now().Local()}

			instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryMemory, sTime, eTime, instancesNotOptimizable, compartment, MemoryThreshold)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
		}
	}

	var recommendation model.OciErcoleRecommendation

	for _, a := range instancesNotOptimizable {
		// if an instance is not optimizable I have to remove it from the list
		if val, ok := allInstancesWithMetrics[a.ResourceID]; ok {
			delete(allInstancesWithMetrics, val.ResourceID)
		}
	}

	// build recommendation data for optimizable instances
	for _, inst := range allInstancesWithMetrics {
		if allInstances[inst.ResourceID].Type == "kubernetes" {
			recommendation.Type = model.RecommendationTypeSISRightsizing1
		} else {
			recommendation.Type = model.RecommendationTypeInstanceRightsizing
		}
		recommendation.CompartmentID = inst.CompartmentID
		recommendation.CompartmentName = inst.CompartmentName
		recommendation.ResourceID = inst.ResourceID
		recommendation.Name = inst.Name
		listRec = append(listRec, recommendation)
	}

	for _, b := range allInstancesWithMetrics {
		delete(allInstances, b.ResourceID)
	}

	// build recommendation data for instances withoout monitoring
	for _, in := range allInstances {
		if in.Type == "kubernetes" {
			recommendation.Type = model.RecommendationTypeSISRightsizing1
		} else {
			recommendation.Type = model.RecommendationTypeInstanceWithoutMonitoring
		}
		recommendation.CompartmentID = in.CompartmentID
		recommendation.CompartmentName = in.CompartmentName
		recommendation.ResourceID = in.ResourceID
		recommendation.Name = in.Name
		listRec = append(listRec, recommendation)
	}

	return listRec, merr
}

func (as *ThunderService) getOciInstancesWithoutBasicShape(allInstances map[string]Instance, compartment model.OciCompartment, tenancyOCID string, customConfigProvider common.ConfigurationProvider) (map[string]Instance, error) {
	client, err := core.NewComputeClientWithConfigurationProvider(customConfigProvider)
	if err != nil {
		return allInstances, err
	}

	req := core.ListInstancesRequest{
		CompartmentId: &compartment.CompartmentID,
	}

	// Send the request using the service client
	resp, err := client.ListInstances(context.Background(), req)
	if err != nil {
		return allInstances, err
	}

	for _, s := range resp.Items {
		var tmpInstance Instance

		// we have to exclude the instances with basic shape
		if *s.Shape != "VM.Standard2.1" && *s.Shape != "VM.StandardE2.1" {
			tmpInstance.CompartmentID = compartment.CompartmentID
			tmpInstance.CompartmentName = compartment.Name
			tmpInstance.ResourceID = *s.Id
			tmpInstance.Name = *s.DisplayName
			tmpInstance.Shape = *s.Shape
			if _, ok := s.Metadata["oke-pool-id"]; ok {
				tmpInstance.Type = "kubernetes"
			} else {
				tmpInstance.Type = "normal"
			}
			allInstances[*s.Id] = tmpInstance
		}
	}

	return allInstances, nil
}

func (as *ThunderService) countEventsOccurence(client monitoring.MonitoringClient, strQuery string, sTime common.SDKTime, eTime common.SDKTime, instances map[string]Instance, compartment model.OciCompartment, threshold int) (map[string]Instance, error) {

	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: &compartment.CompartmentID,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			StartTime: &sTime,
			EndTime:   &eTime,
			Namespace: common.String("oci_computeagent"),
			Query:     &strQuery,
		},
	}

	// Send the request using the service client
	resp, err := client.SummarizeMetricsData(context.Background(), req)
	if err != nil {
		return instances, err
	}

	var instance Instance
	for _, s := range resp.Items {
		// reset the counter
		cnt := 0
		for _, a := range s.AggregatedDatapoints {
			if *a.Value == 1.0 {
				cnt++
			}
		}
		if cnt > threshold {
			// the instance is not eligible for optimization
			if s.Dimensions["shape"] != "VM.Standard2.1" && s.Dimensions["shape"] != "VM.StandardE2.1" {
				if val, ok := instances[s.Dimensions["resourceId"]]; ok {
					val.Cnt += 1
					instances[s.Dimensions["resourceId"]] = val

				} else {
					instance.CompartmentID = compartment.CompartmentID
					instance.CompartmentName = compartment.Name
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

func (as *ThunderService) getOciInstancesWithMetricsWithoutBasicShape(instances map[string]Instance, compartment model.OciCompartment, tenancyOCID string, customConfigProvider common.ConfigurationProvider) (map[string]Instance, error) {

	client, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)
	if err != nil {
		return instances, err
	}

	req := monitoring.ListMetricsRequest{
		CompartmentId:      &compartment.CompartmentID,
		ListMetricsDetails: monitoring.ListMetricsDetails{Namespace: common.String("oci_computeagent")},
	}

	// Send the request using the service client
	resp, err := client.ListMetrics(context.Background(), req)
	if err != nil {
		return instances, err
	}

	for _, s := range resp.Items {
		if s.Dimensions["shape"] != "VM.StandardE2.1" && s.Dimensions["shape"] != "VM.Standard2.1" {
			// if the instance is not in the list I have to put it
			if _, ok := instances[s.Dimensions["resourceId"]]; !ok {
				tmpInstance := Instance{*s.CompartmentId, compartment.Name, s.Dimensions["resourceId"], s.Dimensions["resourceDisplayName"], s.Dimensions["shape"], 1, ""}
				instances[s.Dimensions["resourceId"]] = tmpInstance
			}
		}
	}

	return instances, nil
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
				resp, err := as.getMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeReadThroughput[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeWriteThroughput[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeReadOps[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "oci_blockstore", "VolumeWriteOps[5d].max()")

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

func (as *ThunderService) getMetricResponse(client monitoring.MonitoringClient, compartmentId string, namespace string, query string) (*monitoring.SummarizeMetricsDataResponse, error) {
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
	retList := make(map[string]string, 0)

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
