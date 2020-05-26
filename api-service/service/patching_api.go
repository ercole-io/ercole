// Copyright (c) 2019 Sorint.lab S.p.A.
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
	"Tags": {
		"dbname1": ["tag1", "tag2", ...],
		"dbname2": ["tag1", "tag2", ...],
		"..."
	}
}
*/
const DatabaseTagsAdderMarker = "DATABASE_TAGS_ADDER"
const DatabaseTagsAdderCode = `
	/*<DATABASE_TAGS_ADDER>*/
	hostdata.Extra.Databases.forEach(function addTag(db) {
		if (db.Name in vars.Tags) {
			db.Tags = vars.Tags[db.Name];
		}
	});
	/*</DATABASE_TAGS_ADDER>*/
`

// DatabaseLicensesFixerCode is the code used to modify the value of a license when the DatabaseLicensesFixerMarker is missing
// It assumes that the vars has the following structure:
/*
{
	"LicenseModifiers": {
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
	hostdata.Extra.Databases.forEach(function fixLicensesDb(db) {
		db.Licenses.forEach(function fixLicense(lic) {
			if (db.Name in vars.LicenseModifiers && lic.Name in vars.LicenseModifiers[db.Name]) {
				if (! ("OldCount" in lic)) {
					lic.OldCount = lic.Count;
				}
				lic.Count = vars.LicenseModifiers[db.Name][lic.Name];
			} else if ("OldCount" in lic) {
				lic.Count = lic.OldCount;
				delete lic.OldCount;
			}
		});
	});
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

// AddTagToDatabase add the tag to the database if it hasn't the tag
func (as *APIService) AddTagToDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface {
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
			"Tags": map[string]interface{}{
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
		if val, ok := pf.Vars["Tags"]; !ok {
			pf.Vars["Tags"] = make(map[string]interface{})
		} else if _, ok := val.(map[string]interface{}); !ok {
			pf.Vars["Tags"] = make(map[string]interface{})
		}
		tags := pf.Vars["Tags"].(map[string]interface{})

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

// DeleteTagOfDatabase delete the tag from the database if it hasn't the tag
func (as *APIService) DeleteTagOfDatabase(hostname string, dbname string, tagname string) utils.AdvancedErrorInterface {
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
	if val, ok := pf.Vars["Tags"]; !ok {
		return nil
	} else if _, ok := val.(map[string]interface{}); !ok {
		return nil
	}
	tags := pf.Vars["Tags"].(map[string]interface{})

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
	if aerr != nil {
		return aerr
	}

	//If it's found, patch it
	if data == nil {
		return nil
	}

	//Patch it
	if as.Config.DataService.LogDataPatching {
		as.Log.Printf("Patching %s hostdata with the patch %s\n", pf.Hostname, pf.ID)
	}
	utils.PatchHostdata(pf, data)

	//Save the patched data
	return as.Database.ReplaceHostData(data)
}

// SetLicenseModifier set the value of certain license to newValue
func (as *APIService) SetLicenseModifier(hostname string, dbname string, licenseName string, newValue int) utils.AdvancedErrorInterface {
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
			"LicenseModifiers": map[string]interface{}{
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
		if val, ok := pf.Vars["LicenseModifiers"]; !ok {
			pf.Vars["LicenseModifiers"] = make(map[string]interface{})
		} else if _, ok := val.(map[string]interface{}); !ok {
			pf.Vars["LicenseModifiers"] = make(map[string]interface{})
		}
		licenseModifiers := pf.Vars["LicenseModifiers"].(map[string]interface{})

		//Check the presence of the database with a slice inside
		if val, ok := licenseModifiers[dbname]; !ok {
			licenseModifiers[dbname] = make(map[string]interface{}, 0)
		} else if _, ok := val.(bson.A); !ok {
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

// DeleteLicenseModifier delete the modifier of a certain license
func (as *APIService) DeleteLicenseModifier(hostname string, dbname string, licenseName string) utils.AdvancedErrorInterface {
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
	if val, ok := pf.Vars["LicenseModifiers"]; !ok {
		return nil
	} else if _, ok := val.(map[string]interface{}); !ok {
		return nil
	}
	licenseModifiers := pf.Vars["LicenseModifiers"].(map[string]interface{})

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
