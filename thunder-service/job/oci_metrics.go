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
package job

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/core"
	"github.com/oracle/oci-go-sdk/v45/monitoring"

	"github.com/ercole-io/ercole/v2/model"
)

type Instance struct {
	CompartmentID   string  `json:"compartmentID"`
	CompartmentName string  `json:"compartmentName"`
	ProfileID       string  `json:"profileID"`
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

func (job *OciDataRetrieveJob) GetOciComputeInstancesIdle(profiles []string, seqValue uint64) {
	var listRec []model.OciRecommendation

	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		// retrieve metrics data for each compartment
		var strNamespace = "oci_compute_infrastructure_health"

		for _, compartment := range listCompartments {
			instances, err := job.getOciInstances(customConfigProvider, compartment.CompartmentID, profileId)
			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

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
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

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
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			var recommendation model.OciRecommendation

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
				recommendation.ProfileID = profileId
				recommendation.SeqValue = seqValue

				if value.Type == "kubernetes" {
					recommendation.Category = model.OciUnusedServiceDecommisioning
					recommendation.Suggestion = model.OciDeleteKubernetesNodeNotActive
					recommendation.ObjectType = model.OciObjectTypeClusterKubernetes
				} else {
					recommendation.Category = model.OciComputeInstanceIdle
					recommendation.Suggestion = model.OciDeleteComputeInstanceNotActive
					recommendation.ObjectType = model.OciObjectTypeComputeInstance
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

				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddOciRecommendations(listRec)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}

	if len(errors) > 0 {
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}

func (job *OciDataRetrieveJob) GetOciComputeInstanceRightsizing(profiles []string, seqValue uint64) {
	job.getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles, 50, 50, 90, true, "rightsizing", seqValue)
}

func (job *OciDataRetrieveJob) GetOciUnusedServiceDecommisioning(profiles []string, seqValue uint64) {
	job.getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles, 5, 5, 40, false, "decommisioning", seqValue)
}

func (job *OciDataRetrieveJob) getOciDataForCoumputeInstanceAndServiceDecommisioning(profiles []string, percAvgCPU int, percPeakCPU int, percMemoryUtilization int, verifyShape bool, recommType string, seqValue uint64) {
	var listRec []model.OciRecommendation

	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	var AvgCPUThreshold = 3

	var PeakCPUThreshold = 180

	var MemoryThreshold = 1

	var recClusterName model.RecDetail

	instancesNotOptimizable := make(map[string]Instance)
	allInstancesWithMetrics := make(map[string]Instance)
	allInstances := make(map[string]Instance)

	allInstanceMetrics := make(map[string]MetricData)

	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		// retrieve metrics data for each compartment
		for _, compartment := range listCompartments {
			allInstances, err = job.getOciInstancesList(allInstances, compartment, profileId, customConfigProvider, verifyShape)
			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			allInstancesWithMetrics, err = job.getOciInstancesWithMetrics(allInstances, allInstancesWithMetrics, compartment, profileId, customConfigProvider, verifyShape)
			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			// if an instance with metrics is not in the list of all instances I have to remove it
			for _, b := range allInstancesWithMetrics {
				if _, ok := allInstances[b.ResourceID]; !ok {
					delete(allInstancesWithMetrics, b.ResourceID)
				}
			}

			for _, inst := range allInstancesWithMetrics {
				// first query is about average CPU utilization in the last  90 days
				var strQueryAvgCPU = fmt.Sprintf("CpuUtilization[1d]{resourceId=\"%s\"}.avg()", inst.ResourceID)

				var sTime common.SDKTime

				var eTime common.SDKTime

				sTime.Time = time.Now().Local().AddDate(0, 0, -89)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = job.countEventsOccurence(allInstances, monClient, strQueryAvgCPU, sTime, eTime, instancesNotOptimizable, compartment, profileId, AvgCPUThreshold, percAvgCPU, verifyShape, allInstanceMetrics, "AvgCPU")
				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)
				}

				// second query is about CPU utilization peak in the last 7 days
				var strQueryPeakCPU = fmt.Sprintf("CpuUtilization[1m]{resourceId=\"%s\"}.max()", inst.ResourceID)

				sTime.Time = time.Now().Local().AddDate(0, 0, -7)
				sTime.Time = sTime.Add(time.Hour * 1)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = job.countEventsOccurence(allInstances, monClient, strQueryPeakCPU, sTime, eTime, instancesNotOptimizable, compartment, profileId, PeakCPUThreshold, percPeakCPU, verifyShape, allInstanceMetrics, "PeakCPU")
				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)
				}

				// third query is about memory utilization in the last 7 days
				var strQueryMemory = fmt.Sprintf("MemoryUtilization[1m]{resourceId=\"%s\"}.max()", inst.ResourceID)

				sTime.Time = time.Now().Local().AddDate(0, 0, -7)
				eTime.Time = time.Now().Local()

				instancesNotOptimizable, err = job.countEventsOccurence(allInstances, monClient, strQueryMemory, sTime, eTime, instancesNotOptimizable, compartment, profileId, MemoryThreshold, percMemoryUtilization, verifyShape, allInstanceMetrics, "AvgMemory")
				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)
				}
			}
		}
	}

	var recommendation model.OciRecommendation

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

			recommendation.ProfileID = inst.ProfileID
			recommendation.SeqValue = seqValue

			if allInstances[inst.ResourceID].Type == "kubernetes" {
				recClusterName = model.RecDetail{Name: "Oke Cluster Name", Value: allInstances[inst.ResourceID].ClusterName}

				if recommType == "rightsizing" {
					recommendation.Category = model.OciSISRightsizing1
					recommendation.Suggestion = model.OciResizeOversizedKubernetesCluster
				} else {
					recommendation.Category = model.OciUnusedServiceDecommisioning
					recommendation.Suggestion = model.OciDeleteKubernetesNodeNotUsed
				}

				recommendation.ObjectType = model.OciObjectTypeClusterKubernetes
			} else {
				if recommType == "rightsizing" {
					recommendation.Category = model.OciInstanceRightsizing
					recommendation.Suggestion = model.OciResizeOversizedComputeInstance
				} else {
					recommendation.Category = model.OciComputeInstanceDecommisioning
					recommendation.Suggestion = model.OciDeleteComputeInstanceNotUsed
				}

				recommendation.ObjectType = model.OciObjectTypeComputeInstance
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

			recommendation.CreatedAt = time.Now().UTC()

			listRec = append(listRec, recommendation)
		}
	}

	// build recommendation data for instances without monitoring
	// only if recommType is "rightsizing"
	if recommType == "rightsizing" {
		for _, in := range allInstancesWithoutMetrics {
			if in.Status != "STOPPED" {
				recommendation.Details = make([]model.RecDetail, 0)

				recommendation.ProfileID = in.ProfileID
				recommendation.SeqValue = seqValue

				if in.Type == "kubernetes" {
					recClusterName = model.RecDetail{Name: "Oke Cluster Name", Value: allInstances[in.ResourceID].ClusterName}
					recommendation.Category = model.OciSISRightsizing1
					recommendation.Suggestion = model.OciResizeOversizedKubernetesCluster
					recommendation.ObjectType = model.OciObjectTypeClusterKubernetes
				} else {
					recommendation.Category = model.OciInstanceWithoutMonitoring
					recommendation.Suggestion = model.OciResizeOversizedComputeInstance
					recommendation.ObjectType = model.OciObjectTypeComputeInstance
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

				recommendation.CreatedAt = time.Now().UTC()

				listRec = append(listRec, recommendation)
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddOciRecommendations(listRec)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}

	if len(errors) > 0 {
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}

func (job *OciDataRetrieveJob) countEventsOccurence(allInstances map[string]Instance, client monitoring.MonitoringClient, strQuery string, sTime common.SDKTime, eTime common.SDKTime, instancesNotOptimizable map[string]Instance, compartment model.OciCompartment, profileId string, threshold int, percThreshold int, verifyShape bool, allInstanceMetrics map[string]MetricData, sType string) (map[string]Instance, error) {
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
		return instancesNotOptimizable, err
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
			if !verifyShape || (allInstances[s.Dimensions["resourceId"]].OCPUs > 1) {
				if val, ok := instancesNotOptimizable[s.Dimensions["resourceId"]]; ok {
					val.Cnt += 1
					instancesNotOptimizable[s.Dimensions["resourceId"]] = val
				} else {
					instance.CompartmentID = compartment.CompartmentID
					instance.CompartmentName = compartment.Name
					instance.ProfileID = profileId
					instance.ResourceID = s.Dimensions["resourceId"]
					instance.Name = s.Dimensions["resourceDisplayName"]
					instance.Shape = s.Dimensions["shape"]
					instancesNotOptimizable[s.Dimensions["resourceId"]] = instance
				}
			}
		}
	}

	return instancesNotOptimizable, nil
}

func (job *OciDataRetrieveJob) getOciInstancesList(allInstances map[string]Instance, compartment model.OciCompartment, profileId string, customConfigProvider common.ConfigurationProvider, verifyShape bool) (map[string]Instance, error) {
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

		if !verifyShape || (*s.ShapeConfig.Ocpus > 1) {
			tmpInstance.CompartmentID = compartment.CompartmentID
			tmpInstance.CompartmentName = compartment.Name
			tmpInstance.ProfileID = profileId
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

func (job *OciDataRetrieveJob) getOciInstancesWithMetrics(allInstances map[string]Instance, allInstancesWithMetrics map[string]Instance, compartment model.OciCompartment, profileId string, customConfigProvider common.ConfigurationProvider, verifyShape bool) (map[string]Instance, error) {
	client, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)

	if err != nil {
		return allInstancesWithMetrics, err
	}

	req := monitoring.ListMetricsRequest{
		CompartmentId:      &compartment.CompartmentID,
		ListMetricsDetails: monitoring.ListMetricsDetails{Namespace: common.String("oci_computeagent")},
	}

	// Send the request using the service client
	resp, err := client.ListMetrics(context.Background(), req)

	if err != nil {
		return allInstancesWithMetrics, err
	}

	for _, s := range resp.Items {
		if !verifyShape || (allInstances[s.Dimensions["resourceId"]].OCPUs > 1) {
			// if the instance is not in the list I have to put it
			if _, ok := allInstancesWithMetrics[s.Dimensions["resourceId"]]; !ok {
				tmpInstance := Instance{*s.CompartmentId, compartment.Name, profileId, s.Dimensions["resourceId"], s.Dimensions["resourceDisplayName"], "", s.Dimensions["shape"], 1, "", "", 0.0}
				allInstancesWithMetrics[s.Dimensions["resourceId"]] = tmpInstance
			}
		}
	}

	return allInstancesWithMetrics, nil
}

func (job *OciDataRetrieveJob) GetOciBlockStorageRightsizing(profiles []string, seqValue uint64) {
	var listRec []model.OciRecommendation

	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	var recommendation model.OciRecommendation

	var vol model.OciResourcePerformance

	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)

		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)

		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		var resTmp model.OciResourcePerformance

		var ok bool

		monClient, err := monitoring.NewMonitoringClientWithConfigurationProvider(customConfigProvider)

		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)
			errDb := job.Database.AddOciRecommendationErrors(errors)

			if errDb != nil {
				job.Log.Error(errDb)
			}

			return
		}

		// retrieve metrics data for each compartment
		for _, compartment := range listCompartments {
			var vols = make(map[string]model.OciResourcePerformance)

			coreClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)

			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			req := core.ListVolumesRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp1, err := coreClient.ListVolumes(context.Background(), req)

			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

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
				resp, err := job.getMetricResponse(monClient, compartment.CompartmentID, "VolumeReadThroughput[5d].max()")

				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)

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
				resp, err = job.getMetricResponse(monClient, compartment.CompartmentID, "VolumeWriteThroughput[5d].max()")

				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)

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
				resp, err = job.getMetricResponse(monClient, compartment.CompartmentID, "VolumeReadOps[5d].max()")

				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)

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
				resp, err = job.getMetricResponse(monClient, compartment.CompartmentID, "VolumeWriteOps[5d].max()")

				if err != nil {
					recError := ore.SetOciRecommendationError(seqValue, profileId, model.OciObjectStorageOptimization, time.Now().UTC(), err.Error())
					errors = append(errors, recError)

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
						isOpt, ociPerfs := job.isOptimizable(v)
						if isOpt {
							recommendation.Details = make([]model.RecDetail, 0)
							recommendation.SeqValue = seqValue
							recommendation.ProfileID = profileId
							recommendation.Category = model.OciBlockStorageRightsizing
							recommendation.Suggestion = model.OciResizeOversizedBlockStorage
							recommendation.CompartmentID = compartment.CompartmentID
							recommendation.CompartmentName = compartment.Name
							recommendation.ResourceID = v.ResourceID
							recommendation.Name = v.Name
							recommendation.ObjectType = model.OciObjectTypeBlockStorage

							detail1 := model.RecDetail{Name: "Block Storage Name", Value: v.Name}
							detail2 := model.RecDetail{Name: "VPU", Value: fmt.Sprintf("%d", v.VpusPerGB)}
							detail3 := model.RecDetail{Name: "Size", Value: fmt.Sprintf("%d GB", v.Size)}
							detail4 := model.RecDetail{Name: "VPU Target", Value: fmt.Sprintf("%.0f MB/s %s%d iops", math.Round(ociPerfs.Performances[0].Values.MaxThroughput), " - ", ociPerfs.Performances[0].Values.MaxIOPS)}
							detail5 := model.RecDetail{Name: "Throughput R/W Max 5dd", Value: fmt.Sprintf("%.0f MB/s %s%.0f MB/s", math.Round(v.Throughput), " - ", math.Round(ociPerfs.Performances[0].Values.MaxThroughput))}
							detail6 := model.RecDetail{Name: "Iops Max 5dd", Value: fmt.Sprintf("%d%s%d", v.VpusPerGB, " / ", ociPerfs.Performances[0].Values.MaxIOPS)}

							recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4, detail5, detail6)
							recommendation.CreatedAt = time.Now().UTC()

							listRec = append(listRec, recommendation)
						}
					}
				}
			}
		}
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddOciRecommendations(listRec)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}

	if len(errors) > 0 {
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}

func (job *OciDataRetrieveJob) getMetricResponse(client monitoring.MonitoringClient, compartmentId string, query string) (*monitoring.SummarizeMetricsDataResponse, error) {
	req := monitoring.SummarizeMetricsDataRequest{
		CompartmentId: &compartmentId,
		SummarizeMetricsDataDetails: monitoring.SummarizeMetricsDataDetails{
			Namespace: common.String("oci_blockstore"),
			Query:     &query,
		},
	}

	resp, err := client.SummarizeMetricsData(context.Background(), req)

	if err != nil {
		return nil, err
	}

	return &resp, nil
}

func (job *OciDataRetrieveJob) isOptimizable(res model.OciResourcePerformance) (bool, *model.OciVolumePerformance) {
	var ociPerfs *model.OciVolumePerformance

	if res.VpusPerGB == 0 {
		return false, nil
	}

	ociPerfs = job.getOciVolumePerformance(res.VpusPerGB, res.Size)

	if res.Throughput < (ociPerfs.Performances[0].Values.MaxThroughput/2.0) && res.Iops < (ociPerfs.Performances[0].Values.MaxIOPS)/2.0 {
		return true, ociPerfs
	} else {
		return false, nil
	}
}

func (job *OciDataRetrieveJob) getOciVolumePerformance(vpu int, size int) *model.OciVolumePerformance {
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

func (job *OciDataRetrieveJob) getOciInstances(customConfigProvider common.ConfigurationProvider, compartmentID string, profileId string) (map[string]Instance, error) {
	retList := make(map[string]Instance)

	client, err := core.NewComputeClientWithConfigurationProvider(customConfigProvider)

	if err != nil {
		return nil, err
	}

	req := core.ListInstancesRequest{
		CompartmentId: &compartmentID,
	}

	resp, err := client.ListInstances(context.Background(), req)

	if err != nil {
		return nil, err
	}

	for _, s := range resp.Items {
		var tmpInstance Instance

		tmpInstance.CompartmentID = compartmentID
		tmpInstance.CompartmentName = ""
		tmpInstance.ProfileID = profileId
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
