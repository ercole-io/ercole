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

	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/core"

	"github.com/ercole-io/ercole/v2/model"
)

func (job *OciDataRetrieveJob) GetOciUnusedStorage(profiles []string, seqValue uint64) {
	var ore model.OciRecommendationError

	var volumeList map[string]model.OciVolume

	var attachedVolumeList []model.OciVolume

	var listRec []model.OciRecommendation

	var recommendation model.OciRecommendation

	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	volumeList, err := job.getOciVolumeList(profiles)
	if err != nil {
		recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
		errors = append(errors, recError)
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}

		return
	}

	attachedVolumeList, err = job.getOciAttachedVolumeList(profiles)
	if err != nil {
		recError := ore.SetOciRecommendationError(seqValue, "", model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
		errors = append(errors, recError)
		errDb := job.Database.AddOciRecommendationErrors(errors)

		if errDb != nil {
			job.Log.Error(errDb)
		}

		return
	}

	for _, avl := range attachedVolumeList {
		if val, ok := volumeList[avl.ResourceID]; ok {
			delete(volumeList, val.ResourceID)
		}
	}

	for _, vl := range volumeList {
		recommendation.Details = make([]model.RecDetail, 0)
		recommendation.SeqValue = seqValue
		recommendation.ProfileID = vl.ProfileID
		recommendation.Category = model.UnusedStorage
		recommendation.Suggestion = model.DeleteBlockStorageNotUsed
		recommendation.CompartmentID = vl.CompartmentID
		recommendation.CompartmentName = vl.CompartmentName
		recommendation.ResourceID = vl.ResourceID
		recommendation.Name = vl.Name
		recommendation.ObjectType = model.ObjectTypeBlockStorage
		detail1 := model.RecDetail{Name: "Block Storage Name", Value: vl.Name}
		detail2 := model.RecDetail{Name: "Size", Value: vl.Size}
		detail3 := model.RecDetail{Name: "Vpu", Value: vl.VpusPerGB}
		detail4 := model.RecDetail{Name: "Attached", Value: "No"}

		recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4)
		recommendation.CreatedAt = time.Now().UTC()
		listRec = append(listRec, recommendation)
	}

	if len(listRec) > 0 {
		errDb := job.Database.AddOciRecommendations(listRec)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}

func (job *OciDataRetrieveJob) getOciVolumeList(profiles []string) (map[string]model.OciVolume, error) {
	var merr error

	var listCompartments []model.OciCompartment

	var vol model.OciVolume

	var vols = make(map[string]model.OciVolume)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
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
					ProfileID:          profileId,
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

func (job *OciDataRetrieveJob) getOciAttachedVolumeList(profiles []string) ([]model.OciVolume, error) {
	var merr error

	var listCompartments []model.OciCompartment

	var vol model.OciVolume

	var vols []model.OciVolume

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
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
						ProfileID:          profileId,
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

func (job *OciDataRetrieveJob) GetOciOldSnapshotDecommissioning(profiles []string, seqValue uint64) {
	var ore model.OciRecommendationError

	var listCompartments []model.OciCompartment

	var recommendation model.OciRecommendation

	var listRec []model.OciRecommendation

	listRec = make([]model.OciRecommendation, 0)
	errors := make([]model.OciRecommendationError, 0)

	for _, profileId := range profiles {
		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(profileId)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		coreClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			recError := ore.SetOciRecommendationError(seqValue, profileId, model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
			errors = append(errors, recError)

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
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			nowt := common.SDKTime{Time: time.Now().Local()}

			for _, s := range resp.Items {
				tDiff := int(nowt.Sub(s.TimeCreated.Time).Hours() / 24)
				if s.SourceType == "MANUAL" && tDiff > 30 {
					recommendation.Details = make([]model.RecDetail, 0)
					recommendation.SeqValue = seqValue
					recommendation.ProfileID = profileId
					recommendation.Category = model.OldSnapshot
					recommendation.Suggestion = model.DeleteSnapshotOlder
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.ResourceID = *s.Id
					recommendation.Name = *s.DisplayName
					recommendation.ObjectType = model.ObjectTypeSnapshot
					detail1 := model.RecDetail{Name: "Snapshot Name", Value: *s.DisplayName}
					detail2 := model.RecDetail{Name: "Compartment Name", Value: compartment.Name}
					detail3 := model.RecDetail{Name: "Size", Value: fmt.Sprintf("%d GB", *s.SizeInGBs)}
					detail4 := model.RecDetail{Name: "Creation Date", Value: s.TimeCreated.String()}
					detail5 := model.RecDetail{Name: "Source Type", Value: "Manual"}

					recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4, detail5)
					recommendation.CreatedAt = time.Now().UTC()
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
				recError := ore.SetOciRecommendationError(seqValue, profileId, model.ObjectStorageOptimization, time.Now().UTC(), err.Error())
				errors = append(errors, recError)

				continue
			}

			for _, s := range resp1.Items {
				tDiff := int(nowt.Sub(s.TimeCreated.Time).Hours() / 24)
				if s.SourceType == "MANUAL" && tDiff > 30 {
					recommendation.Details = make([]model.RecDetail, 0)
					recommendation.SeqValue = seqValue
					recommendation.ProfileID = profileId
					recommendation.Category = model.OldSnapshot
					recommendation.Suggestion = model.DeleteSnapshotOlder
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.ResourceID = *s.Id
					recommendation.Name = *s.DisplayName
					recommendation.ObjectType = model.ObjectTypeSnapshot
					detail1 := model.RecDetail{Name: "Snapshot Name", Value: *s.DisplayName}
					detail2 := model.RecDetail{Name: "Compartment Name", Value: compartment.Name}
					detail3 := model.RecDetail{Name: "Size", Value: fmt.Sprintf("%d GB", *s.SizeInGBs)}
					detail4 := model.RecDetail{Name: "Creation Date", Value: s.TimeCreated.String()}
					detail5 := model.RecDetail{Name: "Source Type", Value: "Manual"}

					recommendation.Details = append(recommendation.Details, detail1, detail2, detail3, detail4, detail5)
					recommendation.CreatedAt = time.Now().UTC()
					listRec = append(listRec, recommendation)
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
