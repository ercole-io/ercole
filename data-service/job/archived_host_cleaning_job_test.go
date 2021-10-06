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

package job

import (
	"testing"

	"go.mongodb.org/mongo-driver/bson/primitive"

	gomock "github.com/golang/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestArchivedHostCleaningJobRun_SuccessNoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ahcj := ArchivedHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				ArchivedHostCleaningJob: config.ArchivedHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldArchivedHosts(utils.P("2019-11-05T4:02:03Z")).Return([]primitive.ObjectID{}, nil).Times(1)

	ahcj.Run()
}

func TestArchivedHostCleaningJobRun_WithOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ahcj := ArchivedHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				ArchivedHostCleaningJob: config.ArchivedHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldArchivedHosts(utils.P("2019-11-05T4:02:03Z")).Return([]primitive.ObjectID{
		utils.Str2oid("5dcad88b24273d3489310b64"),
		utils.Str2oid("5dcad8933b243f80e2ed8538"),
	}, nil).Times(1)
	db.EXPECT().DeleteHostData(utils.Str2oid("5dcad88b24273d3489310b64")).Return(nil).Times(1)
	db.EXPECT().DeleteHostData(utils.Str2oid("5dcad8933b243f80e2ed8538")).Return(nil).Times(1)
	db.EXPECT().DeleteHostData(gomock.Any()).Times(0)

	ahcj.Run()
}

func TestArchivedHostCleaningJobRun_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ahcj := ArchivedHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				ArchivedHostCleaningJob: config.ArchivedHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldArchivedHosts(utils.P("2019-11-05T4:02:03Z")).Return([]primitive.ObjectID{
		utils.Str2oid("5dcad88b24273d3489310b64"),
		utils.Str2oid("5dcad8933b243f80e2ed8538"),
	}, aerrMock).Times(1)
	db.EXPECT().DeleteHostData(gomock.Any()).Times(0)

	ahcj.Run()
}

func TestArchivedHostCleaningJobRun_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	ahcj := ArchivedHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				ArchivedHostCleaningJob: config.ArchivedHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldArchivedHosts(utils.P("2019-11-05T4:02:03Z")).Return([]primitive.ObjectID{
		utils.Str2oid("5dcad88b24273d3489310b64"),
		utils.Str2oid("5dcad8933b243f80e2ed8538"),
	}, nil).Times(1)
	db.EXPECT().DeleteHostData(utils.Str2oid("5dcad88b24273d3489310b64")).Return(aerrMock).Times(1)
	db.EXPECT().DeleteHostData(gomock.Any()).Times(0)

	ahcj.Run()
}
