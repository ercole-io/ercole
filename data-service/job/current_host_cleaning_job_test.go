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
	"testing"

	gomock "github.com/golang/mock/gomock"

	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/logger"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestCurrentHostCleaningJobRun_SuccessNoOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	chcj := CurrentHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				CurrentHostCleaningJob: config.CurrentHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldCurrentHostnames(utils.P("2019-11-05T4:02:03Z")).Return([]string{}, nil).Times(1)

	chcj.Run()
}

func TestCurrentHostCleaningJobRun_SuccessOldCurrentHosts(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	chcj := CurrentHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				CurrentHostCleaningJob: config.CurrentHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldCurrentHostnames(utils.P("2019-11-05T4:02:03Z")).Return([]string{"superhost", "pippohost"}, nil).Times(1)
	db.EXPECT().DismissHost("superhost").Return(nil, nil).Times(1)
	db.EXPECT().DismissHost("pippohost").Return(nil, nil).Times(1)

	chcj.Run()
}

func TestCurrentHostCleaningJobRun_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	chcj := CurrentHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				CurrentHostCleaningJob: config.CurrentHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldCurrentHostnames(utils.P("2019-11-05T4:02:03Z")).Return([]string{"invalid"}, aerrMock).Times(1)
	db.EXPECT().DismissHost(gomock.Any()).Times(0)

	chcj.Run()
}

func TestCurrentHostCleaningJobRun_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	chcj := CurrentHostCleaningJob{
		TimeNow:  utils.Btc(utils.P("2019-11-05T14:02:03Z")),
		Database: db,
		Config: config.Configuration{
			DataService: config.DataService{
				CurrentHostCleaningJob: config.CurrentHostCleaningJob{
					HourThreshold: 10,
				},
			},
		},
		Log: logger.NewLogger("TEST"),
	}

	db.EXPECT().FindOldCurrentHostnames(utils.P("2019-11-05T4:02:03Z")).Return([]string{"superhost", "pippohost"}, nil).Times(1)
	db.EXPECT().DismissHost("superhost").Return(nil, aerrMock).Times(1)
	db.EXPECT().DismissHost("pippohost").Times(0)

	chcj.Run()
}
