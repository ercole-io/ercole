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

	"github.com/ercole-io/ercole/v2/model"
	multierror "github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/objectstorage"
)

func (as *ThunderService) GetOciObjectStorageOptimization(profiles []string) ([]model.OciErcoleRecommendation, error) {
	var merr error

	var listCompartments []model.OciCompartment

	var recommendation model.OciErcoleRecommendation

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

		osClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		// retrieve buckets data for each compartment
		for _, compartment := range listCompartments {
			req1 := objectstorage.GetNamespaceRequest{
				CompartmentId: &compartment.CompartmentID,
			}
			resp1, err := osClient.GetNamespace(context.Background(), req1)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			req2 := objectstorage.ListBucketsRequest{
				CompartmentId: &compartment.CompartmentID,
				NamespaceName: common.String(*resp1.Value),
			}
			resp2, err := osClient.ListBuckets(context.Background(), req2)

			if err != nil {
				merr = multierror.Append(merr, err)
				continue
			}

			for _, bucket := range resp2.Items {
				req3 := objectstorage.GetBucketRequest{
					BucketName:    common.String(*bucket.Name),
					Fields:        []objectstorage.GetBucketFieldsEnum{objectstorage.GetBucketFieldsAutotiering, objectstorage.GetBucketFieldsApproximatecount, objectstorage.GetBucketFieldsApproximatesize},
					NamespaceName: common.String(*resp1.Value),
				}
				resp3, err := osClient.GetBucket(context.Background(), req3)

				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				if resp3.AutoTiering == "Disabled" {
					recommendation.Details = make([]model.RecDetail, 0)
					recommendation.Category = model.ObjectStorageOptimization
					recommendation.Suggestion = model.EnableBucketAutoTiering
					recommendation.CompartmentID = compartment.CompartmentID
					recommendation.CompartmentName = compartment.Name
					recommendation.ResourceID = *resp3.Id
					recommendation.Name = *resp3.Name
					recommendation.ObjectType = "Object Storage"
					detail1 := model.RecDetail{Name: "Bucket Name", Value: *resp3.Name}
					detail2 := model.RecDetail{Name: "Size", Value: as.getBucketSize(*resp3.ApproximateSize)}
					detail3 := model.RecDetail{Name: "Optimization", Value: "Enable auto-tiering"}

					recommendation.Details = append(recommendation.Details, detail1, detail2, detail3)
					listRec = append(listRec, recommendation)
				}
			}
		}
	}

	return listRec, nil
}

func (as *ThunderService) getBucketSize(sizeVal int64) string {
	var valRet string

	var newVal float64

	var valTmp float64

	cnt := 0
	newVal = float64(sizeVal)

	for {
		valTmp = newVal
		newVal = newVal / 1024.0
		cnt = cnt + 1

		if newVal <= 1 {
			switch cnt {
			case 1:
				valRet = fmt.Sprintf("%d bytes", sizeVal)
			case 2:
				valRet = fmt.Sprintf("%.2f KiB", valTmp)
			case 3:
				valRet = fmt.Sprintf("%.2f MiB", valTmp)
			case 4:
				valRet = fmt.Sprintf("%.2f GiB", valTmp)
			case 5:
				valRet = fmt.Sprintf("%.2f TiB", valTmp)
			case 6:
				valRet = fmt.Sprintf("%.2f PiB", valTmp)
			default:
				valRet = fmt.Sprintf("%d bytes", sizeVal)
			}

			break
		}
	}

	return valRet
}
