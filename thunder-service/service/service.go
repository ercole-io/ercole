// Copyright (c) 2023 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"math/rand"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"
	"github.com/ercole-io/ercole/v2/thunder-service/dto"
	"github.com/ercole-io/ercole/v2/thunder-service/job"

	"github.com/ercole-io/ercole/v2/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThunderServiceInterface interface {
	Init()
	GetOciNativeRecommendations(profiles []string) ([]model.OciNativeRecommendation, error)
	AddOciProfile(profile model.OciProfile) (*model.OciProfile, error)
	UpdateOciProfile(profile model.OciProfile) (*model.OciProfile, error)
	GetOciProfiles() ([]model.OciProfile, error)
	DeleteOciProfile(id primitive.ObjectID) error
	GetOciObjects() ([]model.OciObjects, error)
	GetOciRecommendations() ([]model.OciRecommendation, error)
	WriteOciRecommendationsXlsx(recommendations []model.OciRecommendation) (*excelize.File, error)
	ForceGetOciRecommendations() error
	GetOciRecommendationErrors(seqNum uint64) ([]model.OciRecommendationError, error)
	SelectOciProfile(profileId string, selected bool) error
	AddAwsProfile(profile model.AwsProfile) (*model.AwsProfile, error)
	UpdateAwsProfile(profile model.AwsProfile) (*model.AwsProfile, error)
	GetAwsProfiles() ([]model.AwsProfile, error)
	DeleteAwsProfile(id primitive.ObjectID) error
	SelectAwsProfile(profileId string, selected bool) error
	GetAwsRecommendations() ([]model.AwsRecommendation, error)
	GetLastAwsRecommendations() ([]model.AwsRecommendation, error)
	GetAwsRecommendationsBySeqValue(seqValue uint64) ([]model.AwsRecommendation, error)
	ForceGetAwsRecommendations() error
	WriteAwsRecommendationsXlsx(recommendations []dto.AwsRecommendationDto) (*excelize.File, error)
	GetLastAwsObjects() ([]model.AwsObject, error)
	AddAzureProfile(profile model.AzureProfile) (*model.AzureProfile, error)
	UpdateAzureProfile(profile model.AzureProfile) (*model.AzureProfile, error)
	GetAzureProfiles() ([]model.AzureProfile, error)
	DeleteAzureProfile(id primitive.ObjectID) error
	SelectAzureProfile(profileId string, selected bool) error
	GetAwsRDS() ([]model.AwsRDS, error)

	GetGcpProfiles() ([]model.GcpProfile, error)
	AddGcpProfile(profile dto.GcpProfileRequest) error
	SelectGcpProfile(idhex string, selected bool) error
	UpdateGcpProfile(profileID string, profile dto.GcpProfileRequest) error
	RemoveGcpProfile(profileID string) error

	ListGcpRecommendations() ([]model.GcpRecommendation, error)
	ForceGetGcpRecommendations()
	ListGcpError() ([]model.GcpError, error)
	CreateGcpRecommendationsXlsx() (*excelize.File, error)
}

type ThunderService struct {
	Config      config.Configuration
	Database    database.MongoDatabaseInterface
	TimeNow     func() time.Time
	Log         logger.Logger
	Random      *rand.Rand
	NewObjectID func() primitive.ObjectID
	Job         job.OciDataRetrieveJob
}

func (ts *ThunderService) Init() {
	ts.Random = rand.New(rand.NewSource(ts.TimeNow().UnixNano()))

	ts.NewObjectID = func() primitive.ObjectID {
		return primitive.NewObjectIDFromTimestamp(ts.TimeNow())
	}
}
