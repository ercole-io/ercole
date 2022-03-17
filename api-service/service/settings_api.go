// Copyright (c) 2022 Sorint.lab S.p.A.
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

package service

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/ercole-io/ercole/v2/model"
	"github.com/ercole-io/ercole/v2/utils"
)

// loadManagedTechnologiesList loads the list of the managed techlogies from file and store it to as.TechnologyInfos.
func (as *APIService) loadManagedTechnologiesList() {
	// read the list content
	listContentRaw, err := ioutil.ReadFile(as.Config.ResourceFilePath + "/technologies/list.json")
	if err != nil {
		as.Log.Warnf("Unable to read %s: %v\n", as.Config.ResourceFilePath+"/technologies/list.json", err)
		return
	}

	// unmarshal to TechnologyInfos
	err = json.Unmarshal(listContentRaw, &as.TechnologyInfos)
	if err != nil {
		as.Log.Warnf("Unable to unmarshal %s: %v\n", as.Config.ResourceFilePath+"/technologies/list.json", err)
		return
	}

	// Load every image and encode it to base64
	for i, info := range as.TechnologyInfos {
		// read image content
		raw, err := ioutil.ReadFile(as.Config.ResourceFilePath + "/technologies/" + info.Product + ".png")
		if err != nil {
			as.Log.Warnf("Unable to read %s: %v\n", as.Config.ResourceFilePath+"/technologies/"+info.Product+".png", err)
		} else {
			// encode it!
			as.TechnologyInfos[i].Logo = base64.StdEncoding.EncodeToString(raw)
		}
	}
}

// GetDefaultDatabaseTags return the default list of database tags from configuration
func (as *APIService) GetDefaultDatabaseTags() ([]string, error) {
	return as.Config.APIService.DefaultDatabaseTags, nil
}

// GetErcoleFeatures return a map of active/inactive features
func (as *APIService) GetErcoleFeatures() (map[string]bool, error) {
	partialList, err := as.Database.GetHostsCountUsingTechnologies("", "", utils.MAX_TIME)
	if err != nil {
		return nil, err
	}

	out := map[string]bool{}

	out[model.TechnologyOracleDatabase] = partialList[model.TechnologyOracleDatabase] > 0
	out[model.TechnologyOracleExadata] = partialList[model.TechnologyOracleExadata] > 0

	return out, nil
}

// GetTechnologyList return the list of technologies
func (as *APIService) GetTechnologyList() ([]model.TechnologyInfo, error) {
	return as.TechnologyInfos, nil
}
