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
	"github.com/oracle/oci-go-sdk/identity"
	"github.com/oracle/oci-go-sdk/v45/common"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

func (as *ThunderService) GetOciCompartments(profiles []string) ([]model.OciCompartment, error) {
	var merr error
	var listCompartments []model.OciCompartment

	// retrieve data for configurated profiles
	dbProfiles, err := as.GetMapOciProfiles()
	if err != nil {
		return nil, err
	}

	// retrieve compartments list for all selected profiles
	for _, profileId := range profiles {

		objId, err := primitive.ObjectIDFromHex(profileId)
		if err != nil {
			merr = multierror.Append(merr, utils.NewErrorf("%w - invalid profileId %q", err, profileId))
			continue
		}

		dbProfile, found := dbProfiles[objId]
		if !found {
			merr = multierror.Append(merr, utils.NewErrorf("profile %q not found", profileId))
			continue
		}

		customConfigProvider := common.NewRawConfigurationProvider(dbProfile.TenancyOCID, dbProfile.UserOCID, dbProfile.Region, dbProfile.KeyFingerprint, *dbProfile.PrivateKey, nil)
		// Create a custom authentication provider that uses the profile passed as parameter
		// Refer to <see href="https://docs.cloud.oracle.com/en-us/iaas/Content/API/Concepts/sdkconfig.htm#SDK_and_CLI_Configuration_File>the public documentation</see> on how to prepare a configuration file.

		client, err := identity.NewIdentityClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		req := identity.ListCompartmentsRequest{
			CompartmentId: &dbProfile.TenancyOCID,
		}

		resp, err := client.ListCompartments(context.Background(), req)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		var compTmp model.OciCompartment
		for _, s := range resp.Items {
			compTmp.CompartmentID = *s.Id
			compTmp.Name = *s.Name
			compTmp.Description = *s.Description
			compTmp.TimeCreating = s.TimeCreated.String()
			listCompartments = append(listCompartments, compTmp)
		}
	}
	return listCompartments, merr
}

func (as *ThunderService) getOciProfileCompartments(tenancyOCID string, customConfigProvider common.ConfigurationProvider) ([]model.OciCompartment, error) {
	var merr error
	var listCompartments []model.OciCompartment

	client, err := identity.NewIdentityClientWithConfigurationProvider(customConfigProvider)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	req := identity.ListCompartmentsRequest{
		CompartmentId: &tenancyOCID,
	}

	resp, err := client.ListCompartments(context.Background(), req)
	if err != nil {
		merr = multierror.Append(merr, err)
		return nil, merr
	}

	var compTmp model.OciCompartment
	for _, s := range resp.Items {
		compTmp.CompartmentID = *s.Id
		compTmp.Name = *s.Name
		compTmp.Description = *s.Description
		compTmp.TimeCreating = s.TimeCreated.String()
		listCompartments = append(listCompartments, compTmp)
	}

	return listCompartments, merr
}

func (as *ThunderService) getOciCustomConfigProviderAndTenancy(profileId string) (common.ConfigurationProvider, string, error) {
	var merr error

	// retrieve data for configurated profiles
	dbProfiles, err := as.GetMapOciProfiles()
	if err != nil {
		return nil, "", err
	}

	objId, err := primitive.ObjectIDFromHex(profileId)
	if err != nil {
		merr = multierror.Append(merr, utils.NewErrorf("%w - invalid profileId %q", err, profileId))
		return nil, "", merr
	}

	dbProfile, found := dbProfiles[objId]
	if !found {
		merr = multierror.Append(merr, utils.NewErrorf("profile %q not found", profileId))
		return nil, "", merr
	}

	customConfigProvider := common.NewRawConfigurationProvider(dbProfile.TenancyOCID, dbProfile.UserOCID, dbProfile.Region, dbProfile.KeyFingerprint, *dbProfile.PrivateKey, nil)

	return customConfigProvider, dbProfile.TenancyOCID, nil
}
