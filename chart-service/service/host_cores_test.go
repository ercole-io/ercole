package service

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
)

func TestGetHostsHistory_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	db := NewMockMongoDatabaseInterface(mockCtrl)
	as := ChartService{
		Database: db,
	}

	location := ""
	environment := ""
	olderThan := utils.MAX_TIME
	newerThan := utils.MIN_TIME
	expectedRes := []dto.HostCores{
		{
			Date:  utils.P("2020-04-15T00:00:00Z"),
			Cores: 1,
		},
	}

	db.EXPECT().GetHostCores(location, environment, olderThan, newerThan).
		Return(expectedRes, nil).Times(1)

	res, err := as.GetHostCores(location, environment, olderThan, newerThan)

	require.NoError(t, err)
	assert.Equal(t, expectedRes, res)
}
