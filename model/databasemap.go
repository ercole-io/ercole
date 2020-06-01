package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// DatabaseMap holds information about a database.
type DatabaseMap map[string]interface{}

// Licenses getter
func (db *DatabaseMap) Licenses() []LicenseMap {
	interf := (*db)["Licenses"]

	interfSlice := []interface{}(interf.(primitive.A))

	var licenses []LicenseMap
	for _, interf := range interfSlice {
		mapStrInterf := interf.(map[string]interface{})
		licenses = append(licenses, mapStrInterf)
	}

	return licenses
}

// Name getter
func (db *DatabaseMap) Name() string {
	name := (*db)["Name"]
	return name.(string)
}

// Features getter
func (db *DatabaseMap) Features() []FeatureMap {
	interf := (*db)["Features"]

	interfSlice := []interface{}(interf.(primitive.A))

	var features []FeatureMap
	for _, interf := range interfSlice {
		mapStrInterf := interf.(map[string]interface{})
		features = append(features, mapStrInterf)
	}

	return features
}

// DatabaseMapArrayAsMap return the equivalent map of the database array with Database.Name as Key
func DatabaseMapArrayAsMap(dbs []Database) map[string]Database {
	out := make(map[string]Database)
	for _, db := range dbs {
		out[db.Name] = db
	}
	return out
}

// HasEnterpriseLicense return true if the database has enterprise license.
func (db *DatabaseMap) HasEnterpriseLicense() bool {
	//Search for a enterprise license
	for _, lic := range db.Licenses() {
		if (lic.Name() == "Oracle ENT" || lic.Name() == "oracle ENT" || lic.Name() == "Oracle EXT" || lic.Name() == "oracle EXT") && lic.Count() > 0 {
			return true
		}
	}

	return false
}
