package service

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"

	"github.com/ercole-io/ercole/model"
	"github.com/ercole-io/ercole/utils"
)

// LoadManagedTechnologiesList loads the list of the managed techlogies from file and store it to as.TechnologyInfos.
func (as *APIService) LoadManagedTechnologiesList() {
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
func (as *APIService) GetDefaultDatabaseTags() ([]string, utils.AdvancedErrorInterface) {
	return as.Config.APIService.DefaultDatabaseTags, nil
}

// GetErcoleFeatures return a map of active/inactive features
func (as *APIService) GetErcoleFeatures() (map[string]bool, utils.AdvancedErrorInterface) {
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
func (as *APIService) GetTechnologyList() ([]model.TechnologyInfo, utils.AdvancedErrorInterface) {
	return as.TechnologyInfos, nil
}
