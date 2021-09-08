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

	"github.com/ercole-io/ercole/v2/model"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/optimizer"
)

func (as *ThunderService) GetOCRecommendations(compartmentId string) ([]model.Recommendation, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())

	if err != nil {
		return nil, err
	}
	req := optimizer.ListRecommendationsRequest{
		CompartmentId:          &compartmentId,
		CategoryId:             common.String("ocid1.optimizercategory.oc1..aaaaaaaaqeiskhuyp4pr7tohuooyujgyjmcq6cibc3btq6na62ev4ytz7ppa"),
		CompartmentIdInSubtree: common.Bool(true),
		Limit:                  common.Int(964),
	}

	resp, err := client.ListRecommendations(context.Background(), req)
	if err != nil {
		return nil, err
	}
	var cnt int
	var recTmp model.Recommendation
	var listRec []model.Recommendation

	for _, s := range resp.Items {
		for _, p := range s.ResourceCounts {
			if p.Status == "PENDING" {
				cnt = *p.Count
			}
			recTmp = model.Recommendation{*s.Name, strconv.Itoa(cnt), fmt.Sprintf("%.2f", *s.EstimatedCostSaving), fmt.Sprintf("%v", s.Status), fmt.Sprintf("%v", s.Importance), *s.Id}

		}
		listRec = append(listRec, recTmp)
	}
	return listRec, nil
}

func (as *ThunderService) GetOCRecommendationsWithCategory(compartmentId string) ([]model.RecommendationWithCategory, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())

	if err != nil {
		return nil, err
	}
	var listRecWithCat []model.RecommendationWithCategory
	listCategory, err := GetOCListCategories(compartmentId)

	if err != nil {
		return nil, err
	}
	for _, q := range listCategory {
		req := optimizer.ListRecommendationsRequest{
			CompartmentId:          &compartmentId,
			CategoryId:             &q.CategoryId,
			CompartmentIdInSubtree: common.Bool(true),
		}

		resp, err := client.ListRecommendations(context.Background(), req)
		if err != nil {
			return nil, err
		}
		var cnt int
		var recTmp model.Recommendation
		var recWithCatTmp model.RecommendationWithCategory
		var listRec []model.Recommendation

		for _, s := range resp.Items {
			for _, p := range s.ResourceCounts {
				if p.Status == "PENDING" {
					cnt = *p.Count
				}
				recTmp = model.Recommendation{*s.Name, strconv.Itoa(cnt), fmt.Sprintf("%.2f", *s.EstimatedCostSaving), fmt.Sprintf("%v", s.Status), fmt.Sprintf("%v", s.Importance), *s.Id}

			}
			listRec = append(listRec, recTmp)
		}
		recWithCatTmp = model.RecommendationWithCategory{q.Name, listRec}
		listRecWithCat = append(listRecWithCat, recWithCatTmp)

	}
	return listRecWithCat, nil
}

func GetOCListCategories(compartmentId string) ([]model.Category, error) {
	// Create a default authentication provider that uses the DEFAULT
	// profile in the configuration file.
	// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.
	client, err := optimizer.NewOptimizerClientWithConfigurationProvider(common.DefaultConfigProvider())

	if err != nil {
		return nil, err
	}
	req := optimizer.ListCategoriesRequest{
		CompartmentId:          &compartmentId,
		CompartmentIdInSubtree: common.Bool(true),
	}

	resp, err := client.ListCategories(context.Background(), req)

	if err != nil {
		return nil, err
	}
	var catTmp model.Category
	var listCategory []model.Category

	for _, s := range resp.Items {
		catTmp = model.Category{*s.Name, *s.Id}
		listCategory = append(listCategory, catTmp)
	}
	return listCategory, nil
}
