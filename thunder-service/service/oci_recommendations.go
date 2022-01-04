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
	"strconv"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/optimizer"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (as *ThunderService) GetOciRecommendations(profiles []string) ([]model.OciRecommendation, error) {
	var listRec []model.OciRecommendation
	listRec = make([]model.OciRecommendation, 0)

	dbProfiles, err := as.GetMapOciProfiles()
	if err != nil {
		return nil, err
	}

	var merr error

	for _, profileId := range profiles {

		objId, err := primitive.ObjectIDFromHex(profileId)
		if err != nil {
			merr = multierror.Append(merr, utils.NewErrorf("%w %q ", utils.ErrInvalidProfileId, profileId))
			continue
		}

		dbProfile, found := dbProfiles[objId]
		if !found {
			merr = multierror.Append(merr, utils.NewErrorf("%w: profileId %q", utils.ErrNotFound, profileId))
			continue
		}

		customConfigProvider := common.NewRawConfigurationProvider(dbProfile.TenancyOCID, dbProfile.UserOCID, dbProfile.Region, dbProfile.KeyFingerprint, *dbProfile.PrivateKey, nil)
		// Create a custom authentication provider that uses the profile passed as parameter
		// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
		client, err := optimizer.NewOptimizerClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		categoryId, err := as.getOciCategoryId(client, dbProfile.TenancyOCID)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		req := optimizer.ListRecommendationsRequest{
			CompartmentId: &dbProfile.TenancyOCID,
			CategoryId:    &categoryId,
			//CategoryId:    common.String("ocid1.optimizercategory.oc1..aaaaaaaa5w33jrqjheaxiczxxo2aguut6clwbn7aq2ujwjmyyzfz7b63uppq"),
			//CategoryId:             common.String("ocid1.optimizercategory.oc1..aaaaaaaaqeiskhuyp4pr7tohuooyujgyjmcq6cibc3btq6na62ev4ytz7ppa"),
			CompartmentIdInSubtree: common.Bool(true),
			Limit:                  common.Int(964),
		}

		resp, err := client.ListRecommendations(context.Background(), req)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		var cnt int
		var recTmp model.OciRecommendation

		if len(resp.Items) != 0 {
			for _, s := range resp.Items {
				for _, p := range s.ResourceCounts {
					if p.Status == "PENDING" {
						cnt = *p.Count
						break
					}
				}

				recTmp = model.OciRecommendation{
					TenancyOCID:         dbProfile.TenancyOCID,
					Name:                *s.Name,
					NumPending:          strconv.Itoa(cnt),
					EstimatedCostSaving: fmt.Sprintf("%.2f", *s.EstimatedCostSaving),
					Status:              fmt.Sprintf("%v", s.Status),
					Importance:          fmt.Sprintf("%v", s.Importance),
					RecommendationId:    *s.Id,
				}
				listRec = append(listRec, recTmp)
			}
		}
	}

	return listRec, merr
}

func (as *ThunderService) getOciCategoryId(optClient optimizer.OptimizerClient, tenancyOCID string) (string, error) {
	var merr error

	req := optimizer.ListCategoriesRequest{
		CompartmentId:          &tenancyOCID,
		CompartmentIdInSubtree: common.Bool(true),
	}
	// Send the request using the service client
	resp, err := optClient.ListCategories(context.Background(), req)
	if err != nil {
		merr := multierror.Append(merr, err)
		return "", merr
	}

	for _, s := range resp.Items {
		if *s.Name == "cost-management-name" {
			return *s.Id, nil
		}
	}

	return "", nil
}
