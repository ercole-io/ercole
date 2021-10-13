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

// Package service is a package that provides methods for manipulating host informations
package service

import (
	"math/rand"
	"time"

	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/thunder-service/database"

	"github.com/ercole-io/ercole/v2/config"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ThunderServiceInterface interface {
	Init()
	GetOciRecommendations(profiles []string) ([]model.Recommendation, error)
	AddOciProfile(profile model.OciProfile) (*model.OciProfile, error)
	UpdateOciProfile(profile model.OciProfile) (*model.OciProfile, error)
	GetOciProfiles() ([]model.OciProfile, error)
	DeleteOciProfile(id primitive.ObjectID) error
}

type ThunderService struct {
	Config      config.Configuration
	Database    database.MongoDatabaseInterface
	TimeNow     func() time.Time
	Log         logger.Logger
	Random      *rand.Rand
	NewObjectID func() primitive.ObjectID
}

func (as *ThunderService) Init() {
	as.Random = rand.New(rand.NewSource(as.TimeNow().UnixNano()))

	as.NewObjectID = func() primitive.ObjectID {
		return primitive.NewObjectIDFromTimestamp(as.TimeNow())
	}
}
