package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// ExtraInfoMap holds various informations.
type ExtraInfoMap map[string]interface{}

// Databases getter
func (ei *ExtraInfoMap) Databases() []DatabaseMap {
	interf := (*ei)["Databases"]

	interfSlice := []interface{}(interf.(primitive.A))

	var dbs []DatabaseMap
	for _, interf := range interfSlice {
		mapStrInterf := interf.(map[string]interface{})
		dbs = append(dbs, mapStrInterf)
	}

	return dbs
}

// SetDatabases setter
func (ei *ExtraInfoMap) SetDatabases(databases []DatabaseMap) {
	(*ei)["Databases"] = databases
}
