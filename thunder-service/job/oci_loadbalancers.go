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
	"time"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/oracle/oci-go-sdk/v45/loadbalancer"
)

func (job *OciDataRetrieveJob) GetOciUnusedLoadBalancers(profiles []string, seqValue uint64) {
	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	var recommendation model.OciRecommendation

	tempListRec := make(map[string]model.OciRecommendation, 0)
	listRec := make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)

		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)
			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)

		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)
			continue
		}

		lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)
			continue
		}

		// retrieve load balancer data for each compartment
		for _, compartment := range listCompartments {
			req := loadbalancer.ListLoadBalancerHealthsRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp, err := lbClient.ListLoadBalancerHealths(context.Background(), req)

			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)
				continue
			}

			for _, s := range resp.Items {
				if s.Status == "CRITICAL" || s.Status == "UNKNOWN" {
					recommendation.Details = make([]model.RecDetail, 0)
					recommendation.SeqValue = seqValue
					recommendation.ProfileID = profileId
					recommendation.Category = model.UnusedResource
					recommendation.Suggestion = model.DeleteLoadBalancerNotActive
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.Name = ""
					recommendation.ResourceID = *s.LoadBalancerId
					recommendation.ObjectType = model.ObjectTypeLoadBalancer
					detail1 := model.RecDetail{Name: "Resource Id", Value: *s.LoadBalancerId}
					detail2 := model.RecDetail{Name: "Resource Type", Value: "Load Balancer"}
					detail3 := model.RecDetail{Name: "Resource Status", Value: fmt.Sprintf("%v", s.Status)}

					recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)
					recommendation.CreatedAt = time.Now().UTC()
					tempListRec[*s.LoadBalancerId] = recommendation
				}
			}

			req1 := loadbalancer.ListLoadBalancersRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp1, err := lbClient.ListLoadBalancers(context.Background(), req1)

			if err != nil {
				recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)
				continue
			}

			for _, r := range resp1.Items {
				if rec, ok := tempListRec[*r.Id]; ok {
					rec.Name = *r.DisplayName
					listRec = append(listRec, rec)
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
