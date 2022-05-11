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
	"fmt"
	"math"
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/core"
	"github.com/oracle/oci-go-sdk/v45/monitoring"

	"github.com/ercole-io/ercole/v2/model"
)

type Instance struct {
	CompartmentID   string  `json:"compartmentID"`
	CompartmentName string  `json:"compartmentName"`
	ResourceID      string  `json:"resourceID"`
	Name            string  `json:"name"`
	ClusterName     string  `json:"clusterName"`
	Shape           string  `json:"shape"`
	Cnt             int     `json:"cnt"`
	Type            string  `json:"type"`
	Status          string  `json:"status"`
	OCPUs           float32 `json:"ocpus"`
}

type MetricData map[string]ValThresh

type ValThresh struct {
	Perc      string
	Max       string
	Value     string
	Threshold string
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

			var sTime common.SDKTime

			var eTime common.SDKTime

			sTime.Time = time.Now().Local().AddDate(0, 0, -8)
			eTime.Time = time.Now().Local()

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
				recommendation.Details = make([]model.RecDetail, 0)
				if value.Type == "kubernetes" {
					recommendation.Category = model.UnusedServiceDecommisioning
					recommendation.Suggestion = model.DeleteKubernetesNodeNotActive
					recommendation.ObjectType = model.ObjectTypeClusterKubernetes
				} else {
					recommendation.Category = model.ComputeInstanceIdle
					recommendation.Suggestion = model.DeleteComputeInstanceNotActive
					recommendation.ObjectType = model.ObjectTypeComputeInstance
				}

				recommendation.CompartmentID = compartment.CompartmentID
				recommendation.CompartmentName = compartment.Name
				recommendation.Name = value.Name
				recommendation.ResourceID = id
				detail1 := model.RecDetail{Name: "Instance Name", Value: value.Name}
				detail2 := model.RecDetail{Name: "Instance Shape", Value: value.Shape}

				recommendation.Details = append(recommendation.Details, detail1, detail2)

				if value.Type == "kubernetes" {
					detail3 := model.RecDetail{Name: "Oke Cluster Name", Value: value.ClusterName}

					recommendation.Details = append(recommendation.Details, detail3)
				}

				listRec = append(listRec, recommendation)
			}
		}
	}

	return listRec, merr
}

func (as *ThunderService) GetOciComputeInstanceRightsizing(profiles []string) ([]model.OciErcoleRecommendation, error) {
	return as.getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles, 50, 50, 90, true, "rightsizing")
}

func (as *ThunderService) GetOciUnusedServiceDecommisioning(profiles []string) ([]model.OciErcoleRecommendation, error) {
	return as.getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles, 5, 5, 40, false, "decommisioning")
}

func (as *ThunderService) getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles []string, percAvgCPU int, percPeakCPU int, percMemoryUtilization int, verifyShape bool, recommType string) ([]model.OciErcoleRecommendation, error) {
	var listRec []model.OciErcoleRecommendation

	var merr error

	var listCompartments []model.OciCompartment

	var AvgCPUThreshold = 3

	var PeakCPUThreshold = 180

	var MemoryThreshold = 1

	var recClusterName model.RecDetail

	instancesNotOptimizable := make(map[string]Instance)
	allInstancesWithMetrics := make(map[string]Instance)
	allInstances := make(map[string]Instance)

	allInstanceMetrics := make(map[string]MetricData)

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
			allInstances, err = as.getOciInstancesList(allInstances, compartment, customConfigProvider, verifyShape)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			allInstancesWithMetrics, err = as.getOciInstancesWithMetrics(allInstancesWithMetrics, compartment, customConfigProvider, verifyShape)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			for _, inst := range allInstancesWithMetrics {
				// first query is about average CPU utilization in the last  90 days
				var strQueryAvgCPU = fmt.Sprintf("CpuUtilization[1d]{resourceId=\"%s\"}.avg()", inst.ResourceID)

				var sTime common.SDKTime

				var eTime common.SDKTime

				sTime.Time = time.Now().Local().AddDate(0, 0, -89)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryAvgCPU, sTime, eTime, instancesNotOptimizable, compartment, AvgCPUThreshold, percAvgCPU, verifyShape, allInstanceMetrics, "AvgCPU")
				if err != nil {
					merr = multierror.Append(merr, err)
				}

				// second query is about CPU utilization peak in the last 7 days
				var strQueryPeakCPU = fmt.Sprintf("CpuUtilization[1m]{resourceId=\"%s\"}.max()", inst.ResourceID)

				sTime.Time = time.Now().Local().AddDate(0, 0, -7)
				sTime.Time = sTime.Add(time.Hour * 1)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryPeakCPU, sTime, eTime, instancesNotOptimizable, compartment, PeakCPUThreshold, percPeakCPU, verifyShape, allInstanceMetrics, "PeakCPU")
				if err != nil {
					merr = multierror.Append(merr, err)
				}

				// third query is about memory utilization in the last 7 days
				var strQueryMemory = fmt.Sprintf("MemoryUtilization[1m]{resourceId=\"%s\"}.max()", inst.ResourceID)

				sTime.Time = time.Now().Local().AddDate(0, 0, -7)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = as.countEventsOccurence(monClient, strQueryMemory, sTime, eTime, instancesNotOptimizable, compartment, MemoryThreshold, percMemoryUtilization, verifyShape, allInstanceMetrics, "AvgMemory")
				if err != nil {
					merr = multierror.Append(merr, err)
				}
			}
		}
	}

	var recommendation model.OciErcoleRecommendation

	//allInstancesWithoutMetrics := allInstances
	allInstancesWithoutMetrics := make(map[string]Instance)
	for key, value := range allInstances {
		allInstancesWithoutMetrics[key] = value
	}

	for _, b := range allInstancesWithMetrics {
		delete(allInstancesWithoutMetrics, b.ResourceID)
	}

	for _, a := range instancesNotOptimizable {
		// if an instance is not optimizable I have to remove it from the list
		if val, ok := allInstancesWithMetrics[a.ResourceID]; ok {
			delete(allInstancesWithMetrics, val.ResourceID)
		}
	}

	// build recommendation data for optimizable instances
	for _, inst := range allInstancesWithMetrics {
		if allInstances[inst.ResourceID].Status != "STOPPED" {
			recommendation.Details = make([]model.RecDetail, 0)

			if allInstances[inst.ResourceID].Type == "kubernetes" {
				recClusterName = model.RecDetail{Name: "Oke Cluster Name", Value: allInstances[inst.ResourceID].ClusterName}

				if recommType == "rightsizing" {
					recommendation.Category = model.SISRightsizing1
					recommendation.Suggestion = model.ResizeOversizedKubernetesCluster
				} else {
					recommendation.Category = model.UnusedServiceDecommisioning
					recommendation.Suggestion = model.DeleteKubernetesNodeNotUsed
				}

				recommendation.ObjectType = model.ObjectTypeClusterKubernetes
			} else {
				if recommType == "rightsizing" {
					recommendation.Category = model.InstanceRightsizing
					recommendation.Suggestion = model.ResizeOversizedComputeInstance
				} else {
					recommendation.Category = model.ComputeInstanceDecommisioning
					recommendation.Suggestion = model.DeleteComputeInstanceNotUsed
				}

				recommendation.ObjectType = model.ObjectTypeComputeInstance
			}

			recommendation.CompartmentID = inst.CompartmentID
			recommendation.CompartmentName = inst.CompartmentName
			recommendation.ResourceID = inst.ResourceID
			recommendation.Name = inst.Name
			detail1 := model.RecDetail{Name: "Instance Name", Value: inst.Name}
			detail2 := model.RecDetail{Name: "Cpu Core Count", Value: fmt.Sprintf("%.2f", allInstances[inst.ResourceID].OCPUs)}
			detail3 := model.RecDetail{Name: "%Cpu Average 90dd(daily)", Value: allInstanceMetrics[inst.ResourceID]["AvgCPU"].Perc}
			detail4 := model.RecDetail{Name: fmt.Sprintf("%%Cpu Average 90dd - Number of Threshold Reached (>%d%%)", percAvgCPU), Value: fmt.Sprintf("%s/%s", allInstanceMetrics[inst.ResourceID]["AvgCPU"].Value, allInstanceMetrics[inst.ResourceID]["AvgCPU"].Threshold)}

			var detail5, detail6 model.RecDetail
			if valPeak, ok := allInstanceMetrics[inst.ResourceID]["PeakCPU"]; ok {
				detail5 = model.RecDetail{Name: fmt.Sprintf("%%Cpu Average 7dd - Number of Threshold Reached (>%d%%)", percPeakCPU), Value: fmt.Sprintf("%s/%s", valPeak.Value, valPeak.Threshold)}
			} else {
				detail5 = model.RecDetail{Name: fmt.Sprintf("%%Cpu Average 7dd - Number of Threshold Reached (>%d%%)", percPeakCPU), Value: "NO DATA"}
			}

			if valMem, ok := allInstanceMetrics[inst.ResourceID]["AvgMemory"]; ok {
				detail6 = model.RecDetail{Name: fmt.Sprintf("%%memory Average 7dd - Number of Threshold Reached (>%d%%)", percMemoryUtilization), Value: fmt.Sprintf("%s/%s", valMem.Value, valMem.Threshold)}
			} else {
				detail6 = model.RecDetail{Name: fmt.Sprintf("%%memory Average 7dd - Number of Threshold Reached (>%d%%)", percMemoryUtilization), Value: "NO DATA"}
			}

			if recClusterName.Name != "" && recClusterName.Value != "" {
				recommendation.Details = append(recommendation.Details, detail1, detail2, recClusterName, detail3, detail4, detail5, detail6)
			} else {
				recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4, detail5, detail6)
			}

			listRec = append(listRec, recommendation)
		}
	}

	// build recommendation data for instances without monitoring
	// only if recommType is "rightsizing"
	if recommType == "rightsizing" {
		for _, in := range allInstancesWithoutMetrics {
			if in.Status != "STOPPED" {
				recommendation.Details = make([]model.RecDetail, 0)

				if in.Type == "kubernetes" {
					recClusterName = model.RecDetail{Name: "Oke Cluster Name", Value: allInstances[in.ResourceID].ClusterName}
					recommendation.Category = model.SISRightsizing1
					recommendation.Suggestion = model.ResizeOversizedKubernetesCluster
					recommendation.ObjectType = model.ObjectTypeClusterKubernetes
				} else {
					recommendation.Category = model.InstanceWithoutMonitoring
					recommendation.Suggestion = model.ResizeOversizedComputeInstance
					recommendation.ObjectType = model.ObjectTypeComputeInstance
				}

				recommendation.CompartmentID = in.CompartmentID
				recommendation.CompartmentName = in.CompartmentName
				recommendation.ResourceID = in.ResourceID
				recommendation.Name = in.Name

				detail1 := model.RecDetail{Name: "Instance Name", Value: in.Name}
				detail2 := model.RecDetail{Name: "Cpu Core Count", Value: fmt.Sprintf("%.2f", allInstances[in.ResourceID].OCPUs)}

				if recClusterName.Name != "" && recClusterName.Value != "" {
					recommendation.Details = append(recommendation.Details, detail1, detail2, recClusterName)
				} else {
					recommendation.Details = append(recommendation.Details, detail1, detail2)
				}

				listRec = append(listRec, recommendation)
			}
		}
	}

	return listRec, merr
}

func (as *ThunderService) countEventsOccurence(client monitoring.MonitoringClient, strQuery string, sTime common.SDKTime, eTime common.SDKTime, instances map[string]Instance, compartment model.OciCompartment, threshold int, percThreshold int, verifyShape bool, allInstanceMetrics map[string]MetricData, sType string) (map[string]Instance, error) {
	var metricData map[string]ValThresh

	var valThresh ValThresh

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
		totValue := 0.0
		maxValue := 0.0

		for _, a := range s.AggregatedDatapoints {
			totValue = totValue + *a.Value

			if maxValue > *a.Value {
				maxValue = *a.Value
			}

			if *a.Value > float64(percThreshold) {
				cnt++
			}
		}

		avgValue := totValue / float64(len(s.AggregatedDatapoints))
		valThresh.Perc = fmt.Sprintf("%.2f", avgValue)
		valThresh.Max = fmt.Sprintf("%.2f", maxValue)

		metricData = allInstanceMetrics[s.Dimensions["resourceId"]]

		if metricData == nil {
			metricData = make(map[string]ValThresh)
		}

		valThresh.Value = fmt.Sprintf("%d", cnt)
		valThresh.Threshold = fmt.Sprintf("%d", threshold)
		metricData[sType] = valThresh
		allInstanceMetrics[s.Dimensions["resourceId"]] = metricData

		if cnt > threshold {
			// the instance is not eligible for optimization
			if !verifyShape || (s.Dimensions["shape"] != "VM.Standard2.1" && s.Dimensions["shape"] != "VM.StandardE2.1") {
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

func (as *ThunderService) getOciInstancesList(allInstances map[string]Instance, compartment model.OciCompartment, customConfigProvider common.ConfigurationProvider, verifyShape bool) (map[string]Instance, error) {
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

		if !verifyShape || (*s.Shape != "VM.Standard2.1" && *s.Shape != "VM.StandardE2.1") {
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

			tmpInstance.Status = fmt.Sprintf("%v", s.LifecycleState)
			tmpInstance.OCPUs = *s.ShapeConfig.Ocpus
			allInstances[*s.Id] = tmpInstance
		}
	}

	return allInstances, nil
}

func (as *ThunderService) getOciInstancesWithMetrics(instances map[string]Instance, compartment model.OciCompartment, customConfigProvider common.ConfigurationProvider, verifyShape bool) (map[string]Instance, error) {
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
		if !verifyShape || (s.Dimensions["shape"] != "VM.StandardE2.1" && s.Dimensions["shape"] != "VM.Standard2.1") {
			// if the instance is not in the list I have to put it
			if _, ok := instances[s.Dimensions["resourceId"]]; !ok {
				tmpInstance := Instance{*s.CompartmentId, compartment.Name, s.Dimensions["resourceId"], s.Dimensions["resourceDisplayName"], "", s.Dimensions["shape"], 1, "", "", 0.0}
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
				resp, err := as.getMetricResponse(monClient, compartment.CompartmentID, "VolumeReadThroughput[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "VolumeWriteThroughput[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "VolumeReadOps[5d].max()")

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
				resp, err = as.getMetricResponse(monClient, compartment.CompartmentID, "VolumeWriteOps[5d].max()")

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
						isOpt, ociPerfs := as.isOptimizable(v)
						if isOpt {
							recommendation.Details = make([]model.RecDetail, 0)
							recommendation.Category = model.BlockStorageRightsizing
							recommendation.Suggestion = model.ResizeOversizedBlockStorage
							recommendation.CompartmentID = compartment.CompartmentID
							recommendation.CompartmentName = compartment.Name
							recommendation.ResourceID = v.ResourceID
							recommendation.Name = v.Name
							recommendation.ObjectType = model.ObjectTypeBlockStorage

							detail1 := model.RecDetail{Name: "Block Storage Name", Value: v.Name}
							detail2 := model.RecDetail{Name: "VPU", Value: fmt.Sprintf("%d", v.VpusPerGB)}
							detail3 := model.RecDetail{Name: "Size", Value: fmt.Sprintf("%d GB", v.Size)}
							detail4 := model.RecDetail{Name: "VPU Target", Value: fmt.Sprintf("%.0f MB/s %s%d iops", math.Round(ociPerfs.Performances[0].Values.MaxThroughput), " - ", ociPerfs.Performances[0].Values.MaxIOPS)}
							detail5 := model.RecDetail{Name: "Throughput R/W Max 5dd", Value: fmt.Sprintf("%.0f MB/s %s%.0f MB/s", math.Round(v.Throughput), " - ", math.Round(ociPerfs.Performances[0].Values.MaxThroughput))}
							detail6 := model.RecDetail{Name: "Iops Max 5dd", Value: fmt.Sprintf("%d%s%d", v.VpusPerGB, " / ", ociPerfs.Performances[0].Values.MaxIOPS)}

							recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4, detail5, detail6)

							listRec = append(listRec, recommendation)
						}
					}
				}
			}
		}
	}

	return listRec, merr
}

func (as *ThunderService) getMetricResponse(client monitoring.MonitoringClient, compartmentId string, query string) (*monitoring.SummarizeMetricsDataResponse, error) {
	var merr error

	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: &compartmentId,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Namespace: common.String("oci_blockstore"),
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

func (as *ThunderService) isOptimizable(res model.OciResourcePerformance) (bool, *model.OciVolumePerformance) {
	var ociPerfs *model.OciVolumePerformance

	if res.VpusPerGB == 0 {
		return false, nil
	}

	ociPerfs = as.getOciVolumePerformance(res.VpusPerGB, res.Size)

	if res.Throughput < (ociPerfs.Performances[0].Values.MaxThroughput/2.0) && res.Iops < (ociPerfs.Performances[0].Values.MaxIOPS)/2.0 {
		return true, ociPerfs
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

func (as *ThunderService) getOciInstances(customConfigProvider common.ConfigurationProvider, compartmentID string) (map[string]Instance, error) {
	var merr error

	retList := make(map[string]Instance)

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
		var tmpInstance Instance

		tmpInstance.CompartmentID = ""
		tmpInstance.CompartmentName = ""
		tmpInstance.ResourceID = *s.Id
		tmpInstance.Name = *s.DisplayName
		tmpInstance.Shape = *s.Shape

		if _, ok := s.Metadata["oke-pool-id"]; ok {
			tmpInstance.Type = "kubernetes"
			tmpInstance.ClusterName = s.Metadata["oke-cluster-display-name"]
		} else {
			tmpInstance.Type = "normal"
			tmpInstance.ClusterName = ""
		}

		tmpInstance.Status = fmt.Sprintf("%v", s.LifecycleState)
		tmpInstance.OCPUs = *s.ShapeConfig.Ocpus
		retList[*s.Id] = tmpInstance
	}

	return retList, nil
}
