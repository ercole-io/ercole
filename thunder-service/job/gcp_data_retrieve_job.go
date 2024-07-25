// Copyright (c) 2024 Sorint.lab S.p.A.
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
	"fmt"
	"sync"
	"time"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
	"google.golang.org/api/cloudresourcemanager/v1"
	"google.golang.org/api/option"
)

type GcpDataRetrieveJob struct {
	Database database.MongoDatabaseInterface
	Config   config.Configuration
	Log      logger.Logger
	Opt      *option.ClientOption
}

func (job *GcpDataRetrieveJob) GetClientOption(profile model.GcpProfile) *option.ClientOption {
	cred := fmt.Sprintf(`{
		"type": "service_account",
		"private_key": "%s",
		"client_email": "%s"
	  }`, profile.PrivateKey, profile.ClientEmail)

	opt := option.WithCredentialsJSON([]byte(cred))

	return &opt
}

func worker[T comparable](data T, wg *sync.WaitGroup, ch chan T) {
	defer wg.Done()
	ch <- data
}

func (job *GcpDataRetrieveJob) Run() {
	tstart := time.Now()

	ctx := context.TODO()
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	// defer cancel()

	seqValue, err := job.Database.GetLastGcpSeqValue()
	if err != nil {
		job.Log.Errorf("gcp seq value %v", err)
		return
	}

	seqValue = seqValue + 1

	job.Log.Debugf("seqvalue is %d", seqValue)

	profiles, err := job.Database.GetActiveGcpProfiles()
	if err != nil {
		job.Log.Errorf("gcp data retriever job active profile error %v", err)
		return
	}

	var profileWg sync.WaitGroup

	profileCh := make(chan model.GcpProfile, len(profiles))

	for _, profile := range profiles {
		profileWg.Add(1)

		go worker(profile, &profileWg, profileCh)
	}

	go func() {
		profileWg.Wait()
		close(profileCh)
	}()

	for profile := range profileCh {
		job.Opt = job.GetClientOption(profile)

		projects, err := job.Getprojects(ctx)
		if err != nil {
			gcperr := model.GcpError{
				SeqValue:  seqValue,
				ProfileID: profile.ID,
				Category:  "project error",
				CreatedAt: time.Now(),
				Msg:       err.Error(),
			}

			if err := job.AddError(gcperr); err != nil {
				job.Log.Warn(err)
			}

			continue
		}

		var projectWg sync.WaitGroup

		projectCh := make(chan *cloudresourcemanager.Project, len(projects))

		for _, project := range projects {
			if project.ProjectId == job.Config.ThunderService.GcpDataRetrieveJob.ProjectID {
				projectWg.Add(1)

				go worker(project, &projectWg, projectCh)
			}
		}

		go func() {
			projectWg.Wait()
			close(projectCh)
		}()

		for project := range projectCh {
			instances, err := job.GetInstances(ctx, project.ProjectId)
			if err != nil {
				gcperr := model.GcpError{
					SeqValue:  seqValue,
					ProfileID: profile.ID,
					Category:  "instance error",
					CreatedAt: time.Now(),
					Msg:       err.Error(),
				}

				if err := job.AddError(gcperr); err != nil {
					job.Log.Warn(err)
				}

				continue
			}

			var instanceWg sync.WaitGroup

			instanceCh := make(chan model.GcpInstance, len(instances))

			for _, instance := range instances {
				instanceWg.Add(1)

				go worker(model.GcpInstance{
					Instance:  instance,
					Project:   project,
					ProfileID: profile.ID,
				}, &instanceWg, instanceCh)
			}

			go func() {
				instanceWg.Wait()
				close(instanceCh)
			}()

			for instance := range instanceCh {
				var recInstanceWg sync.WaitGroup

				recInstanceCh := make(chan model.GcpRecommendation, len(instances))

				recInstanceWg.Add(1)

				go job.FetchGcpInstanceRightsizing(ctx, instance, seqValue, &recInstanceWg, recInstanceCh)

				go func() {
					recInstanceWg.Wait()
					close(recInstanceCh)
				}()

				for rec := range recInstanceCh {
					if err := job.AddRecommendation(rec); err != nil {
						gcperr := model.GcpError{
							SeqValue:  seqValue,
							ProfileID: profile.ID,
							Category:  "recommendation error",
							CreatedAt: time.Now(),
							Msg:       err.Error(),
						}

						if err := job.AddError(gcperr); err != nil {
							job.Log.Warn(err)
						}

						job.Log.Debugf("added new error - seqvalue: %d - category: %s - msg %d", gcperr.SeqValue, gcperr.Category, gcperr.Msg)

						continue
					}

					job.Log.Debugf("added new recommendation - seqvalue: %d - project name: %s - instanceID %d", rec.SeqValue, rec.ProjectName, rec.ResourceID)
				}

				var diskWg sync.WaitGroup

				diskCh := make(chan model.GcpDisk, len(instance.Disks))

				for _, attachedDisk := range instance.Disks {
					diskWg.Add(1)

					disk, err := job.GetDisk(ctx, project.ProjectId, attachedDisk.GetDeviceName(), instance.Zone())
					if err != nil {
						gcperr := model.GcpError{
							SeqValue:  seqValue,
							ProfileID: profile.ID,
							Category:  "recommendation error",
							CreatedAt: time.Now(),
							Msg:       err.Error(),
						}

						if err := job.AddError(gcperr); err != nil {
							job.Log.Warn(err)
						}

						job.Log.Debugf("added new error - seqvalue: %d - category: %s - msg %d", gcperr.SeqValue, gcperr.Category, gcperr.Msg)

						continue
					}

					gcpDisk := model.GcpDisk{
						InstanceID:   instance.GetId(),
						InstanceZone: instance.Zone(),
						MachineType:  instance.MachineType(),
						ProfileID:    instance.ProfileID,
						Project:      project,
						Disk:         disk,
					}

					go worker(gcpDisk, &diskWg, diskCh)
				}

				go func() {
					diskWg.Wait()
					close(diskCh)
				}()

				for disk := range diskCh {
					var recDiskWg sync.WaitGroup

					recDiskCh := make(chan model.GcpRecommendation, len(instance.Disks))

					recDiskWg.Add(1)

					go job.FetchGcpStorageDisk(ctx, disk, seqValue, &recDiskWg, recDiskCh)

					go func() {
						recDiskWg.Wait()
						close(recDiskCh)
					}()

					for rec := range recDiskCh {
						if err := job.AddRecommendation(rec); err != nil {
							gcperr := model.GcpError{
								SeqValue:  seqValue,
								ProfileID: profile.ID,
								Category:  "recommendation error",
								CreatedAt: time.Now(),
								Msg:       err.Error(),
							}

							if err := job.AddError(gcperr); err != nil {
								job.Log.Warn(err)
							}

							job.Log.Debugf("added new error - seqvalue: %d - category: %s - msg %d", gcperr.SeqValue, gcperr.Category, gcperr.Msg)

							continue
						}

						job.Log.Debugf("added new recommendation - seqvalue: %d - project name: %s - instanceID %d", rec.SeqValue, rec.ProjectName, rec.ResourceID)
					}
				}
			}
		}
	}

	dend := time.Since(tstart)

	job.Log.Debugf("gcp job took %v minutes", dend.Minutes())
}
