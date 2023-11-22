package service

import (
	"testing"

	"github.com/ercole-io/ercole/v2/api-service/dto"
	"github.com/ercole-io/ercole/v2/config"
	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func TestListOracleGrantDbaByHostname_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := APIService{
		Database: db,
		Config:   config.Configuration{},
	}

	filter := dto.GlobalFilter{
		Location:    "",
		Environment: "",
		OlderThan:   utils.MAX_TIME,
	}

	expected := []dto.OracleGrantDbaDto{
		{
			Hostname:     "hostname",
			Databasename: "databasename",
			OracleGrantDba: model.OracleGrantDba{
				Grantee:     "test#001",
				AdminOption: "yes",
				DefaultRole: "no",
			},
		},
	}
	db.EXPECT().FindGrantDbaByHostname("hostname", filter).Return(expected, nil)

	res, err := as.ListOracleGrantDbaByHostname("hostname", filter)
	require.NoError(t, err)
	assert.Equal(t, expected, res)
}
