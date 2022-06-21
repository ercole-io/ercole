package database

import (
	"context"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
)

func (m *MongodbSuite) TestInsertOracleLicenseType_Success() {
	defer m.db.Client.Database(m.dbname).Collection("oracle_database_license_types").DeleteMany(context.TODO(), bson.M{})

	licenseTypes := []model.OracleDatabaseLicenseType{
		{
			ID:              "A00001",
			ItemDescription: "Oracle Test Item Description",
			Metric:          "Test Metric",
			Cost:            47500,
			Aliases:         []string{"Test Alias #01", "Test Alias #02"},
			Option:          false,
		},
		{
			ID:              "A00002",
			ItemDescription: "Oracle Test Item Description",
			Metric:          "Test Metric",
			Cost:            47501,
			Aliases:         []string{"Test Alias #01", "Test Alias #02"},
			Option:          false,
		},
	}

	for _, lt := range licenseTypes {
		err := m.db.InsertOracleLicenseType(lt)
		require.NoError(m.T(), err)
	}
}
