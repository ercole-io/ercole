package dto

import "github.com/ercole-io/ercole/v2/model"

type OracleGrantDbaDto struct {
	OracleGrantDba model.OracleGrantDba `json:"oracleGrantDba" bson:"oracleGrantDba"`
	Hostname       string               `json:"hostname" bson:"hostname"`
	Databasename   string               `json:"databasename" bson:"databasename"`
}
