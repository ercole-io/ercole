// Copyright (c) 2021 Sorint.lab S.p.A.
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

package job

import (
	"context"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	db "github.com/ercole-io/ercole/v2/thunder-service/database"
	"github.com/hashicorp/go-multierror"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/oracle/oci-go-sdk/database"
	"github.com/oracle/oci-go-sdk/filestorage"
	"github.com/oracle/oci-go-sdk/identity"
	"github.com/oracle/oci-go-sdk/v45/common"
	"github.com/oracle/oci-go-sdk/v45/loadbalancer"
	"github.com/oracle/oci-go-sdk/v45/objectstorage"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OciDataRetrieveJob is the job used to retrieve data from Oracle Cloud installation
type OciDataRetrieveJob struct {
	// Database contains the database layer
	Database db.MongoDatabaseInterface
	// TimeNow contains a function that return the current time
	TimeNow func() time.Time
	// Config contains the dataservice global configuration
	Config config.Configuration
	// Log contains logger formatted
	Log logger.Logger
}

// Run archive every hostdata that is older than a amount
func (job *OciDataRetrieveJob) Run() {
	var merr error

	var listCompartments []model.OciCompartment

	dbProfiles, err := job.Database.GetOciProfiles(true)
	if err != nil {
		job.Log.Error(err)
		return
	}

	for _, p := range dbProfiles {
		// reset all the counters
		cntInstances := 0
		cntDatabases := 0
		cntLoadBalancers := 0
		cntNetworks := 0
		cntBlockVolume := 0
		cntBuckets := 0
		cntNFSs := 0

		customConfigProvider, tenancyOCID, err := job.getOciCustomConfigProviderAndTenancy(p.ID.Hex())
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		listCompartments, err = job.getOciProfileCompartments(tenancyOCID, customConfigProvider)
		if err != nil {
			merr = multierror.Append(merr, err)
			continue
		}

		bComputeClient := true
		bDatabaseClient := true
		bLoadBalancerClient := true
		bNetworkClient := true
		bBlockVolumeClient := true
		bObjectStorageClient := true
		bIdentityClient := true
		bFileStorageClient := true

		// first of all I create the client for each profile
		computeClient, err := core.NewComputeClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bComputeClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		databaseClient, err := database.NewDatabaseClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bDatabaseClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		loadBalancerClient, err := loadbalancer.NewLoadBalancerClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bLoadBalancerClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		networkClient, err := core.NewVirtualNetworkClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bNetworkClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		blockVolumeClient, err := core.NewBlockstorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bBlockVolumeClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		objectStorageClient, err := objectstorage.NewObjectStorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bObjectStorageClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		identityClient, err := identity.NewIdentityClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bIdentityClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		fileStorageClient, err := filestorage.NewFileStorageClientWithConfigurationProvider(customConfigProvider)
		if err != nil {
			bFileStorageClient = false
		} else {
			merr = multierror.Append(merr, err)
		}

		// then I have to retrieve informazion for each compartment
		for _, compartment := range listCompartments {
			if bComputeClient {
				// Create a request and dependent object(s).
				reqInstance := core.ListInstancesRequest{
					CompartmentId: common.String(compartment.CompartmentID),
				}

				// Send the request using the service client
				respInstance, err := computeClient.ListInstances(context.Background(), reqInstance)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				cntInstances = cntInstances + len(respInstance.Items)
			}

			if bDatabaseClient {
				// Create a request and dependent object(s).
				reqDbHomes := database.ListDbHomesRequest{
					CompartmentId: &compartment.CompartmentID,
				}

				// Send the request using the service client
				respDbHomes, err := databaseClient.ListDbHomes(context.Background(), reqDbHomes)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				for _, d := range respDbHomes.Items {
					// Create a request and dependent object(s).
					reqDatabases := database.ListDatabasesRequest{
						//SystemId:      common.String("ocid1.dbsystem.oc1.eu-frankfurt-1.abtheljsnloyfeoefvmw3mfraftfocrwhifjxuyu25gadsulfgdr2vcschia"),
						CompartmentId: &compartment.CompartmentID,
						DbHomeId:      common.String(*d.Id),
					}

					// Send the request using the service client
					respDatabases, err := databaseClient.ListDatabases(context.Background(), reqDatabases)
					if err != nil {
						merr = multierror.Append(merr, err)
						continue
					}

					cntDatabases = cntDatabases + len(respDatabases.Items)
				}
			}

			if bLoadBalancerClient {
				// Create a request and dependent object(s).
				reqLoadBalancer := loadbalancer.ListLoadBalancersRequest{
					CompartmentId: &compartment.CompartmentID,
				}

				// Send the request using the service client
				respLoadBalancer, err := loadBalancerClient.ListLoadBalancers(context.Background(), reqLoadBalancer)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				cntLoadBalancers = cntLoadBalancers + len(respLoadBalancer.Items)
			}

			if bNetworkClient {
				// Create a request and dependent object(s).
				reqNetwork := core.ListVcnsRequest{
					CompartmentId: common.String(compartment.CompartmentID),
				}

				// Send the request using the service client
				respNetwork, err := networkClient.ListVcns(context.Background(), reqNetwork)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				cntNetworks = cntNetworks + len(respNetwork.Items)
			}

			if bBlockVolumeClient {
				// Create a request and dependent object(s).
				reqBlockVolume := core.ListVolumesRequest{
					CompartmentId: &compartment.CompartmentID,
				}

				// Send the request using the service client
				respBlockVolume, err := blockVolumeClient.ListVolumes(context.Background(), reqBlockVolume)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				cntBlockVolume = cntBlockVolume + len(respBlockVolume.Items)
			}

			if bObjectStorageClient {
				reqNamespasce := objectstorage.GetNamespaceRequest{
					CompartmentId: common.String(compartment.CompartmentID),
				}

				// Send the request using the service client
				respNamespasce, err := objectStorageClient.GetNamespace(context.Background(), reqNamespasce)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				// Create a request and dependent object(s)
				reqBucket := objectstorage.ListBucketsRequest{
					CompartmentId: common.String(compartment.CompartmentID),
					NamespaceName: common.String(*respNamespasce.Value),
				}

				// Send the request using the service client
				respBucket, err := objectStorageClient.ListBuckets(context.Background(), reqBucket)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				cntBuckets = cntBuckets + len(respBucket.Items)
			}

			if bIdentityClient && bFileStorageClient {
				// Create a request and dependent object(s).
				reqAvailabilityDomain := identity.ListAvailabilityDomainsRequest{
					CompartmentId: common.String(compartment.CompartmentID),
				}
				// Send the request using the service client
				respAvailabilityDomain, err := identityClient.ListAvailabilityDomains(context.Background(), reqAvailabilityDomain)
				if err != nil {
					merr = multierror.Append(merr, err)
					continue
				}

				//fmt.Println("----------------- AVAILABILITY DOMAIN -----------------")
				//fmt.Println(resp1)

				for _, f := range respAvailabilityDomain.Items {
					// Create a request and dependent object(s).
					reqFileSystem := filestorage.ListFileSystemsRequest{
						AvailabilityDomain: common.String(*f.Name),
						CompartmentId:      common.String(compartment.CompartmentID),
					}
					//	fmt.Println("Availability Domain: ", *f.Name)

					// Send the request using the service client
					resp2, err := fileStorageClient.ListFileSystems(context.Background(), reqFileSystem)
					if err != nil {
						merr = multierror.Append(merr, err)
						continue
					}
					//fmt.Println("len: ", len(resp2.Items))

					cntNFSs = cntNFSs + len(resp2.Items)
				}
			}
		}

		var ociObject1, ociObject2, ociObject3, ociObject4, ociObject5, ociObject6 model.OciObject

		ociObject1.ObjectName = "instance"
		ociObject1.ObjectNumber = cntInstances
		ociObject2.ObjectName = "database"
		ociObject2.ObjectNumber = cntDatabases
		ociObject3.ObjectName = "network"
		ociObject3.ObjectNumber = cntLoadBalancers + cntNetworks
		ociObject4.ObjectName = "block volume"
		ociObject4.ObjectNumber = cntBlockVolume
		ociObject5.ObjectName = "bucket"
		ociObject5.ObjectNumber = cntBuckets
		ociObject6.ObjectName = "NFS"
		ociObject6.ObjectNumber = cntNFSs

		var ociList []model.OciObject

		ociList = append(ociList, ociObject1, ociObject2, ociObject3, ociObject4, ociObject5, ociObject6)

		var ociObjects model.OciObjects

		strError := merr.Error()
		if strError[0] != '0' {
			ociObjects = model.OciObjects{
				ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
				ProfileID: p.ID.Hex(),
				CreatedAt: time.Now().UTC(),
				Error:     merr.Error(),
				Objects:   ociList,
			}
		} else {
			ociObjects = model.OciObjects{
				ID:        primitive.NewObjectIDFromTimestamp(time.Now()),
				ProfileID: p.ID.Hex(),
				CreatedAt: time.Now().UTC(),
				Objects:   ociList,
			}
		}

		errDb := job.Database.AddOciObjects(ociObjects)

		if errDb != nil {
			job.Log.Error(errDb)
		}
	}
}
