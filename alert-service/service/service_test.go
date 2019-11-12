package service

import (
	"testing"

	"github.com/amreo/ercole-services/model"
	"github.com/leandro-lugaresi/hub"
	"github.com/stretchr/testify/assert"

	"github.com/golang/mock/gomock"
)

func TestProcessHostDataInsertion_SuccessNewHost(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, nil).Do(func(alert model.Alert) {
		assert.Equal(t, "The server 'superhost1' was added to ercole", alert.Description)
		assert.Equal(t, p("2019-11-05T14:02:03Z"), alert.Date)
	}).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError1(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DatabaseError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, aerrMock).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Times(0)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}

func TestProcessHostDataInsertion_DiffHostError3(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := AlertService{
		Database: db,
		TimeNow:  btc(p("2019-11-05T14:02:03Z")),
	}

	db.EXPECT().FindHostData(str2oid("5dc3f534db7e81a98b726a52")).Return(hostData1, nil).Times(1)
	db.EXPECT().FindHostData(gomock.Any()).Times(0)
	db.EXPECT().FindMostRecentHostDataOlderThan("superhost1", p("2019-11-05T14:02:03Z")).Return(model.HostData{}, nil).Times(1)
	db.EXPECT().FindMostRecentHostDataOlderThan(gomock.Any(), gomock.Any()).Return(model.HostData{}, nil).Times(0)
	db.EXPECT().InsertAlert(gomock.Any()).Return(nil, aerrMock).Times(1)

	as.ProcessHostDataInsertion(hub.Fields{
		"id": str2oid("5dc3f534db7e81a98b726a52"),
	})
}
