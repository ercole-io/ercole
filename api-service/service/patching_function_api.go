// Copyright (c) 2020 Sorint.lab S.p.A.
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

// Package service is a package that provides methods for querying data
package service

import (
	"strings"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DatabaseTagsAdderCode is the code used to add a tag to a database when the DatabaseTagsAdderMarker is missing
// It assumes that the vars has the following structure:
/*
{
	"tags": {
		"dbname1": ["tag1", "tag2", ...],
		"dbname2": ["tag1", "tag2", ...],
		"..."
	}
}
*/
const DatabaseTagsAdderMarker = "DATABASE_TAGS_ADDER"
const DatabaseTagsAdderCode = `
	/*<DATABASE_TAGS_ADDER>*/
	if (hostdata.features.oracle != undefined && hostdata.features.oracle.database != undefined) {
		hostdata.features.oracle.database.databases.forEach(function addTag(db) {
			if (db.name in vars.tags) {
				db.tags = vars.tags[db.name];
			}
		});
	}
	/*</DATABASE_TAGS_ADDER>*/
`

// DatabaseLicensesFixerCode is the code used to modify the value of a license when the DatabaseLicensesFixerMarker is missing
// It assumes that the vars has the following structure:
/*
{
	"licenseModifiers": {
		"dbname1": {
			"license1": 0,
			"license2": 2,
			...
		},
		"dbname2": {
			"license1": 0,
			"license2": 2,
			...
		},
		"..."
	}
}
*/
const DatabaseLicensesFixerMarker = "DATABASE_LICENSES_FIXER"
const DatabaseLicensesFixerCode = `
	/*<DATABASE_LICENSES_FIXER>*/
	if (hostdata.features.oracle != undefined && hostdata.features.oracle.database != undefined) { 
		hostdata.features.oracle.database.databases.forEach(function fixLicensesDb(db) {
			db.licenses.forEach(function fixLicense(lic) {
				if (db.name in vars.licenseModifiers && lic.name in vars.licenseModifiers[db.name]) {
					if (! ("oldCount" in lic)) {
						lic.oldCount = lic.count;
					}
					lic.count = vars.licenseModifiers[db.name][lic.name];
				} else if ("oldCount" in lic) {
					lic.count = lic.oldCount;
					delete lic.oldCount;
				}
			});
		});
	}
	/*</DATABASE_LICENSES_FIXER>*/
`

// GetPatchingFunction return the patching function specified in the hostname param
func (as *APIService) GetPatchingFunction(hostname string) (interface{}, utils.AdvancedErrorInterface) {
	//Check host existence
	exist, err := as.Database.ExistHostdata(hostname)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, utils.AerrHostNotFound
	}

	//Get the data
	return as.Database.FindPatchingFunction(hostname)
}

// SetPatchingFunction set the patching function of a host
func (as *APIService) SetPatchingFunction(hostname string, pf model.PatchingFunction) (interface{}, utils.AdvancedErrorInterface) {
	//Check host existence
	exist, err := as.Database.ExistHostdata(hostname)
	if err != nil {
		return nil, err
	} else if !exist {
		return nil, utils.AerrHostNotFound
	}

	//Get old patching function of the same host
	oldPf, err := as.Database.FindPatchingFunction(hostname)
	if err != nil {
		return nil, err
	}
	const DefaultPatchingCode = DatabaseTagsAdderCode + DatabaseLicensesFixerCode

	//Fill missing fields in the new pf
	pf.Hostname = hostname
	pf.CreatedAt = as.TimeNow()
	if oldPf.Hostname != hostname || oldPf.Code == "" {
		id := new(primitive.ObjectID)
		*id = primitive.NewObjectIDFromTimestamp(as.TimeNow())
		pf.ID = id
	} else {
		pf.ID = oldPf.ID
	}

	//Save the pf
	err = as.Database.SavePatchingFunction(pf)
	if err != nil {
		return nil, err
	}

	//Apply the patch
	err = as.ApplyPatch(pf)
	if err != nil {
		return nil, err
	}

	return pf.ID, nil
}

// DeletePatchingFunction delete the patching function of a host
func (as *APIService) DeletePatchingFunction(hostname string) utils.AdvancedErrorInterface {
	//Check host existence
	exist, err := as.Database.ExistHostdata(hostname)
	if err != nil {
		return err
	} else if !exist {
		return utils.AerrHostNotFound
	}

	//Delete the patching function
	return as.Database.DeletePatchingFunction(hostname)
}

// AddTagToOracleDatabase add the tag to the database if it hasn't the tag
func (as *APIService) AddTagToOracleDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface {
	//Find the patching function
	pf, err := as.Database.FindPatchingFunction(hostname)
	if err != nil {
		return err
	}

	//Check if the pf was found
	if pf.Hostname != hostname || pf.Code == "" {
		//No, initialze pf
		pf.Hostname = hostname
		pf.Code = DatabaseTagsAdderCode
		pf.Vars = map[string]interface{}{
			"tags": map[string]interface{}{
				dbname: []interface{}{
					tagname,
				},
			},
		}
		pf.CreatedAt = as.TimeNow()
	} else {
		//Check the presence of the marker in the code
		if !strings.Contains(pf.Code, "<"+DatabaseTagsAdderMarker+">") {
			pf.Code += DatabaseTagsAdderCode
		}

		//Check the presence and the type of Tags key in Vars. If not (re)initialize it!
		if val, ok := pf.Vars["tags"]; !ok {
			pf.Vars["tags"] = make(map[string]interface{})
		} else if _, ok := val.(map[string]interface{}); !ok {
			pf.Vars["tags"] = make(map[string]interface{})
		}
		tags := pf.Vars["tags"].(map[string]interface{})

		//Check the presence of the database with a slice inside
		if val, ok := tags[dbname]; !ok {
			tags[dbname] = make(bson.A, 0)
		} else if _, ok := val.(bson.A); !ok {
			tags[dbname] = make(bson.A, 0)
		}

		//Get the slice inside
		dbTags := tags[dbname].(bson.A)

		//Check if it already contains the tag
		for _, val := range dbTags {
			if val == tagname {
				return nil
			}
		}

		//Add it because the pf cannot contains it
		dbTags = append(dbTags, tagname)
		tags[dbname] = dbTags
	}

	// Save the modified patch
	if pf.ID == nil {
		oi := primitive.NewObjectIDFromTimestamp(as.TimeNow())
		pf.ID = &oi
	}
	err = as.Database.SavePatchingFunction(pf)
	if err != nil {
		return err
	}

	return as.ApplyPatch(pf)
}

// DeleteTagOfOracleDatabase delete the tag from the database if it hasn't the tag
func (as *APIService) DeleteTagOfOracleDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface {
	//Find the patching function
	pf, err := as.Database.FindPatchingFunction(hostname)
	if err != nil {
		return err
	}

	//Check if the pf was found
	if pf.Hostname != hostname || pf.Code == "" {
		return nil
	}

	//Check the presence of the marker in the code
	if !strings.Contains(pf.Code, "<"+DatabaseTagsAdderMarker+">") {
		pf.Code += DatabaseTagsAdderCode
	}

	//Check the presence and the type of Tags key in Vars. If not (re)initialize it!
	if val, ok := pf.Vars["tags"]; !ok {
		return nil
	} else if _, ok := val.(map[string]interface{}); !ok {
		return nil
	}
	tags := pf.Vars["tags"].(map[string]interface{})

	//Check the presence of the database with a slice inside
	if val, ok := tags[dbname]; !ok {
		return nil
	} else if _, ok := val.(bson.A); !ok {
		return nil
	}

	//Get the slice inside
	dbTags := tags[dbname].(bson.A)

	//Check if it contain the tag
	tagIndex := -1
	for i, val := range dbTags {
		if val == tagname {
			tagIndex = i
			break
		}
	}
	if tagIndex == -1 {
		return nil
	}

	//Remove it because the pf contains it
	tags[dbname] = append(dbTags[:tagIndex], dbTags[tagIndex+1:]...)

	//XXX: if there are no tags in the vars, we can remove the DatabaseTagsAdderCode. If there are no code, we can also remove the PF

	// Save the modified patch
	err = as.Database.SavePatchingFunction(pf)
	if err != nil {
		return err
	}

	return as.ApplyPatch(pf)
}

// ApplyPatch apply the patch pf to the relative host
func (as *APIService) ApplyPatch(pf model.PatchingFunction) utils.AdvancedErrorInterface {
	//Get the data
	data, aerr := as.Database.FindHostData(pf.Hostname)
	if aerr == utils.AerrHostNotFound {
		return nil
	} else if aerr != nil {
		return aerr
	}

	//Patch it
	if as.Config.DataService.LogDataPatching {
		as.Log.Printf("Patching %s hostdata with the patch %s\n", pf.Hostname, pf.ID)
	}
	data, aerr = utils.PatchHostdata(pf, data)
	if aerr != nil {
		return aerr
	}

	//Save the patched data
	return as.Database.ReplaceHostData(data)
}

// SetOracleDatabaseLicenseModifier set the value of certain license to newValue
func (as *APIService) SetOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string, newValue int) utils.AdvancedErrorInterface {
	//Find the patching function
	pf, err := as.Database.FindPatchingFunction(hostname)
	if err != nil {
		return err
	}

	//Check if the pf was found
	if pf.Hostname != hostname || pf.Code == "" {
		//No, initialze pf
		pf.Hostname = hostname
		pf.Code = DatabaseLicensesFixerCode
		pf.Vars = map[string]interface{}{
			"licenseModifiers": map[string]interface{}{
				dbname: map[string]interface{}{
					licenseName: newValue,
				},
			},
		}
		pf.CreatedAt = as.TimeNow()
	} else {
		//Check the presence of the marker in the code
		if !strings.Contains(pf.Code, "<"+DatabaseLicensesFixerMarker+">") {
			pf.Code += DatabaseLicensesFixerCode
		}

		//Check the presence and the type of LicenseModifiers key in Vars. If not (re)initialize it!
		if val, ok := pf.Vars["licenseModifiers"]; !ok {
			pf.Vars["licenseModifiers"] = make(map[string]interface{})
		} else if _, ok := val.(map[string]interface{}); !ok {
			pf.Vars["licenseModifiers"] = make(map[string]interface{})
		}
		licenseModifiers := pf.Vars["licenseModifiers"].(map[string]interface{})

		// Check the presence of the database with a slice inside
		switch licenseModifiers[dbname].(type) {
		case nil:
			licenseModifiers[dbname] = make(map[string]interface{}, 0)
		case map[string]interface{}:
		default:
			licenseModifiers[dbname] = make(map[string]interface{}, 0)
		}

		//Get the modifiers mof the db
		licenseModifiersOfDb := licenseModifiers[dbname].(map[string]interface{})

		//Check if it already contains the licenseName with the same value
		if val, ok := licenseModifiersOfDb[licenseName]; !ok {
			licenseModifiersOfDb[licenseName] = newValue
		} else if val == newValue {
			return nil
		} else {
			licenseModifiersOfDb[licenseName] = newValue
		}

	}

	// Save the modified patch
	if pf.ID == nil {
		oi := primitive.NewObjectIDFromTimestamp(as.TimeNow())
		pf.ID = &oi
	}
	err = as.Database.SavePatchingFunction(pf)
	if err != nil {
		return err
	}

	return as.ApplyPatch(pf)
}

// DeleteOracleDatabaseLicenseModifier delete the modifier of a certain license
func (as *APIService) DeleteOracleDatabaseLicenseModifier(hostname string, dbname string, licenseName string) utils.AdvancedErrorInterface {
	//Find the patching function
	pf, err := as.Database.FindPatchingFunction(hostname)
	if err != nil {
		return err
	}

	//Check if the pf was found
	if pf.Hostname != hostname || pf.Code == "" {
		return nil
	}

	//Check the presence of the marker in the code
	if !strings.Contains(pf.Code, "<"+DatabaseLicensesFixerMarker+">") {
		pf.Code += DatabaseLicensesFixerCode
	}

	//Check the presence and the type of LicenseModifiers key in Vars. If not (re)initialize it!
	if val, ok := pf.Vars["licenseModifiers"]; !ok {
		return nil
	} else if _, ok := val.(map[string]interface{}); !ok {
		return nil
	}
	licenseModifiers := pf.Vars["licenseModifiers"].(map[string]interface{})

	//Check the presence of the database with a slice inside
	if val, ok := licenseModifiers[dbname]; !ok {
		return nil
	} else if _, ok := val.(map[string]interface{}); !ok {
		return nil
	}

	//Get the modifiers mof the db
	licenseModifiersOfDb := licenseModifiers[dbname].(map[string]interface{})

	//Check if it already contains the licenseName with the same value
	if _, ok := licenseModifiersOfDb[licenseName]; !ok {
		return nil
	}
	//TODO: improved the cleanup of the license modifier
	delete(licenseModifiersOfDb, licenseName)

	// Save the modified patch
	err = as.Database.SavePatchingFunction(pf)
	if err != nil {
		return err
	}

	return as.ApplyPatch(pf)
}

// SearchOracleDatabaseLicenseModifiers search license modifiers
func (as *APIService) SearchOracleDatabaseLicenseModifiers(search string, sortBy string, sortDesc bool, page int, pageSize int) ([]map[string]interface{}, utils.AdvancedErrorInterface) {
	return as.Database.SearchOracleDatabaseLicenseModifiers(strings.Split(search, " "), sortBy, sortDesc, page, pageSize)
}
