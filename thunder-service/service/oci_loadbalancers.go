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

	multierror "github.com/hashicorp/go-multierror"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/oracle/oci-go-sdk/v45/loadbalancer"
)

func (as *ThunderService) GetOciUnusedLoadBalancers(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var merr error
	var listCompartments []model.OciCompartment
	var recommendation model.OciErcoleRecommendation

	tempListRec := make(map[string]model.OciErcoleRecommendation, 0)
	listRec := make([]model.OciErcoleRecommendation, 0)

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

		lbClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			return nil, merr
		}

		// retrieve load balancer data for each compartment
		for _, compartment := range listCompartments {

			req := loadbalancer.ListLoadBalancerHealthsRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp, err := lbClient.ListLoadBalancerHealths(context.Background(), req)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			for _, s := range resp.Items {
				if s.Status == "CRITICAL" || s.Status == "UNKNOWN" {
					recommendation.Type = model.RecommendationTypeUnusedResource
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.Name = ""
					recommendation.ResourceID = *s.LoadBalancerId
					tempListRec[*s.LoadBalancerId] = recommendation
					//listRec = append(listRec, recommendation)
				}
			}

			req1 := loadbalancer.ListLoadBalancersRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp1, err := lbClient.ListLoadBalancers(context.Background(), req1)

			if err != nil {
				merr = multierror.Append(merr, err)
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
	return listRec, merr
}
