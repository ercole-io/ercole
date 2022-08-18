package dto

import "github.com/ercole-io/ercole/v2/model"

type OracleDatabasePluggableDatabase struct {
	Hostname                              string `json:"hostname" bson:"hostname"`
	model.OracleDatabasePluggableDatabase `json:"pdb" bson:"pdb"`
}
