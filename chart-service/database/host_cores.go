package database

import (
	"context"
	"github.com/amreo/mu"
	"github.com/ercole-io/ercole/v2/chart-service/dto"
	"github.com/ercole-io/ercole/v2/utils"
	"go.mongodb.org/mongo-driver/bson"
	"time"
)

func (md *MongoDatabase) GetHostCores(location, environment string, olderThan, newerThan time.Time) ([]dto.HostCores, error) {
	cur, err := md.Client.Database(md.Config.Mongodb.DBName).Collection("hosts").Aggregate(
		context.TODO(),
		mu.MAPipeline(
			mu.APMatch(bson.M{
				"createdAt": bson.M{
					"$gte": newerThan,
					"$lte": olderThan,
				},
			}),
			FilterByLocationAndEnvironmentSteps(location, environment),
			mu.APGroup(
				bson.M{
					"_id":   bson.M{"$dateToString": bson.M{"format": "%Y-%m-%d", "date": "$createdAt"}},
					"cores": bson.M{"$sum": "$info.cpuCores"},
				},
			),
			mu.APSort(bson.M{
				"_id": 1,
			}),
			mu.APProject(bson.M{"date": bson.M{"$dateFromString": bson.M{"dateString": "$_id"}}, "cores": 1, "_id": 0}),
		),
	)
	if err != nil {
		return nil, utils.NewError(err, "DB ERROR")
	}

	var items = make([]dto.HostCores, 0)
	if err := cur.All(context.TODO(), &items); err != nil {
		return nil, utils.NewError(err, "Decode ERROR")
	}

	return items, nil
}
