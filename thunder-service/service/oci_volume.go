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
	"time"

	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/core"

	"github.com/ercole-io/ercole/v2/model"
)

func (as *ThunderService) GetOciUnusedStorage(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var merr error
	var err error
	var volumeList map[string]model.OciVolume
	var attachedVolumeList []model.OciVolume
	var listRec []model.OciErcoleRecommendation
	var recommendation model.OciErcoleRecommendation

	listRec = make([]model.OciErcoleRecommendation, 0)

	volumeList, err = as.GetOciVolumeList(profiles)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	attachedVolumeList, err = as.GetOciAttachedVolumeList(profiles)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	for _, avl := range attachedVolumeList {
		if val, ok := volumeList[avl.ResourceID]; ok {
			delete(volumeList, val.ResourceID)
		}
	}

	for _, vl := range volumeList {
		recommendation.CompartmentID = vl.CompartmentID
		recommendation.CompartmentName = vl.CompartmentName
		recommendation.ResourceID = vl.ResourceID
		recommendation.Name = vl.Name
		listRec = append(listRec, recommendation)
	}

	return listRec, merr
}

func (as *ThunderService) GetOciVolumeList(profiles []string) (map[string]model.OciVolume, error) {
	var merr error
	var listCompartments []model.OciCompartment
	var vol model.OciVolume
	var vols = make(map[string]model.OciVolume)

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

		// retrieve volume data for each compartment
		for _, compartment := range listCompartments {

			coreClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			req := core.ListVolumesRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp, err := coreClient.ListVolumes(context.Background(), req)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			for _, r := range resp.Items {
				vol = model.OciVolume{
					CompartmentID:      compartment.CompartmentID,
					CompartmentName:    compartment.Name,
					ResourceID:         *r.Id,
					Name:               *r.DisplayName,
					Size:               fmt.Sprintf("%d", *r.SizeInGBs),
					VpusPerGB:          fmt.Sprintf("%d", *r.VpusPerGB),
					AvailabilityDomain: *r.AvailabilityDomain,
					State:              fmt.Sprintf("%v", r.LifecycleState),
				}
				vols[*r.Id] = vol
			}
		}
	}
	return vols, merr
}

func (as *ThunderService) GetOciAttachedVolumeList(profiles []string) ([]model.OciVolume, error) {
	var merr error
	var listCompartments []model.OciCompartment
	var vol model.OciVolume
	var vols []model.OciVolume

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

		// retrieve attached volume data for each compartment
		for _, compartment := range listCompartments {

			coreClient, err := core.NewComputeClientWithConfigurationProvider(customConfigProvider)
			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			req := core.ListVolumeAttachmentsRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			resp, err := coreClient.ListVolumeAttachments(context.Background(), req)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}
			for _, r := range resp.Items {
				if fmt.Sprintf("%v", r.GetLifecycleState()) == "ATTACHED" {
					vol = model.OciVolume{
						CompartmentID:      compartment.CompartmentID,
						CompartmentName:    compartment.Name,
						ResourceID:         *r.GetVolumeId(),
						Name:               *r.GetDisplayName(),
						Size:               "",
						VpusPerGB:          "",
						AvailabilityDomain: *r.GetAvailabilityDomain(),
						State:              fmt.Sprintf("%v", r.GetLifecycleState()),
					}
					vols = append(vols, vol)
				}
			}
		}
	}
	return vols, merr
}

func (as *ThunderService) GetOciOldSnapshotDecommissioning(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var merr error
	var listCompartments []model.OciCompartment
	var recommendation model.OciErcoleRecommendation
	var listRec []model.OciErcoleRecommendation

	listRec = make([]model.OciErcoleRecommendation, 0)

	for _, profileId := range profiles {

		customConfigProvider, tenancyOCID, err := as.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		coreClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			return nil, merr
		}

		listCompartments, err = as.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		for _, compartment := range listCompartments {

			// request for list volume backup
			req := core.ListVolumeBackupsRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			// Send the request using the service client
			resp, err := coreClient.ListVolumeBackups(context.Background(), req)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			nowt := common.SDKTime{Time: time.Now().Local()}

			for _, s := range resp.Items {
				tDiff := int(nowt.Sub(s.TimeCreated.Time) / 24)
				if s.SourceType == "MANUAL" && tDiff > 30 {
					recommendation.Type = model.RecommendationTypeOldSnapshot
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.ResourceID = *s.Id
					recommendation.Name = *s.DisplayName

					listRec = append(listRec, recommendation)
				}
			}

			// request for list boot volume backup
			req1 := core.ListBootVolumeBackupsRequest{
				CompartmentId: &compartment.CompartmentID,
			}

			// Send the request using the service client
			resp1, err := coreClient.ListBootVolumeBackups(context.Background(), req1)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			for _, s := range resp1.Items {
				tDiff := int(nowt.Sub(s.TimeCreated.Time) / 24)
				if s.SourceType == "MANUAL" && tDiff > 30 {
					recommendation.Type = model.RecommendationTypeOldSnapshot
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.ResourceID = *s.Id
					recommendation.Name = *s.DisplayName

					listRec = append(listRec, recommendation)
				}
			}
		}

	}
	//	return volumeList, merr
	return listRec, merr
}
